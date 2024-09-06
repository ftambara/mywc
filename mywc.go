package main

import (
	"flag"
	"fmt"
)

type CountMode int

const (
	BYTES CountMode = iota
	CHARS
	LINES
	WORDS
	INVALID
)

func (m CountMode) String() string {
	switch m {
	case BYTES:
		return "BYTES"
	case CHARS:
		return "CHARS"
	case LINES:
		return "LINES"
	case WORDS:
		return "WORDS"
	default:
		panic("Invalid mode")
	}
}

func main() {
	bytesMode := flag.Bool("bytes", false, "print the byte counts")
	charsMode := flag.Bool("chars", false, "print the character counts")
	linesMode := flag.Bool("lines", false, "print the newline counts")
	wordsMode := flag.Bool("words", false, "print the word counts")

	flag.Parse()

	var countModes []CountMode
	if *bytesMode {
		countModes = append(countModes, BYTES)
	}
	if *charsMode {
		countModes = append(countModes, CHARS)
	}
	if *linesMode {
		countModes = append(countModes, LINES)
	}
	if *wordsMode {
		countModes = append(countModes, WORDS)
	}
	if len(countModes) == 0 {
		countModes = []CountMode{BYTES, CHARS, LINES, WORDS}
	}

	fmt.Printf("%v\n", countModes)
}
