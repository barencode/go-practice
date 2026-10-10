package main

import (
	"errors"
	"fmt"
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

	for {
		var id int
		fmt.Print("Введите ID шахтёра (0 - выход): ")
		_, err := fmt.Scan(&id)

		if err != nil {
			fmt.Println("Некорректный ID")
			return
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
