package main

import "fmt"

func main() {
	miners := []Miner{
		{Name: "Борис", Energy: 2, Production: 3},
		{Name: "Анна", Energy: 1, Production: 5},
	}

	total := 0
	total += mineRound(miners)
	total += mineRound(miners)

	fmt.Println("Total balance:", total, ". Boris' energy:", miners[0].Energy, ". Annas energy:", miners[1].Energy)
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

type Miner struct {
	Name       string
	Energy     int
	Production int
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
