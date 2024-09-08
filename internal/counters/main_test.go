package counters

import (
	"bytes"
	"testing"
)

func assertStats(t *testing.T, expected, actual Stats) {
	if expected.lines != actual.lines {
		t.Errorf("expected %d lines, got %d", expected.lines, actual.lines)
	}
	if expected.words != actual.words {
		t.Errorf("expected %d words, got %d", expected.words, actual.words)
	}
	if expected.chars != actual.chars {
		t.Errorf("expected %d chars, got %d", expected.chars, actual.chars)
	}
	if expected.bytes != actual.bytes {
		t.Errorf("expected %d bytes, got %d", expected.bytes, actual.bytes)
	}
}

func TestCountLines(t *testing.T) {
	r := bytes.NewReader([]byte(""))
	stats, err := CountLines(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{lines: 0}, stats)

	r = bytes.NewReader([]byte("a\n"))
	stats, err = CountLines(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{lines: 1}, stats)

	r = bytes.NewReader([]byte("ab\ncd"))
	stats, err = CountLines(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{lines: 1}, stats)

	r = bytes.NewReader([]byte("ab\n\ncd\n"))
	stats, err = CountLines(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{lines: 3}, stats)
}

func TestCountWords(t *testing.T) {
	r := bytes.NewReader([]byte(""))
	stats, err := CountWords(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{words: 0}, stats)

	r = bytes.NewReader([]byte("abc"))
	stats, err = CountWords(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{words: 1}, stats)

	r = bytes.NewReader([]byte("a b\nc"))
	stats, err = CountWords(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{words: 3}, stats)

	r = bytes.NewReader([]byte("a \n b\nc"))
	stats, err = CountWords(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{words: 3}, stats)

	r = bytes.NewReader([]byte("a b\n"))
	stats, err = CountWords(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{words: 2}, stats)
}

func TestCountChars(t *testing.T) {
	r := bytes.NewReader([]byte(""))
	stats, err := CountChars(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{chars: 0}, stats)

	r = bytes.NewReader([]byte("abc"))
	stats, err = CountChars(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{chars: 3}, stats)

	r = bytes.NewReader([]byte("Hello, 世界"))
	stats, err = CountChars(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{chars: 9}, stats)
}

func TestCountBytes(t *testing.T) {
	r := bytes.NewReader([]byte{})
	stats, err := CountBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{bytes: 0}, stats)

	r = bytes.NewReader([]byte{1, 100, 0})
	stats, err = CountBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{bytes: 3}, stats)

	r = bytes.NewReader([]byte("abc"))
	stats, err = CountBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{bytes: 3}, stats)

	r = bytes.NewReader([]byte("Á¥\næ"))
	stats, err = CountBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{bytes: 7}, stats)
}

func TestCountLinesBytes(t *testing.T) {
	r := bytes.NewReader([]byte(""))
	stats, err := CountLinesBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{lines: 0, bytes: 0}, stats)

	r = bytes.NewReader([]byte("a\n"))
	stats, err = CountLinesBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{lines: 1, bytes: 2}, stats)

	r = bytes.NewReader([]byte("ab\ncd"))
	stats, err = CountLinesBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{lines: 1, bytes: 5}, stats)

	r = bytes.NewReader([]byte("ab\n\ncd\n"))
	stats, err = CountLinesBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{lines: 3, bytes: 7}, stats)
}

func TestCountLinesWordsBytes(t *testing.T) {
	r := bytes.NewReader([]byte(""))
	stats, err := CountLinesWordsBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{lines: 0, words: 0, bytes: 0}, stats)

	r = bytes.NewReader([]byte("a\n"))
	stats, err = CountLinesWordsBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{lines: 1, words: 1, bytes: 2}, stats)

	r = bytes.NewReader([]byte("ab\ncd"))
	stats, err = CountLinesWordsBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{lines: 1, words: 2, bytes: 5}, stats)

	r = bytes.NewReader([]byte("ab\n\ncd\n"))
	stats, err = CountLinesWordsBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{lines: 3, words: 2, bytes: 7}, stats)
}

func TestCountLinesWordsCharsBytes(t *testing.T) {
	r := bytes.NewReader([]byte(""))
	stats, err := CountLinesWordsCharsBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{lines: 0, words: 0, chars: 0, bytes: 0}, stats)

	r = bytes.NewReader([]byte("a\n"))
	stats, err = CountLinesWordsCharsBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{lines: 1, words: 1, chars: 2, bytes: 2}, stats)

	r = bytes.NewReader([]byte("世界"))
	stats, err = CountLinesWordsCharsBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{lines: 0, words: 1, chars: 2, bytes: 6}, stats)

	r = bytes.NewReader([]byte("ab\n\n测试\n"))
	stats, err = CountLinesWordsCharsBytes(r)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	assertStats(t, Stats{lines: 3, words: 2, chars: 7, bytes: 11}, stats)
}
