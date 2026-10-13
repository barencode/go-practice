package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Command struct {
	Name    string
	MinerID int
}

func printMineSummary(mine *Mine) {
	fmt.Printf("На складе: %d, в залежах: %d\n", mine.Coal, mine.Reserves)

	ids := mine.MinerIDs()
	for _, id := range ids {
		miner, err := mine.FindMiner(id)
		if err != nil {
			fmt.Printf("%d: %v\n", id, err)
			continue
		}
		fmt.Printf("%d: %s: энергия %d, добыча за смену %d\n", id, miner.Name, miner.Energy, miner.Production)
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
