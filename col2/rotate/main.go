// slice rotation

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func fmtSlice(slice []rune) string {
	sliceStr := ""
	for _, r := range slice {
		sliceStr += fmt.Sprintf("%c", r)
	}
	return sliceStr
}

func getDirection(right bool) string {
	if right {
		return "right"
	}
	return "left"
}

// reverse takes a slice of runes, and a range (i and j) and reverses slice[i : j].
func reverse(slice []rune, i, j uint) {
	// If i >= j, no reverse to perform.
	if i >= j {
		return
	}
	// For i < j, swap slice[i] with slice[j]
	for ; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// rotate takes a slice of runes, a boolean right direction (default is left), and the number of
// positions to rotate the slice in place.
func rotate(slice []rune, right bool, rotPos uint) {
	sliceLen := uint(len(slice))

	// Don't rotate an empty slice.
	if sliceLen == 0 {
		return
	}

	// If rotPos > sliceLen, treat it as rotPos/sliceLen full rotations and only rotate the
	// remainder.
	rotPos = rotPos % sliceLen

	// Don't try to rotate by 0.
	if rotPos == 0 {
		return
	}

	// If rotating right, perform a full reversal first, otherwise do it last.
	if right {
		reverse(slice, 0, sliceLen-1)
	} else {
		defer reverse(slice, 0, sliceLen-1)
	}
	reverse(slice, 0, rotPos-1)
	reverse(slice, rotPos, sliceLen-1)
}

func main() {
	// Override default flag usage message.
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] string\n", os.Args[0])
		flag.PrintDefaults()
	}

	// Grab number of positions to rotate, and direction.
	rotPos := flag.Uint("rot", 0, "rotate slice by specified positions")
	right := flag.Bool("right", false, "rotate to the right")
	flag.Parse()

	// Grab string from which to form slice.
	if flag.NArg() < 1 {
		log.Fatal("string not provided")
	}

	sliceStr := flag.Arg(0)
	slice := []rune(sliceStr)

	fmt.Printf("original slice:\t\t\t\"%s\"\n", fmtSlice(slice))

	rotate(slice, *right, *rotPos)

	fmt.Printf("slice rotated %s by %d:\t\"%s\"\n", getDirection(*right), *rotPos, fmtSlice(slice))
}
