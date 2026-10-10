package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
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

	printMineSummary(&mine)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Введите ID шахтёра (0 - выход): ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Println("Ошибка чтения:", err)
			}
			return
		}

		input := scanner.Text()

		id, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			fmt.Println("Некорректный ID")
			continue
		}
		if id == 0 {
			return
		}

		fmt.Println("Ищем шахтёра с ID:", id)
		coal, err := mine.MineByID(id)
		if err != nil {
			if errors.Is(err, ErrMineClosed) {
				fmt.Println("Добыча остановлена")
				return
			}
			if errors.Is(err, ErrMinerNotFound) {
				fmt.Println("Выберите другого шахтёра")
				continue
			}
			fmt.Println("Operation caused an error:", err)
			return
		}
		fmt.Println("Coal mined:", coal)
		printMineSummary(&mine)
	}
}
