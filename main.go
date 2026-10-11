package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
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

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Введите команду (mine, status, rest, exit): ")
		input, err := readLine(scanner)
		if errors.Is(err, io.EOF) {
			return
		}
		if err != nil {
			fmt.Println("Ошибка чтения:", err)
			return
		}

		command := strings.TrimSpace(input)

		switch command {
		case "status":
			printMineSummary(&mine)
		case "exit":
			return
		case "mine":
			fmt.Print("Введите ID шахтёра: ")
			input, err := readLine(scanner)
			if errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				fmt.Println("Ошибка чтения:", err)
				return
			}
			id, err := parseMinerID(input)
			if err != nil {
				fmt.Println("Некорректный ID")
				continue
			}
			fmt.Println("Ищем шахтёра с ID:", id)
			coal, err := mine.MineByID(id)
			if err != nil {
				if errors.Is(err, ErrMineClosed) {
					fmt.Println("Добыча остановлена")
					continue
				}
				if errors.Is(err, ErrMinerNotFound) {
					fmt.Println("Выберите другого шахтёра")
					continue
				}
				fmt.Println("Operation caused an error:", err)
				continue
			}
			fmt.Println("Coal mined:", coal)
			printMineSummary(&mine)
		case "rest":
			fmt.Print("Введите ID шахтёра: ")
			input, err := readLine(scanner)
			if errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				fmt.Println("Ошибка чтения:", err)
				return
			}
			id, err := parseMinerID(input)
			if err != nil {
				fmt.Println("Некорректный ID")
				continue
			}
			fmt.Println("Ищем шахтёра с ID:", id)
			err = mine.RestMiner(id)
			if err != nil {
				if errors.Is(err, ErrMinerNotFound) {
					fmt.Println("Выберите другого шахтёра")
					continue
				}
				if errors.Is(err, ErrNotEnoughCoal) {
					fmt.Println("Недостаточно угля")
					continue
				}
				fmt.Println("Operation caused an error:", err)
				continue
			}
			fmt.Println("Шахтёр с ID:", id, "отдохнул")
			printMineSummary(&mine)
		default:
			fmt.Println("Неизвестная команда")
		}
	}
}
