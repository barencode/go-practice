package main

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestParseMinerIDWithSpaces(t *testing.T) {
	str := " 101 "
	got, err := parseMinerID(str)

	if err != nil || got != 101 {
		t.Errorf("got %d and err=%v; want 101 and nil", got, err)
	}
}

func TestParseMinerIDInvalid(t *testing.T) {
	testSequence := []string{"Борис", "", "  ", "101abc"}
	for _, input := range testSequence {
		_, err := parseMinerID(input)

		if err == nil {
			t.Errorf("input=%q: got err=nil; want non-nil error", input)
		}
	}
}

func TestParseMinerIDZero(t *testing.T) {
	input := "0"
	got, err := parseMinerID(input)

	if got != 0 || err != nil {
		t.Errorf("got %d and err=%v; want 0 and nil", got, err)
	}
}

func TestReadLine(t *testing.T) {
	scanner := bufio.NewScanner(strings.NewReader("Борис\n"))

	got, err := readLine(scanner)
	if got != "Борис" || err != nil {
		t.Errorf("got text=%q, err=%v, want text=%q, err=nil", got, err, "Борис")
	}

	got, err = readLine(scanner)
	if got != "" || !errors.Is(err, io.EOF) {
		t.Errorf("got text=%q, err=%v, want text=%q, err=EOF", got, err, "")
	}
}

func TestReadLineEmpty(t *testing.T) {
	scanner := bufio.NewScanner(strings.NewReader("\n"))

	got, err := readLine(scanner)
	if got != "" || err != nil {
		t.Errorf("got text=%q, err=%v, want text=%q, err=nil", got, err, "")
	}
}
