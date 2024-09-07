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

type config struct {
	files      []string
	countModes []countMode
}

func parseArgs(args []string) (config, error) {
	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	bytesMode := f.Bool("bytes", false, "print the byte counts")
	charsMode := f.Bool("chars", false, "print the character counts")
	linesMode := f.Bool("lines", false, "print the newline counts")
	wordsMode := f.Bool("words", false, "print the word counts")

	conf := config{}

	err := f.Parse(args)
	if err != nil {
		return conf, err
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

	conf.countModes = countModes
	conf.files = f.Args()
	return conf, nil
}

func main() {
	conf, err := parseArgs(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("modes: %v, files: %v\n", conf.countModes, conf.files)
}
