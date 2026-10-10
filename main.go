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
