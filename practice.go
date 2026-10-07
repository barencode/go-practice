package main

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
