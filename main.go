package main

import "fmt"

func main() {
	balance := 0

	miner := Miner{
		Name:       "Борис",
		Energy:     3,
		Production: 2,
	}

	for miner.Energy > 0 {
		balance += miner.Production
		miner.Energy--
	}

	fmt.Println(miner.Name, "has complete his work and has", miner.Energy, "energy left. Total coal mined is:", balance)
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

func mineOnce(miner *Miner) int {
	if miner.Energy <= 0 {
		return 0
	}
	miner.Energy--
	return miner.Production
}
