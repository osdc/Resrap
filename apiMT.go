package resrap

type codeGenReq struct {
	name      string
	startnode string
	seed      uint64
	tokens    int
	id        string
}

// CodeGenRes contains the process id along with the code generated returned from ResrapMT
type CodeGenRes struct {
	Code string
	Id   string
}

// ResrapMT is the multithreaded version of ResrapMT
type ResrapMT struct {
	languageGraph map[string]lang
	poolsize      int // Number of threads in the pool
	waitqueuesize int
	pendingjobs   chan codeGenReq
	codeChannel   chan CodeGenRes
}

// NewResrapMT creates and returns a new Resrap MultiThreaded instance.
// The returned instance starts with no loaded grammars.
func NewResrapMT(poolsize, waitqueuesize int) *ResrapMT {
	return &ResrapMT{
		languageGraph: make(map[string]lang),
		waitqueuesize: waitqueuesize,
		poolsize:      poolsize,
		pendingjobs:   make(chan codeGenReq, waitqueuesize),
		codeChannel:   make(chan CodeGenRes),
	}
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

// GenerateRandom schedules a job to generate content from the grammar identified by 'name'.
// startingNode: the starting symbol in the grammar for generation.
// id: a user-defined process ID that will be associated with the generated content.
// tokens: number of tokens to generate.
// The generation is non-deterministic (random). The generated content will be sent
// asynchronously to the CodeChannel. Users must provide a unique process ID and
// retrieve the result via the get channel function.
func (r *ResrapMT) GenerateRandom(id, name, startingNode string, tokens int) {
	req := codeGenReq{name: name, startnode: startingNode, tokens: tokens, seed: 0, id: id}
	r.pendingjobs <- req
}

// GenerateWithSeeded schedules a job to generate content from the grammar identified by 'name'.
// startingNode: the starting symbol in the grammar for generation.
// seed: a numeric seed to make generation deterministic.
// id: a user-defined process ID that will be associated with the generated content.
// tokens: number of tokens to generate.
// The generated content will be sent asynchronously to the CodeChannel. Users must
// provide a unique process ID and retrieve the result via the get channel function.
func (r *ResrapMT) GenerateWithSeeded(id, name, startingNode string, seed uint64, tokens int) {
	req := codeGenReq{name: name, startnode: startingNode, tokens: tokens, seed: seed, id: id}
	r.pendingjobs <- req
}

func (r *ResrapMT) mtparser() {
	for job := range r.pendingjobs {
		p := newPRNG(job.seed)
		l, ok := r.languageGraph[job.name]
		var code string
		if ok && l.graph != nil {
			code = l.graph.GraphWalk(&p, job.startnode, job.tokens)
		}
		r.codeChannel <- CodeGenRes{Code: code, Id: job.id}
	}
}

// ShutDownResrap gracefully ends the server goroutines running
func (r *ResrapMT) ShutDownResrap() {
	close(r.pendingjobs)
}

// StartResrap boots up goroutines as your specified threadpool
func (r *ResrapMT) StartResrap() {
	for i := 0; i < r.poolsize; i++ {
		go r.mtparser()
	}
}
