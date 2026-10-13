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

func TestParseCommandMine(t *testing.T) {
	got, err := parseCommand("  mine  101  ")
	if got.Name != "mine" || got.MinerID != 101 || err != nil {
		t.Errorf("got command=%+v, err=%v; want Name=mine, MinerID=101, err=nil", got, err)
	}
}

func TestParseCommandStatus(t *testing.T) {
	got, err := parseCommand("status")
	if got.Name != "status" || got.MinerID != 0 || err != nil {
		t.Errorf("got command=%+v, err=%v; want Name=status, MinerID=0, err=nil", got, err)
	}
}

func TestParseCommandEmpty(t *testing.T) {
	inputs := []string{"", "  "}
	for _, input := range inputs {
		got, err := parseCommand(input)
		if !errors.Is(err, ErrEmptyCommand) || got.Name != "" || got.MinerID != 0 {
			t.Errorf("input=%q: got command=%+v, err=%v; want empty command, err=%v", input, got, err, ErrEmptyCommand)
		}
	}
}

func TestParseCommandInvalid(t *testing.T) {
	inputs := []string{"hello", "mine", "rest", "mine 101 extra", "status 101", "exit extra", "open 101", "close 101"}
	for _, input := range inputs {
		got, err := parseCommand(input)
		if !errors.Is(err, ErrInvalidCommand) || got.Name != "" || got.MinerID != 0 {
			t.Errorf("input=%q: got command=%+v, err=%v; want empty command, err=%v", input, got, err, ErrInvalidCommand)
		}
	}
}

func TestParseCommandInvalidID(t *testing.T) {
	got, err := parseCommand("mine Борис")
	if err == nil || got.Name != "" || got.MinerID != 0 {
		t.Errorf("got command=%+v, err=%v; want empty command and non-nil error", got, err)
	}
}

func TestParseCommandRest(t *testing.T) {
	got, err := parseCommand("rest 101")
	if got.Name != "rest" || got.MinerID != 101 || err != nil {
		t.Errorf("got command=%+v, command=%s, minerID=%d, err=%v; want command=\"rest\", ID=101 and nil error", got, got.Name, got.MinerID, err)
	}
}

func TestParseCommandExit(t *testing.T) {
	got, err := parseCommand("exit")
	if got.Name != "exit" || got.MinerID != 0 || err != nil {
		t.Errorf("got command=%+v, command=%s, minerID=%d, err=%v; want command=\"exit\", ID=0 and nil error", got, got.Name, got.MinerID, err)
	}
}

func TestParseCommandOpenClose(t *testing.T) {
	inputs := []string{"open", "close"}
	for _, input := range inputs {
		got, err := parseCommand(input)
		if got.Name != input || got.MinerID != 0 || err != nil {
			t.Errorf("input=%q: got command=%+v, err=%v; want Name=%q, MinerID=0, err=nil", input, got, err, input)
		}
	}
}
