package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
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

func executeCommand(mine *Mine, command Command) bool {
	switch command.Name {
	case "mine":
		id := command.MinerID
		fmt.Println("Ищем шахтёра с ID:", id)
		coal, err := mine.MineByID(id)
		if err != nil {
			if errors.Is(err, ErrMineClosed) {
				fmt.Println("Добыча остановлена")
				return false
			}
			if errors.Is(err, ErrMinerNotFound) {
				fmt.Println("Выберите другого шахтёра")
				return false
			}
			fmt.Println("Operation caused an error:", err)
			return false
		}
		fmt.Println("Coal mined:", coal)
		printMineSummary(os.Stdout, mine)
	case "rest":
		id := command.MinerID
		fmt.Println("Ищем шахтёра с ID:", id)
		err := mine.RestMiner(id)
		if err != nil {
			if errors.Is(err, ErrMinerNotFound) {
				fmt.Println("Выберите другого шахтёра")
				return false
			}
			if errors.Is(err, ErrNotEnoughCoal) {
				fmt.Println("Недостаточно угля")
				return false
			}
			fmt.Println("Operation caused an error:", err)
			return false
		}
		fmt.Println("Шахтёр с ID:", id, "отдохнул")
		printMineSummary(os.Stdout, mine)
	case "status":
		printMineSummary(os.Stdout, mine)
	case "open":
		mine.Closed = false
		fmt.Println("Шахта открыта")
	case "close":
		mine.Closed = true
		fmt.Println("Шахта закрыта")
	case "exit":
		return true
	}
	return false
}
