package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
)

const bufferSize = 8192

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

type inspector struct {
	modes  []countMode
	counts []uint
}

func newInspector(modes []countMode) (*inspector, error) {
	unique := make(map[countMode]bool, len(modes))
	for _, m := range modes {
		_, ok := unique[m]
		if ok {
			return nil, fmt.Errorf("duplicate mode: %v", m)
		}
		unique[m] = true
	}
	return &inspector{
		modes:  modes,
		counts: make([]uint, len(modes)),
	}, nil
}

func (in *inspector) resetCounts() {
	for i := range in.counts {
		in.counts[i] = 0
	}
}

func (in *inspector) readAll(r io.Reader) error {
	var (
		lineCount uint
		wordCount uint
		byteCount uint
		err       error
	)
	if slices.Equal(in.modes, []countMode{byLines}) {
		lineCount, err = countLines(r)
	} else if slices.Equal(in.modes, []countMode{byWords}) {
		wordCount, err = countWords(r)
	} else if slices.Equal(in.modes, []countMode{byBytes}) {
		byteCount, err = countBytes(r)
	} else {
		var counts [3]uint
		counts, err = countLinesWordsBytes(r)
		lineCount = counts[0]
		wordCount = counts[1]
		byteCount = counts[2]
	}
	if err != nil {
		return err
	}
	for i, m := range in.modes {
		switch m {
		case byLines:
			in.counts[i] = lineCount
		case byWords:
			in.counts[i] = wordCount
		case byBytes:
			in.counts[i] = byteCount
		}
	}
	return nil
}

func countLines(r io.Reader) (uint, error) {
	count := 0
	buffer := make([]byte, bufferSize)
	for {
		n, err := r.Read(buffer)
		count += bytes.Count(buffer[:n], []byte("\n"))
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return uint(count), err
		}
	}
	return uint(count), nil
}

func countWords(r io.Reader) (uint, error) {
	count := 0
	buffer := make([]byte, bufferSize)
	inWhitespace := true
	for {
		n, err := r.Read(buffer)
		for _, b := range buffer[:n] {
			if b == ' ' || b == '\n' || b == '\t' {
				if !inWhitespace {
					count++
					inWhitespace = true
				}
			} else {
				inWhitespace = false
			}
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return uint(count), err
		}
	}
	// EOF counts as WS
	if !inWhitespace {
		count++
	}
	return uint(count), nil
}

func countBytes(r io.Reader) (uint, error) {
	count := 0
	buffer := make([]byte, bufferSize)
	for {
		n, err := r.Read(buffer)
		count += n
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return uint(count), err
		}
	}
	return uint(count), nil
}

func countLinesWordsBytes(r io.Reader) ([3]uint, error) {
	var (
		lines        int
		words        int
		bytes        int
		inWhitespace = true
	)
	buffer := make([]byte, bufferSize)
	for {
		n, err := r.Read(buffer)
		bytes += n
		for _, b := range buffer[:n] {
			switch b {
			case '\n':
				lines++
				fallthrough
			case ' ', '\t':
				if !inWhitespace {
					words++
					inWhitespace = true
				}
			default:
				inWhitespace = false
			}
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return [3]uint{}, err
		}
	}
	// EOF counts as WS
	if !inWhitespace {
		words++
	}
	return [...]uint{uint(lines), uint(words), uint(bytes)}, nil
}

type namedReader struct {
	name string
	r    io.Reader
}

func main() {
	conf, err := parseArgs(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	namedReaders := make([]namedReader, max(len(conf.files), 1))
	if len(conf.files) == 0 {
		namedReaders[0] = namedReader{r: os.Stdin}
	} else {
		for i, name := range conf.files {
			f, err := os.Open(name)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			namedReaders[i] = namedReader{name: name, r: f}
		}
	}

	insp, err := newInspector(conf.countModes)
	if err != nil {
		log.Fatal(err)
	}
	for _, nr := range namedReaders {
		err = insp.readAll(nr.r)
		if err != nil {
			log.Fatal(err)
		}
		for _, c := range insp.counts {
			fmt.Printf("%v\t", c)
		}
		fmt.Printf("%v\n", nr.name)
		insp.resetCounts()
	}
}
