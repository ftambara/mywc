package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/ftambara/mywc/internal/counters"
)

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
	files []string
	opts  counters.Options
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

	conf.opts = counters.NewOptions(*linesMode, *wordsMode, *charsMode, *bytesMode)
	conf.files = f.Args()
	return conf, nil
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

	count := conf.opts.SelectCountingFn()
	if err != nil {
		log.Fatal(err)
	}
	for _, nr := range namedReaders {
		stats, err := count(nr.r)
		if err != nil {
			log.Fatal(err)
		}
		stats.Print(conf.opts, nr.name)
	}
}
