package main

import (
	"fmt"
	"time"

	"github.com/osdc/resrap"
)

// func main() {
// 	// Create a new single-threaded Resrap
// 	resrapSync := resrap.NewResrap()

// 	// Parse grammar file (single-threaded)
// 	fmt.Println("Parsing grammar...")
// 	startParse := time.Now()
// 	resrapSync.ParseGrammarFile("C", "../example/c.g4")
// 	parseElapsed := time.Since(startParse)
// 	fmt.Printf("Grammar parsed in %v\n", parseElapsed)

// 	// Generate a huge sequence (4 million tokens)
// 	numTokens := 40000000
// 	fmt.Printf("Generating %d tokens (single-threaded)...\n", numTokens)
// 	startGen := time.Now()
// 	_ = resrapSync.GenerateRandom("C", "program", numTokens)
// 	genElapsed := time.Since(startGen)

//		fmt.Printf("Generated %d tokens in %v\n", numTokens, genElapsed)
//		//fmt.Printf("First 500 chars of output: \n%s\n", code[:500])
//	}
func main() {
	// Configuration
	numJobs := 100
	tokensPerJob := 1000000

	// --- SINGLE-THREADED RESRAP ---
	fmt.Println("Starting single-threaded test...")
	resrapSync := resrap.NewResrap()
	resrapSync.ParseGrammarFile("C", "../example/c.g4")

	start := time.Now()
	for i := 0; i < numJobs; i++ {
		_ = resrapSync.GenerateRandom("C", "program", tokensPerJob)
	}
	elapsedSync := time.Since(start)
	fmt.Printf("Single-threaded: %v for %d jobs of %d tokens\n", elapsedSync, numJobs, tokensPerJob)

	// --- MULTI-THREADED RESRAPMT ---
	fmt.Println("Starting multi-threaded test...")
	poolSize := 20
	waitQueueSize := 10000
	resrapMT := resrap.NewResrapMT(poolSize, waitQueueSize)
	resrapMT.ParseGrammarFile("C", "../example/c.g4")
	resrapMT.StartResrap()

	start = time.Now()
	for i := 0; i < numJobs; i++ {
		jobID := fmt.Sprintf("%d", i)
		resrapMT.GenerateRandom(jobID, "C", "program", tokensPerJob)
	}

	// Collect all results
	for i := 0; i < numJobs; i++ {
		<-resrapMT.GetCodeChannel()
	}
	elapsedMT := time.Since(start)
	fmt.Printf("Multi-threaded: %v for %d jobs of %d tokens\n", elapsedMT, numJobs, tokensPerJob)
}
