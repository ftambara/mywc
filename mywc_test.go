package main

import (
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
		t.Error("unexpected error", err)
	}
	assertSlice(t, []countMode{chars}, conf.countModes)

	conf, err = parseArgs([]string{"-chars", "-words"})
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertSlice(t, []countMode{words, chars}, conf.countModes)

	conf, err = parseArgs([]string{})
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertSlice(t, []countMode{lines, words, bytes}, conf.countModes)

	conf, err = parseArgs([]string{"-lines", "-lines"})
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertSlice(t, []countMode{lines}, conf.countModes)

	conf, err = parseArgs([]string{"file"})
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertSlice(t, []countMode{lines, words, bytes}, conf.countModes)
}

func TestParsingModesAndArgs(t *testing.T) {
	conf, err := parseArgs([]string{"-chars", "file"})
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertSlice(t, []countMode{chars}, conf.countModes)
	assertSlice(t, []string{"file"}, conf.files)

	conf, err = parseArgs([]string{"-chars", "-words", "file1", "file2"})
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertSlice(t, []countMode{words, chars}, conf.countModes)
	assertSlice(t, []string{"file1", "file2"}, conf.files)

	conf, err = parseArgs([]string{"-bytes"})
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertSlice(t, []countMode{bytes}, conf.countModes)
	assertSlice(t, []string{}, conf.files)
}
