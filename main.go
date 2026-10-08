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
		Miners: team,
		Closed: false,
	}

	coal, err := mine.MineByID(101)
	if err != nil {
		if errors.Is(err, ErrMineClosed) {
			fmt.Println("Добыча остановлена")
			return
		}
		if errors.Is(err, ErrMinerNotFound) {
			fmt.Println("Выберите другого шахтёра")
			return
		}
		fmt.Println("Operation caused an error:", err)
		return
	}
	fmt.Println("Coal mined:", coal)

	var newMine Mine
	newMine.AddMiner(101, &Miner{Name: "Борис", Energy: 2, Production: 3})
	newMine.AddMiner(205, &Miner{Name: "Анна", Energy: 3, Production: 2})

	fmt.Println(len(newMine.Miners))
}
