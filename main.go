package main

import "fmt"

func main() {
	fmt.Println(mineShifts(5, 3, 2))
	fmt.Println(mineShifts(5, 8, 2))
	fmt.Println(mineShifts(1, 1, 2))
	fmt.Println(mineShifts(5, 3, 0))
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
