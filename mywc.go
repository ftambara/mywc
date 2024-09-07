package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

const bufferSize = 1024

type countMode int

const (
	byBytes countMode = iota
	byChars
	byLines
	byWords
)

func (m countMode) String() string {
	switch m {
	case byBytes:
		return "bytes"
	case byChars:
		return "chars"
	case byLines:
		return "lines"
	case byWords:
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
		countModes = append(countModes, byLines)
	}
	if *wordsMode {
		countModes = append(countModes, byWords)
	}
	if *charsMode {
		countModes = append(countModes, byChars)
	}
	if *bytesMode {
		countModes = append(countModes, byBytes)
	}
	if len(countModes) == 0 {
		countModes = []countMode{byLines, byWords, byBytes}
	}

	conf.countModes = countModes
	conf.files = f.Args()
	return conf, nil
}

func countLines(r io.Reader) (uint, error) {
	scanner := bufio.NewScanner(r)
	var count uint = 0
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		return count, err
	}
	return count, nil
}

func main() {
	conf, err := parseArgs(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("modes: %v, files: %v\n", conf.countModes, conf.files)
}
