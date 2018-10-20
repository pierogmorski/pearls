// A vectorsort package.

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"pearls/col1/vectorsort/vector"
	"strconv"
	"time"
)

// PopulateVecFromFile will take a filename and Vector, and read one unsigned integer value per line,
// checking the value against the range.  If a value is valid, SetBit will be invoked.  If a failure
// in reading occurrs, or if a value falls outside the specified range, or if SetBit fails, ReadFile
// will immediately return the number of entries processed prior to the failure, and an error to the
// caller.  Otherwise, the number of entries processed and success will be returned.
func PopulateVecFromFile(filename string, vec *vector.Vector) (uint, error) {
	var numEntries uint = 0

	f, err := os.Open(filename)
	if err != nil {
		return numEntries, err
	}
	defer f.Close()

	input := bufio.NewScanner(f)
	// For line in file.
	for input.Scan() {
		// Convert line to integer, checking for negative.
		val, err := strconv.Atoi(input.Text())
		if err != nil {
			return numEntries, errors.New(fmt.Sprintf("failed to read input file: %v", err))
		}
		if val < 0 {
			return numEntries, errors.New(fmt.Sprintf("invalid input read: %v", val))
		}

		// SetBit takes uint.
		_, err = vec.SetBit(uint(val))
		if err != nil {
			return numEntries, errors.New(fmt.Sprintf("failed to set %v: %v", val, err))
		}

		numEntries += 1
	}

	// Success.
	return numEntries, nil
}

// WriteVecToFile will take a filename and Vector, and write the contents of the vector (from least to
// greatest) to the output file, one value per line.  An error will be returned on any write failures,
// otherwise success will be returned.
func WriteVecToFile(filename string, vec *vector.Vector) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Loop over each bit in the vector.  Thus, loop invariant
	maxBits := vec.MaxBits()
	var i uint
	for i = 0; i < maxBits; i++ {
		set, err := vec.TestBit(uint(i))
		if err != nil {
			return errors.New(fmt.Sprintf("failed to write: %v", err))
		}

		if set {
			_, err = f.WriteString(fmt.Sprintf("%v\n", i))
			if err != nil {
				return errors.New(fmt.Sprintf("failed to write %v: %v", i, err))
			}
		}
	}

	// Success.
	return nil
}

func main() {
	start := time.Now()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}
	inFile := flag.String("in", "./input.txt", "input filename")
	outFile := flag.String("out", "./output.txt", "output filename")
	flag.Parse()

	var maxVals uint = 10000000
	vec := vector.New(maxVals)

	// Read inFile and populate the vector.  Any failure is treated as fatal.
	numEntries, err := PopulateVecFromFile(*inFile, vec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}

	// Write vector out in increasing order to outFile .  Any failure is treated as fatal.
	if err := WriteVecToFile(*outFile, vec); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}

	stop := time.Now()
	fmt.Printf("vectorsort took %v processing %v entries.\n", stop.Sub(start), numEntries)
}
