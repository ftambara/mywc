package main

import (
	"bytes"
	"slices"
	"testing"
)

func assertSlice[C comparable](t *testing.T, expected []C, got []C) {
	if !slices.Equal(expected, got) {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestParsingCountModes(t *testing.T) {
	conf, err := parseArgs([]string{"-chars"})
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertSlice(t, []countMode{byChars}, conf.countModes)

	conf, err = parseArgs([]string{"-chars", "-words"})
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertSlice(t, []countMode{byWords, byChars}, conf.countModes)

	conf, err = parseArgs([]string{})
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertSlice(t, []countMode{byLines, byWords, byBytes}, conf.countModes)

	conf, err = parseArgs([]string{"-lines", "-lines"})
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertSlice(t, []countMode{byLines}, conf.countModes)

	conf, err = parseArgs([]string{"file"})
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertSlice(t, []countMode{byLines, byWords, byBytes}, conf.countModes)
}

func TestParsingModesAndArgs(t *testing.T) {
	conf, err := parseArgs([]string{"-chars", "file"})
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertSlice(t, []countMode{byChars}, conf.countModes)
	assertSlice(t, []string{"file"}, conf.files)

	conf, err = parseArgs([]string{"-chars", "-words", "file1", "file2"})
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertSlice(t, []countMode{byWords, byChars}, conf.countModes)
	assertSlice(t, []string{"file1", "file2"}, conf.files)

	conf, err = parseArgs([]string{"-bytes"})
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertSlice(t, []countMode{byBytes}, conf.countModes)
	assertSlice(t, []string{}, conf.files)
}

func TestCountLines(t *testing.T) {
	r := bytes.NewReader([]byte(""))
	lines, err := countLines(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if lines != 0 {
		t.Error("expected 0 lines, got", lines)
	}

	r = bytes.NewReader([]byte("a\n"))
	lines, err = countLines(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if lines != 1 {
		t.Error("expected 1 line, got", lines)
	}

	r = bytes.NewReader([]byte("ab\ncd"))
	lines, err = countLines(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if lines != 1 {
		t.Error("expected 1 lines, got", lines)
	}

	r = bytes.NewReader([]byte("ab\n\ncd\n"))
	lines, err = countLines(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if lines != 3 {
		t.Error("expected 3 lines, got", lines)
	}
}

func TestCountBytes(t *testing.T) {
	r := bytes.NewReader([]byte{})
	bytesN, err := countBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if bytesN != 0 {
		t.Error("expected 0 bytes, got", bytesN)
	}

	r = bytes.NewReader([]byte{1, 100, 0})
	bytesN, err = countBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if bytesN != 3 {
		t.Error("expected 3 bytes, got", bytesN)
	}

	r = bytes.NewReader([]byte("abc"))
	bytesN, err = countBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if bytesN != 3 {
		t.Error("expected 3 bytes, got", bytesN)
	}

	r = bytes.NewReader([]byte("Á¥\næ"))
	bytesN, err = countBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if bytesN != 7 {
		t.Error("expected 7 bytes, got", bytesN)
	}
}

func TestCountWords(t *testing.T) {
	r := bytes.NewReader([]byte(""))
	words, err := countWords(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if words != 0 {
		t.Error("expected 0 words, got", words)
	}

	r = bytes.NewReader([]byte("abc"))
	words, err = countWords(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if words != 1 {
		t.Error("expected 1 words, got", words)
	}

	r = bytes.NewReader([]byte("a b\nc"))
	words, err = countWords(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if words != 3 {
		t.Error("expected 3 words, got", words)
	}

	r = bytes.NewReader([]byte("a \n b\nc"))
	words, err = countWords(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if words != 3 {
		t.Error("expected 3 words, got", words)
	}

	r = bytes.NewReader([]byte("a b\n"))
	words, err = countWords(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if words != 2 {
		t.Error("expected 2 words, got", words)
	}
}

func TestCountLinesBytes(t *testing.T) {
	r := bytes.NewReader([]byte(""))
	counts, err := countLinesBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if counts != [2]uint{0, 0} {
		t.Error("expected [0 0], got", counts)
	}

	r = bytes.NewReader([]byte("a\n"))
	counts, err = countLinesBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if counts != [2]uint{1, 2} {
		t.Error("expected [1 2], got", counts)
	}

	r = bytes.NewReader([]byte("ab\ncd"))
	counts, err = countLinesBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if counts != [2]uint{1, 5} {
		t.Error("expected [1 5], got", counts)
	}

	r = bytes.NewReader([]byte("ab\n\ncd\n"))
	counts, err = countLinesBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if counts != [2]uint{3, 7} {
		t.Error("expected [3 7], got", counts)
	}
}

func TestCountLinesWordsBytes(t *testing.T) {
	r := bytes.NewReader([]byte(""))
	counts, err := countLinesWordsBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if counts != [3]uint{0, 0, 0} {
		t.Error("expected [0 0 0], got", counts)
	}

	r = bytes.NewReader([]byte("a\n"))
	counts, err = countLinesWordsBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if counts != [3]uint{1, 1, 2} {
		t.Error("expected [1 1 2], got", counts)
	}

	r = bytes.NewReader([]byte("ab\ncd"))
	counts, err = countLinesWordsBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if counts != [3]uint{1, 2, 5} {
		t.Error("expected [1 2 5], got", counts)
	}

	r = bytes.NewReader([]byte("ab\n\ncd\n"))
	counts, err = countLinesWordsBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if counts != [3]uint{3, 2, 7} {
		t.Error("expected [3 2 7], got", counts)
	}
}
