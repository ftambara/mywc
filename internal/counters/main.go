package counters

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"
)

const bufferSize = 8192

type Options struct {
	CountLines bool
	CountWords bool
	CountChars bool
	CountBytes bool
}

var DefaultOptions = Options{
	CountLines: true,
	CountWords: true,
	CountBytes: true,
}

func NewOptions(lines, words, chars, bytes bool) Options {
	if !lines && !words && !chars && !bytes {
		return DefaultOptions
	}
	return Options{
		CountLines: lines,
		CountWords: words,
		CountChars: chars,
		CountBytes: bytes,
	}
}

func (opts Options) SelectCountingFn() CountingFn {
	if opts.CountLines && !opts.CountWords && !opts.CountChars && !opts.CountBytes {
		return CountLines
	} else if opts.CountWords && !opts.CountLines && !opts.CountChars && !opts.CountBytes {
		return CountWords
	} else if opts.CountChars && !opts.CountLines && !opts.CountWords && !opts.CountBytes {
		return CountChars
	} else if opts.CountBytes && !opts.CountLines && !opts.CountWords && !opts.CountChars {
		return CountBytes
	} else if opts.CountLines && opts.CountBytes && !opts.CountWords && !opts.CountChars {
		return CountLinesBytes
	} else if !opts.CountChars {
		return CountLinesWordsBytes
	} else {
		return CountLinesWordsCharsBytes
	}
}

type Stats struct {
	lines uint
	words uint
	chars uint
	bytes uint
}

func (s Stats) Print(opts Options, name string) {
	if opts.CountLines {
		fmt.Printf("%d\t", s.lines)
	}
	if opts.CountWords {
		fmt.Printf("%d\t", s.words)
	}
	if opts.CountChars {
		fmt.Printf("%d\t", s.chars)
	}
	if opts.CountBytes {
		fmt.Printf("%d\t", s.bytes)
	}
	if name != "" {
		fmt.Printf("%s\n", name)
	} else {
		fmt.Printf("\n")
	}
}

type CountingFn func(io.Reader) (Stats, error)

func CountLines(r io.Reader) (Stats, error) {
	count := 0
	buffer := make([]byte, bufferSize)
	for {
		n, err := r.Read(buffer)
		count += bytes.Count(buffer[:n], []byte("\n"))
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return Stats{}, err
		}
	}
	return Stats{lines: uint(count)}, nil
}

func CountWords(r io.Reader) (Stats, error) {
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
			return Stats{}, err
		}
	}
	// EOF counts as WS
	if !inWhitespace {
		count++
	}
	return Stats{words: uint(count)}, nil
}

func CountChars(r io.Reader) (Stats, error) {
	count := 0
	buffer := make([]byte, bufferSize)
	writeStart := 0
	for {
		n, err := r.Read(buffer[writeStart:])
		b := buffer[:n]

		// Count runes in buffer[:n]
		for {
			if len(b) == 0 {
				writeStart = 0
				break
			}
			rune, size := utf8.DecodeRune(b)
			if rune == utf8.RuneError {
				if len(b) > 4 {
					// Cannot be an incomplete rune, discard first byte
					b = b[1:]
					continue
				}
				// Last bytes, let's read more and see if it gets fixed
				copy(buffer, b)
				writeStart = len(b)
				break
			}
			count++
			b = b[size:]
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return Stats{}, err
		}
	}
	return Stats{chars: uint(count)}, nil
}

func CountBytes(r io.Reader) (Stats, error) {
	count := 0
	buffer := make([]byte, bufferSize)
	for {
		n, err := r.Read(buffer)
		count += n
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return Stats{}, err
		}
	}
	return Stats{bytes: uint(count)}, nil
}

func CountLinesBytes(r io.Reader) (Stats, error) {
	var (
		linesN int
		bytesN int
	)
	buffer := make([]byte, bufferSize)
	for {
		n, err := r.Read(buffer)
		bytesN += n
		linesN += bytes.Count(buffer[:n], []byte("\n"))
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return Stats{}, err
		}
	}
	return Stats{lines: uint(linesN), bytes: uint(bytesN)}, nil
}

func CountLinesWordsBytes(r io.Reader) (Stats, error) {
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
			return Stats{}, err
		}
	}
	// EOF counts as WS
	if !inWhitespace {
		words++
	}
	return Stats{
		lines: uint(lines),
		words: uint(words),
		bytes: uint(bytes),
	}, nil
}

func CountLinesWordsCharsBytes(r io.Reader) (Stats, error) {
	var (
		lines        int
		words        int
		chars        int
		bytes        int
		inWhitespace = true
	)
	buffer := make([]byte, bufferSize)
	writeStart := 0
	for {
		n, err := r.Read(buffer[writeStart:])
		bytes += n

		// Count runes in buffer[:n]
		b := buffer[:n]
		for {
			if len(b) == 0 {
				writeStart = 0
				break
			}
			rune, size := utf8.DecodeRune(b)
			if rune == utf8.RuneError {
				if len(b) > 4 {
					// Cannot be an incomplete rune, discard first byte
					b = b[1:]
					continue
				}
				// Last bytes, let's read more and see if it gets fixed
				copy(buffer, b)
				writeStart = len(b)
				break
			}
			chars++
			b = b[size:]
		}

		// Count lines, words and bytes
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
			return Stats{}, err
		}
	}
	// EOF counts as WS
	if !inWhitespace {
		words++
	}
	return Stats{
		lines: uint(lines),
		words: uint(words),
		chars: uint(chars),
		bytes: uint(bytes),
	}, nil
}
