package counters

import (
	"bytes"
	"errors"
	"io"
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
