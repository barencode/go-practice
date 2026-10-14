package main

import (
	"bufio"
	"bytes"
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

func TestExecuteCommandClose(t *testing.T) {
	mine := Mine{Closed: false}
	shouldExit := executeCommand(io.Discard, &mine, Command{Name: "close"})
	if !mine.Closed || shouldExit {
		t.Errorf("got closed=%t, should exit=%t, want: closed=true and should exit=false", mine.Closed, shouldExit)
	}
}

func TestExecuteCommandOpen(t *testing.T) {
	mine := Mine{Closed: true}
	shouldExit := executeCommand(io.Discard, &mine, Command{Name: "open"})
	if mine.Closed || shouldExit {
		t.Errorf("got closed=%t, should exit=%t, want: closed=false and should exit=false", mine.Closed, shouldExit)
	}
}

func TestExecuteCommandExit(t *testing.T) {
	mine := Mine{}
	shouldExit := executeCommand(io.Discard, &mine, Command{Name: "exit"})
	if !shouldExit {
		t.Errorf("got should exit=%t, want: should exit=true", shouldExit)
	}
}

func TestExecuteCommandMine(t *testing.T) {
	mine := Mine{Reserves: 5}
	mine.AddMiner(101, &Miner{Name: "Борис", Energy: 1, Production: 3})
	shouldExit := executeCommand(io.Discard, &mine, Command{Name: "mine", MinerID: 101})
	if shouldExit || mine.Reserves != 2 || mine.Coal != 3 || mine.Miners[101].Energy != 0 {
		t.Errorf("got should exit=%t, reserves=%d, coal=%d, and 1st miner energy=%d; want should exit=false, reserves=2, coal=3 and energy=0", shouldExit, mine.Reserves, mine.Coal, mine.Miners[101].Energy)
	}
}

func TestExecuteCommandRest(t *testing.T) {
	mine := Mine{Coal: 3, Reserves: 5}
	mine.AddMiner(101, &Miner{Name: "Борис", Energy: 0})
	shouldExit := executeCommand(io.Discard, &mine, Command{Name: "rest", MinerID: 101})
	if shouldExit || mine.Coal != 1 || mine.Reserves != 5 || mine.Miners[101].Energy != 1 {
		t.Errorf("got should exit=%t, reserves=%d, coal=%d, and 1st miner energy=%d; want should exit=false, reserves=5, coal=1 and energy=1", shouldExit, mine.Reserves, mine.Coal, mine.Miners[101].Energy)
	}
}

func TestExecuteCommandMineClosed(t *testing.T) {
	mine := Mine{Closed: true, Coal: 4, Reserves: 5}
	mine.AddMiner(101, &Miner{Name: "Борис", Energy: 1, Production: 3})
	var out bytes.Buffer
	shouldExit := executeCommand(&out, &mine, Command{Name: "mine", MinerID: 101})
	if shouldExit || mine.Coal != 4 || mine.Reserves != 5 || mine.Miners[101].Energy != 1 {
		t.Errorf("got should exit=%t, reserves=%d, coal=%d, and 1st miner energy=%d; want should exit=false, reserves=5, coal=4 and energy=1", shouldExit, mine.Reserves, mine.Coal, mine.Miners[101].Energy)
	}
	want := "Ищем шахтёра с ID: 101\nДобыча остановлена\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestPrintMineSummaryEmpty(t *testing.T) {
	mine := Mine{Coal: 4, Reserves: 5, Closed: true}
	var out bytes.Buffer

	printMineSummary(&out, &mine)
	got := out.String()
	want := "На складе: 4, в залежах: 5, шахта закрыта: true\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPrintMineSummaryWithMiners(t *testing.T) {
	mine := Mine{}
	mine.AddMiner(205, nil)
	mine.AddMiner(101, &Miner{Name: "Борис", Energy: 2, Production: 3})
	var out bytes.Buffer

	printMineSummary(&out, &mine)
	got := out.String()
	want := "На складе: 0, в залежах: 0, шахта закрыта: false\n" +
		"101: Борис: энергия 2, добыча за смену 3\n" +
		"205: miner 205: шахтёр не найден\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRunConsoleMineAndRest(t *testing.T) {
	mine := Mine{Reserves: 5}
	mine.AddMiner(101, &Miner{Name: "Борис", Energy: 1, Production: 3})

	in := strings.NewReader("mine 101\nrest 101\nexit\nmine 101\n")
	var out bytes.Buffer

	runConsole(in, &out, &mine)

	if mine.Coal != 1 || mine.Reserves != 2 || mine.Miners[101].Energy != 1 {
		t.Errorf("got coal=%d, reserves=%d and 1st miner energy=%d; want coal=1, reserves=2 and energy=1", mine.Coal, mine.Reserves, mine.Miners[101].Energy)
	}
}
