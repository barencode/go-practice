package main

import "fmt"

func printMineSummary(mine *Mine) {
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
