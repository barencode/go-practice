package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	team := map[int]*Miner{
		101: &Miner{Name: "Борис", Energy: 2, Production: 3},
		205: nil,
		309: &Miner{Name: "Анна", Energy: 1, Production: 5},
	}

	mine := Mine{
		Miners:   team,
		Closed:   false,
		Reserves: 10,
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Введите команду (mine <ID>, rest <ID>, status, exit, open, close): ")
		input, err := readLine(scanner)
		if errors.Is(err, io.EOF) {
			return
		}
		if err != nil {
			fmt.Println("Ошибка чтения:", err)
			return
		}

		command, err := parseCommand(input)
		if errors.Is(err, ErrEmptyCommand) {
			continue
		}
		if err != nil {
			fmt.Println("Ошибка команды:", err)
			continue
		}

		if executeCommand(&mine, command) {
			return
		}
	}
}
