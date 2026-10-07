package main

import "testing"

func TestMineOnceWithEnergy(t *testing.T) {
	gotMiner := Miner{
		Name:       "Test",
		Energy:     2,
		Production: 3,
	}
	got := gotMiner.MineOnce()
	if got != 3 {
		t.Errorf("production = %d; want %d", got, 3)
	}
	if gotMiner.Energy != 1 {
		t.Errorf("energy left = %d; want %d", gotMiner.Energy, 1)
	}
}

func TestMineOnceWithoutEnergy(t *testing.T) {
	gotMiner := Miner{
		Name:       "Exhausted Test",
		Energy:     0,
		Production: 2,
	}
	got := gotMiner.MineOnce()
	if got != 0 {
		t.Errorf("production = %d; want %d", got, 0)
	}
	if gotMiner.Energy != 0 {
		t.Errorf("energy left = %d; want %d", gotMiner.Energy, 0)
	}
}
