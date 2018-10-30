// An anagram set finder.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type runes []rune

// Len method returns the length of the runes list.
func (r runes) Len() int {
	return len(r)
}

// Less method compares two runes lexicographically.
func (r runes) Less(i, j int) bool {
	return r[i] < r[j]
}

// Swap method swaps two runes.
func (r runes) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// genAnagramKey takes a word, and returns the word converted to lowercase and sorted by letter.
func genAnagramKey(word string) string {
	// Convert word to lowercase rune slice.
	wordSlice := runes(strings.ToLower(word))
	// Sort rune slice and return joined result to form anagram key.
	sort.Sort(wordSlice)
	return string(wordSlice)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s input_file output_file\nIf no output_file provided, will use stdout.\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	// Get a dictionary input file, and output file from CLI.
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	dictFilename := flag.Arg(0)
	outFilename := flag.Arg(1)

	// Open dictionary file.
	words := []string{}
	dictFile, err := os.Open(dictFilename)
	if err != nil {
		log.Fatalf("failed to open dictionary file: %v\n", err)
	}
	defer dictFile.Close()

	// Read words (one per line) from dictFile and populate 'words' list.
	scanner := bufio.NewScanner(dictFile)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	keyedWords := make(map[string][]string)
	// For each word in 'words' list.
	for _, word := range words {
		anagramKey := genAnagramKey(word)
		//Key 'keyedWords' by anagramKey, and append original word to value list.
		keyedWords[anagramKey] = append(keyedWords[anagramKey], word)
	}

	// Free 'words' list.
	words = nil

	// Sort keys of 'keyedWords', appending to 'anagramKeys' list.
	var anagramKeys []string
	for anagramKey := range keyedWords {
		anagramKeys = append(anagramKeys, anagramKey)
	}
	sort.Strings(anagramKeys)

	// Open output file.  If no output file specified, use os.Stdout.
	outFile := os.Stdout
	if len(outFilename) != 0 {
		outFile, err = os.Create(outFilename)
		if err != nil {
			log.Fatalf("failed to open outpuf file: %v\n", err)
		}
		defer outFile.Close()
	}

	// For each key in 'anagramKeys' list, output list of words (anagram set) to output file.
	for _, anagramKey := range anagramKeys {
		fmt.Fprintf(outFile, "%v\n", strings.Join(keyedWords[anagramKey], " "))
	}
}
