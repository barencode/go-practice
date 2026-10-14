package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Command struct {
	Name    string
	MinerID int
}

func printMineSummary(out io.Writer, mine *Mine) {
	fmt.Fprintf(out, "На складе: %d, в залежах: %d, шахта закрыта: %t\n", mine.Coal, mine.Reserves, mine.Closed)

	ids := mine.MinerIDs()
	for _, id := range ids {
		miner, err := mine.FindMiner(id)
		if err != nil {
			fmt.Fprintf(out, "%d: %v\n", id, err)
			continue
		}
		fmt.Fprintf(out, "%d: %s: энергия %d, добыча за смену %d\n", id, miner.Name, miner.Energy, miner.Production)
	}
}

func parseMinerID(input string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(input))
}

func readLine(scanner *bufio.Scanner) (string, error) {
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}
	return scanner.Text(), nil
}

func parseCommand(input string) (Command, error) {
	parts := strings.Fields(input)

	if len(parts) == 0 {
		return Command{}, ErrEmptyCommand
	}

	switch parts[0] {
	case "status", "exit", "close", "open":
		if len(parts) != 1 {
			return Command{}, ErrInvalidCommand
		}
		return Command{Name: parts[0]}, nil
	case "mine", "rest":
		if len(parts) != 2 {
			return Command{}, ErrInvalidCommand
		}
		id, err := parseMinerID(parts[1])
		if err != nil {
			return Command{}, err
		}
		return Command{Name: parts[0], MinerID: id}, nil
	default:
		return Command{}, ErrInvalidCommand
	}
}

func executeCommand(out io.Writer, mine *Mine, command Command) bool {
	switch command.Name {
	case "mine":
		id := command.MinerID
		fmt.Fprintf(out, "Ищем шахтёра с ID: %d\n", id)
		coal, err := mine.MineByID(id)
		if err != nil {
			if errors.Is(err, ErrMineClosed) {
				fmt.Fprintf(out, "Добыча остановлена\n")
				return false
			}
			if errors.Is(err, ErrMinerNotFound) {
				fmt.Fprintf(out, "Выберите другого шахтёра\n")
				return false
			}
			fmt.Fprintf(out, "Operation caused an error: %v\n", err)
			return false
		}
		fmt.Fprintf(out, "Coal mined: %d\n", coal)
		printMineSummary(out, mine)
	case "rest":
		id := command.MinerID
		fmt.Fprintf(out, "Ищем шахтёра с ID: %d\n", id)
		err := mine.RestMiner(id)
		if err != nil {
			if errors.Is(err, ErrMinerNotFound) {
				fmt.Fprintf(out, "Выберите другого шахтёра\n")
				return false
			}
			if errors.Is(err, ErrNotEnoughCoal) {
				fmt.Fprintf(out, "Недостаточно угля\n")
				return false
			}
			fmt.Fprintf(out, "Operation caused an error: %v\n", err)
			return false
		}
		fmt.Fprintf(out, "Шахтёр с ID: %d отдохнул\n", id)
		printMineSummary(out, mine)
	case "status":
		printMineSummary(out, mine)
	case "open":
		mine.Closed = false
		fmt.Fprintf(out, "Шахта открыта\n")
	case "close":
		mine.Closed = true
		fmt.Fprintf(out, "Шахта закрыта\n")
	case "exit":
		return true
	}
	return false
}

func runConsole(in io.Reader, out io.Writer, mine *Mine) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, "Введите команду (mine <ID>, rest <ID>, status, exit, open, close): ")
		input, err := readLine(scanner)
		if errors.Is(err, io.EOF) {
			return
		}
		if err != nil {
			fmt.Fprintln(out, "Ошибка чтения:", err)
			return
		}

		command, err := parseCommand(input)
		if errors.Is(err, ErrEmptyCommand) {
			continue
		}
		if err != nil {
			fmt.Fprintln(out, "Ошибка команды:", err)
			continue
		}

		if executeCommand(out, mine, command) {
			return
		}
	}
}
