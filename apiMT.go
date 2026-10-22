package resrap

import (
	"errors"
	"sync"
)

type codeGenReq struct {
	name      string
	startnode string
	seed      uint64
	tokens    int
	complete  bool
	id        string
}

// CodeGenRes contains the process id along with the code generated returned from ResrapMT
type CodeGenRes struct {
	Code string
	Id   string
}

// ErrShutDown is returned when a job is submitted after ShutDownResrap.
var ErrShutDown = errors.New("resrap: instance has been shut down")

// ResrapMT is the multithreaded version of ResrapMT
type ResrapMT struct {
	languageGraph map[string]lang
	poolsize      int // Number of threads in the pool
	waitqueuesize int
	pendingjobs   chan codeGenReq
	codeChannel   chan CodeGenRes
	enqueue       chan CodeGenRes
	shutdown      chan struct{}
	startOnce     sync.Once
	shutdownOnce  sync.Once
	mu            sync.RWMutex
	closed        bool
	wg            sync.WaitGroup
	dispatcherWG  sync.WaitGroup
}

// NewResrapMT creates and returns a new Resrap MultiThreaded instance.
// The returned instance starts with no loaded grammars.
func NewResrapMT(poolsize, waitqueuesize int) *ResrapMT {
	if poolsize < 1 {
		poolsize = 1
	}
	if waitqueuesize < 1 {
		waitqueuesize = 1
	}
	r := &ResrapMT{
		languageGraph: make(map[string]lang),
		waitqueuesize: waitqueuesize,
		poolsize:      poolsize,
		pendingjobs:   make(chan codeGenReq, waitqueuesize),
		codeChannel:   make(chan CodeGenRes),
		enqueue:       make(chan CodeGenRes),
		shutdown:      make(chan struct{}),
	}
	// Results are queued in an unbounded internal buffer so workers never
	// block when the consumer is slower than the producer (#5).
	r.dispatcherWG.Add(1)
	go r.dispatch()
	return r
}

// GetCodeChannel is the main endpoint for the user to access the processed tokens
func (r *ResrapMT) GetCodeChannel() chan CodeGenRes {
	return r.codeChannel
}

// ParseGrammar parses a grammar string and stores it under the given name.
// name: a unique identifier for this grammar (e.g., "C").
// grammar: the grammar definition as a string.
func (r *ResrapMT) ParseGrammar(name, grammar string) error {
	l := newLang()
	err := l.ParserString(grammar)
	if err != nil {
		return err
	}
	r.languageGraph[name] = l
	r.languageGraph[name].graph.Normalize()
	return nil
}

// ParseGrammarFile parses a grammar from a file and stores it under the given name.
// name: a unique identifier for this grammar (e.g., "C").
// location: path to the grammar file.
func (r *ResrapMT) ParseGrammarFile(name, location string) error {
	l := newLang()
	err := l.ParserFile(location)
	if err != nil {
		return err
	}
	r.languageGraph[name] = l
	r.languageGraph[name].graph.Normalize()
	return nil
}

// AddWords attaches a list of words to a rule in a parsed grammar.
func (r *ResrapMT) AddWords(name, rule string, words []string) error {
	l, ok := r.languageGraph[name]
	if !ok || l.graph == nil {
		return errUnknownLanguage(name)
	}
	return l.graph.SetWords(rule, words)
}

// AddKeywords is an alias for AddWords.
func (r *ResrapMT) AddKeywords(name, rule string, words []string) error {
	return r.AddWords(name, rule, words)
}

type languageNotFoundError string

func (e languageNotFoundError) Error() string {
	return "resrap: grammar '" + string(e) + "' is not loaded"
}
func errUnknownLanguage(name string) error { return languageNotFoundError(name) }

func (r *ResrapMT) submit(req codeGenReq) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.closed {
		return ErrShutDown
	}
	select {
	case r.pendingjobs <- req:
		return nil
	case <-r.shutdown:
		return ErrShutDown
	}
}

// GenerateRandom schedules a job to generate content from the grammar identified by 'name'.
// startingNode: the starting symbol in the grammar for generation.
// id: a user-defined process ID that will be associated with the generated content.
// tokens: number of tokens to generate.
// The generated content will be sent asynchronously to the CodeChannel. Users must
// provide a unique process ID and retrieve the result via the get channel function.
func (r *ResrapMT) GenerateRandom(id, name, startingNode string, tokens int) error {
	return r.submit(codeGenReq{name: name, startnode: startingNode, tokens: tokens, seed: 0, id: id})
}

// GenerateWithSeeded schedules a job to generate content from the grammar identified by 'name'.
// startingNode: the starting symbol in the grammar for generation.
// seed: a numeric seed to make generation deterministic.
// id: a user-defined process ID that will be associated with the generated content.
// tokens: number of tokens to generate.
// The generated content will be sent asynchronously to the CodeChannel. Users must
// provide a unique process ID and retrieve the result via the get channel function.
func (r *ResrapMT) GenerateWithSeeded(id, name, startingNode string, seed uint64, tokens int) error {
	return r.submit(codeGenReq{name: name, startnode: startingNode, tokens: tokens, seed: seed, id: id})
}

// GenerateComplete schedules a completeness-aware generation job. Output is
// not truncated in the middle of the starting rule.
func (r *ResrapMT) GenerateComplete(id, name, startingNode string, tokens int) error {
	return r.submit(codeGenReq{name: name, startnode: startingNode, tokens: tokens, seed: 0, id: id, complete: true})
}

// GenerateCompleteWithSeeded is the deterministic variant of GenerateComplete.
func (r *ResrapMT) GenerateCompleteWithSeeded(id, name, startingNode string, seed uint64, tokens int) error {
	return r.submit(codeGenReq{name: name, startnode: startingNode, tokens: tokens, seed: seed, id: id, complete: true})
}

func (r *ResrapMT) mtparser() {
	defer r.wg.Done()
	for job := range r.pendingjobs {
		p := newPRNG(job.seed)
		l, ok := r.languageGraph[job.name]
		var code string
		if ok && l.graph != nil {
			if job.complete {
				code = l.graph.GraphWalkComplete(&p, job.startnode, job.tokens)
			} else {
				code = l.graph.GraphWalk(&p, job.startnode, job.tokens)
			}
		}
		r.enqueue <- CodeGenRes{Code: code, Id: job.id}
	}
}

// dispatch implements the unbounded result queue. Workers enqueue results
// without ever blocking on the consumer; the dispatcher only blocks while
// delivering to the public codeChannel.
func (r *ResrapMT) dispatch() {
	defer r.dispatcherWG.Done()
	defer close(r.codeChannel)

	var queue []CodeGenRes
	for {
		var out chan CodeGenRes
		var next CodeGenRes
		if len(queue) > 0 {
			out = r.codeChannel
			next = queue[0]
		}
		select {
		case job, ok := <-r.enqueue:
			if !ok {
				for _, pending := range queue {
					r.codeChannel <- pending
				}
				return
			}
			queue = append(queue, job)
		case out <- next:
			queue = queue[1:]
		}
	}
}

// ShutDownResrap gracefully ends the server goroutines running. It is safe to
// call more than once. Pending results are drained to the code channel; the
// channel is closed once all workers have finished.
func (r *ResrapMT) ShutDownResrap() {
	r.shutdownOnce.Do(func() {
		r.mu.Lock()
		r.closed = true
		r.mu.Unlock()
		close(r.shutdown)

		r.mu.Lock()
		close(r.pendingjobs)
		r.mu.Unlock()

		r.wg.Wait() // workers finish and send their last result
		close(r.enqueue)
		r.dispatcherWG.Wait()
	})
}

// StartResrap boots up goroutines as your specified threadpool
func (r *ResrapMT) StartResrap() {
	r.startOnce.Do(func() {
		r.wg.Add(r.poolsize)
		for i := 0; i < r.poolsize; i++ {
			go r.mtparser()
		}
	})
}
