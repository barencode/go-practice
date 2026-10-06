package main

import (
	"errors"
	"fmt"
)

type Miner struct {
	Name       string
	Energy     int
	Production int
}

type Mine struct {
	Miners map[int]*Miner
	Closed bool
}

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
}

func mineShifts(shifts int, maintenanceShift int, production int) int {
	if shifts <= 0 || production <= 0 {
		return 0
	}

	balance := 0

	for ; shifts > 0; shifts-- {
		if shifts == maintenanceShift {
			continue
		}
		balance += production
	}
	return balance
}

func summarizeProduction(production []int) (int, int) {
	balance := 0
	workingShifts := 0
	for _, coal := range production {
		balance += coal
		if coal > 0 {
			workingShifts++
		}
	}
	return balance, workingShifts
}

func buildProduction(shifts int, maintenanceShift int, coalPerShift int) []int {
	production := []int{}
	for shift := 1; shift <= shifts; shift++ {
		if shift == maintenanceShift {
			production = append(production, 0)
		} else {
			production = append(production, coalPerShift)
		}
	}
	return production
}

func (miner *Miner) MineOnce() int {
	if miner.Energy <= 0 {
		return 0
	}
	miner.Energy--
	return miner.Production
}

func mineRound(miners []Miner) int {
	balance := 0

	for i := range miners {
		balance += miners[i].MineOnce()
	}

	return balance
}

func findMinerIndex(miners []Miner, name string) int {
	for i, miner := range miners {
		if name == miner.Name {
			return i
		}
	}
	return -1
}

func mineMapRound(miners map[int]*Miner) int {
	amount := 0
	for _, miner := range miners {
		if miner != nil {
			amount += miner.MineOnce()
		}
	}
	return amount
}

func (mine *Mine) MineByID(id int) (int, error) {
	if mine.Closed {
		return 0, fmt.Errorf("ошибка шахты: %w", ErrMineClosed)
	}
	miner, ok := mine.Miners[id]
	if !ok || miner == nil {
		return 0, fmt.Errorf("miner %d: %w", id, ErrMinerNotFound)
	}
	return miner.MineOnce(), nil
}
