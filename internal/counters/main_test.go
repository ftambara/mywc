package counters

import (
	"bytes"
	"testing"
)

func TestCountLines(t *testing.T) {
	r := bytes.NewReader([]byte(""))
	lines, err := CountLines(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if lines != 0 {
		t.Error("expected 0 lines, got", lines)
	}

	r = bytes.NewReader([]byte("a\n"))
	lines, err = CountLines(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if lines != 1 {
		t.Error("expected 1 line, got", lines)
	}

	r = bytes.NewReader([]byte("ab\ncd"))
	lines, err = CountLines(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if lines != 1 {
		t.Error("expected 1 lines, got", lines)
	}

	r = bytes.NewReader([]byte("ab\n\ncd\n"))
	lines, err = CountLines(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if lines != 3 {
		t.Error("expected 3 lines, got", lines)
	}
}

func TestCountWords(t *testing.T) {
	r := bytes.NewReader([]byte(""))
	words, err := CountWords(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if words != 0 {
		t.Error("expected 0 words, got", words)
	}

	r = bytes.NewReader([]byte("abc"))
	words, err = CountWords(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if words != 1 {
		t.Error("expected 1 words, got", words)
	}

	r = bytes.NewReader([]byte("a b\nc"))
	words, err = CountWords(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if words != 3 {
		t.Error("expected 3 words, got", words)
	}

	r = bytes.NewReader([]byte("a \n b\nc"))
	words, err = CountWords(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if words != 3 {
		t.Error("expected 3 words, got", words)
	}

	r = bytes.NewReader([]byte("a b\n"))
	words, err = CountWords(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if words != 2 {
		t.Error("expected 2 words, got", words)
	}
}

func TestCountChars(t *testing.T) {
	r := bytes.NewReader([]byte(""))
	chars, err := CountChars(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if chars != 0 {
		t.Error("expected 0 chars, got", chars)
	}

	r = bytes.NewReader([]byte("abc"))
	chars, err = CountChars(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if chars != 3 {
		t.Error("expected 3 chars, got", chars)
	}

	r = bytes.NewReader([]byte("Hello, 世界"))
	chars, err = CountChars(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if chars != 9 {
		t.Error("expected 9 chars, got", chars)
	}
}

func TestCountBytes(t *testing.T) {
	r := bytes.NewReader([]byte{})
	bytesN, err := CountBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if bytesN != 0 {
		t.Error("expected 0 bytes, got", bytesN)
	}

	r = bytes.NewReader([]byte{1, 100, 0})
	bytesN, err = CountBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if bytesN != 3 {
		t.Error("expected 3 bytes, got", bytesN)
	}

	r = bytes.NewReader([]byte("abc"))
	bytesN, err = CountBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if bytesN != 3 {
		t.Error("expected 3 bytes, got", bytesN)
	}

	r = bytes.NewReader([]byte("Á¥\næ"))
	bytesN, err = CountBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if bytesN != 7 {
		t.Error("expected 7 bytes, got", bytesN)
	}
}

func TestCountLinesBytes(t *testing.T) {
	r := bytes.NewReader([]byte(""))
	counts, err := CountLinesBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if counts != [2]uint{0, 0} {
		t.Error("expected [0 0], got", counts)
	}

	r = bytes.NewReader([]byte("a\n"))
	counts, err = CountLinesBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if counts != [2]uint{1, 2} {
		t.Error("expected [1 2], got", counts)
	}

	r = bytes.NewReader([]byte("ab\ncd"))
	counts, err = CountLinesBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if counts != [2]uint{1, 5} {
		t.Error("expected [1 5], got", counts)
	}

	r = bytes.NewReader([]byte("ab\n\ncd\n"))
	counts, err = CountLinesBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if counts != [2]uint{3, 7} {
		t.Error("expected [3 7], got", counts)
	}
}

func TestCountLinesWordsBytes(t *testing.T) {
	r := bytes.NewReader([]byte(""))
	counts, err := CountLinesWordsBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if counts != [3]uint{0, 0, 0} {
		t.Error("expected [0 0 0], got", counts)
	}

	r = bytes.NewReader([]byte("a\n"))
	counts, err = CountLinesWordsBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if counts != [3]uint{1, 1, 2} {
		t.Error("expected [1 1 2], got", counts)
	}

	r = bytes.NewReader([]byte("ab\ncd"))
	counts, err = CountLinesWordsBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if counts != [3]uint{1, 2, 5} {
		t.Error("expected [1 2 5], got", counts)
	}

	r = bytes.NewReader([]byte("ab\n\ncd\n"))
	counts, err = CountLinesWordsBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if counts != [3]uint{3, 2, 7} {
		t.Error("expected [3 2 7], got", counts)
	}
}
