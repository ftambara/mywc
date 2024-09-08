package counters

import (
	"bytes"
	"errors"
	"io"
	"unicode/utf8"
)

const bufferSize = 8192

func CountLines(r io.Reader) (uint, error) {
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

func CountWords(r io.Reader) (uint, error) {
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

func CountChars(r io.Reader) (uint, error) {
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
			return uint(count), err
		}
	}
	return uint(count), nil
}

func CountBytes(r io.Reader) (uint, error) {
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

func CountLinesBytes(r io.Reader) ([2]uint, error) {
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
			return [2]uint{}, err
		}
	}
	return [...]uint{uint(linesN), uint(bytesN)}, nil
}

func CountLinesWordsBytes(r io.Reader) ([3]uint, error) {
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

func CountLinesWordsCharsBytes(r io.Reader) ([4]uint, error) {
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
			return [4]uint{}, err
		}
	}
	// EOF counts as WS
	if !inWhitespace {
		words++
	}
	return [...]uint{uint(lines), uint(words), uint(chars), uint(bytes)}, nil
}
