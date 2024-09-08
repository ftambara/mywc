package main

import (
	"slices"
	"testing"

	"github.com/ftambara/mywc/internal/counters"
)

func assertSlice[C comparable](t *testing.T, expected []C, got []C) {
	if !slices.Equal(expected, got) {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func assertOpts(t *testing.T, expected counters.Options, got counters.Options) {
	if expected.CountLines != got.CountLines {
		t.Errorf("expected CountLines %v, got %v", expected.CountLines, got.CountLines)
	}
	if expected.CountWords != got.CountWords {
		t.Errorf("expected CountWords %v, got %v", expected.CountWords, got.CountWords)
	}
	if expected.CountChars != got.CountChars {
		t.Errorf("expected CountChars %v, got %v", expected.CountChars, got.CountChars)
	}
	if expected.CountBytes != got.CountBytes {
		t.Errorf("expected CountBytes %v, got %v", expected.CountBytes, got.CountBytes)
	}
}

func assertConfig(t *testing.T, expected config, got config) {
	assertSlice(t, expected.files, got.files)
	assertOpts(t, expected.opts, got.opts)
}

func TestParsingOpts(t *testing.T) {
	conf, err := parseArgs([]string{"-chars"})
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertOpts(t, counters.Options{CountChars: true}, conf.opts)

	conf, err = parseArgs([]string{"-chars", "-words"})
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertOpts(t, counters.Options{CountChars: true, CountWords: true}, conf.opts)

	conf, err = parseArgs([]string{})
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertOpts(t, counters.DefaultOptions, conf.opts)

	conf, err = parseArgs([]string{"-lines", "-lines"})
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertOpts(t, counters.Options{CountLines: true}, conf.opts)

	conf, err = parseArgs([]string{"file"})
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertOpts(t, counters.DefaultOptions, conf.opts)
}

func TestParsingArgs(t *testing.T) {
	conf, err := parseArgs([]string{"-chars", "file"})
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	expected := config{
		files: []string{"file"},
		opts:  counters.Options{CountChars: true},
	}
	assertConfig(t, expected, conf)

	conf, err = parseArgs([]string{"-chars", "-words", "file1", "file2"})
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	expected = config{
		files: []string{"file1", "file2"},
		opts:  counters.Options{CountChars: true, CountWords: true},
	}
	assertConfig(t, expected, conf)

	conf, err = parseArgs([]string{"-bytes"})
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	expected = config{
		files: []string{},
		opts:  counters.Options{CountBytes: true},
	}
	assertConfig(t, expected, conf)
}
