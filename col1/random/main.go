// A rather simplistic random number generator.

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	maxUint = ^uint(0)
)

type rans struct {
	sync.Mutex
	m         map[uint]struct{} // use a struct{} instead of bool as it takes up no space
	ransToGen uint
}

// rand returns a failry random number using the nanosecond portion of a timestamp, and the hashing
// algorithm mentioned in
// https://crypto.stackexchange.com/questions/16219/cryptographic-hash-function-for-32-bit-length-input-keys
func rand() uint {
	// Take a timestamp
	ts := time.Now()
	// Obtain nanosecond portion of timestamp.
	nano := uint(ts.Nanosecond())
	var i uint
	for i = 0; i < 12; i++ {
		nano = ((nano>>8)^nano)*0x6b + i
	}

	return nano
}

// genRand takes the 'rans' struct pointer, a scaling factor, the upper bound and sync.WaitGroup
// pointer.  genRand will then loop, obtain a random number, scale the number using a left shift,
// and try to add the number to 'rans'.  If the number is not already present, it will decrement
// the counter of desired values in 'rans' and repeat.  If the counter of desired values in 'rans'
// reaches 0, genRand will signal completion and exit.
func genRand(r *rans, scale uint, upperBound uint, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		// Obtain random number.
		ran := rand()
		// Scale the number and check to ensure it stays within bounds.
		ran >>= scale
		if ran >= upperBound {
			continue
		}

		// Lock.
		r.Lock()

		// If number of values left to obtain == 0, signal finished and terminate, unlocking 'rans'.
		if r.ransToGen == 0 {
			r.Unlock()
			return
		}

		//   Try to add the number to 'rans'.  Check if it already exists.
		if _, ok := r.m[ran]; ok == true {
			// exists!
			r.Unlock()
			continue
		}

		// struct{}{} is passed as a struct{} value to our map.
		r.m[ran] = struct{}{}
		// Decrement number of values left to obtain.
		r.ransToGen--
		r.Unlock()
	}
}

// main parses command line input to obtain the upper bound for random numbers, the quantity of random
// numbers to generate and the output filename.  A 'rans' struct is instantiated, a scaling factor is
// computed based on the upper bound, and the number of workers are computed and started, generating
// random numbers and adding them to 'rans'.  main waits until all workers have finished, and writes out
// the contents of 'rans', one value per line, to the output file.
func main() {
	// Override default flag usage message.
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] N\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	// Obtain upper bound.
	upperBound := flag.Uint("bound", 10000000, "upper bound of generated random numbers")
	// Obtain output file.
	outputFilename := flag.String("out", "", "output file (default os.Stderr)")
	flag.Parse()

	// Obtain quantity of random numbers to generate.
	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "quantity of random numbers to generate not specified\n")
		os.Exit(1)
	}
	ransToGen, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse random numbers to generate: %v\n", err)
		os.Exit(1)
	}
	if ransToGen < 1 {
		fmt.Fprintf(os.Stderr, "invalid quantity of random numbers to generate: %v\n", ransToGen)
		os.Exit(1)
	}

	// Compute number of workers.  10000 numbers to generate per worker.
	workPerWorker := 10000
	workers := ransToGen / workPerWorker
	// We need at least one worker, it's a small business ;)
	if workers == 0 {
		workers++
	}

	// Initialize 'rans' struct with a hash of unsigned integers to structs{}, a mutex, and a counter.
	r := rans{}
	r.m = make(map[uint]struct{})
	r.ransToGen = uint(ransToGen)

	// Determine scale factor by right shifting maximum uint, until it is <= the upper bound.
	var scale uint = 0
	m := maxUint
	for ; m >= *upperBound; scale++ {
		m >>= 1
	}

	// Spool up workers.  sync.WaitGroup must be passed as a pointer.
	var wg sync.WaitGroup
	wg.Add(workers)
	start := time.Now()
	for i := 0; i < workers; i++ {
		go genRand(&r, scale, *upperBound, &wg)
	}

	// Wait for workers to finish.
	wg.Wait()
	stop := time.Now()
	fmt.Printf("generated %v random numbers in %v\n", ransToGen, stop.Sub(start))

	// Write 'rans' to output file, one per line, or to os.Stderr if outputFilename not provided.
	start = time.Now()
	var f *os.File
	if *outputFilename != "" {
		f, err = os.Create(*outputFilename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open %v for writing\n")
			os.Exit(1)
		}
	} else {
		f = os.Stderr
	}
	// Yea, one write per number...
	for ran, _ := range r.m {
		fmt.Fprintf(f, "%v\n", ran)
	}
	stop = time.Now()
	fmt.Printf("finished writing to file after %v\n", stop.Sub(start))
}
