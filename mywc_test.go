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
	assertSlice(t, []countMode{byChars}, conf.countModes)

	conf, err = parseArgs([]string{"-chars", "-words"})
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertSlice(t, []countMode{byWords, byChars}, conf.countModes)

	conf, err = parseArgs([]string{})
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertSlice(t, []countMode{byLines, byWords, byBytes}, conf.countModes)

	conf, err = parseArgs([]string{"-lines", "-lines"})
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertSlice(t, []countMode{byLines}, conf.countModes)

	conf, err = parseArgs([]string{"file"})
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertSlice(t, []countMode{byLines, byWords, byBytes}, conf.countModes)
}

func TestParsingModesAndArgs(t *testing.T) {
	conf, err := parseArgs([]string{"-chars", "file"})
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertSlice(t, []countMode{byChars}, conf.countModes)
	assertSlice(t, []string{"file"}, conf.files)

	conf, err = parseArgs([]string{"-chars", "-words", "file1", "file2"})
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertSlice(t, []countMode{byWords, byChars}, conf.countModes)
	assertSlice(t, []string{"file1", "file2"}, conf.files)

	conf, err = parseArgs([]string{"-bytes"})
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertSlice(t, []countMode{byBytes}, conf.countModes)
	assertSlice(t, []string{}, conf.files)
}
