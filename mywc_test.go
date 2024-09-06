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

func Test_parseModes(t *testing.T) {
	countModes, err := parseModes([]string{"-chars"})
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertSlice(t, []countMode{chars}, countModes)

	countModes, err = parseModes([]string{"-chars", "-words"})
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertSlice(t, []countMode{words, chars}, countModes)

	countModes, err = parseModes([]string{})
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertSlice(t, []countMode{lines, words, bytes}, countModes)

	countModes, err = parseModes([]string{"-lines", "-lines"})
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertSlice(t, []countMode{lines}, countModes)

	_, err = parseModes([]string{"wrongArg"})
	if err == nil {
		t.Error("should return a non-nil error")
	}
}
