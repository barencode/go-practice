package main

import (
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
