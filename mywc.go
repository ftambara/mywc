package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type countMode int

const (
	bytes countMode = iota
	chars
	lines
	words
	INVALID
)

func (m countMode) String() string {
	switch m {
	case bytes:
		return "bytes"
	case chars:
		return "chars"
	case lines:
		return "lines"
	case words:
		return "words"
	default:
		panic("Invalid mode")
	}
}

func parseModes(args []string) ([]countMode, error) {
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	bytesMode := f.Bool("bytes", false, "print the byte counts")
	charsMode := f.Bool("chars", false, "print the character counts")
	linesMode := f.Bool("lines", false, "print the newline counts")
	wordsMode := f.Bool("words", false, "print the word counts")

	err := f.Parse(args)
	if err != nil {
		return nil, err
	}

	if len(f.Args()) > 0 {
		return nil, fmt.Errorf("unrecognized arguments: %v", f.Args())
	}

	var countModes []countMode
	if *linesMode {
		countModes = append(countModes, lines)
	}
	if *wordsMode {
		countModes = append(countModes, words)
	}
	if *charsMode {
		countModes = append(countModes, chars)
	}
	if *bytesMode {
		countModes = append(countModes, bytes)
	}
	if len(countModes) == 0 {
		countModes = []countMode{lines, words, bytes}
	}

	return countModes, nil
}

func main() {
	countModes, err := parseModes(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", countModes)
}
