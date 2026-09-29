package main

import "testing"

func TestMineShiftsWithMaintenance(t *testing.T) {
	got := mineShifts(5, 3, 2)
	want := 8

	if got != want {
		t.Errorf("mineShifts(5, 3, 2) = %d; want %d", got, want)
	}
}

func TestMineShiftsWithoutMaintenance(t *testing.T) {
	got := mineShifts(5, 8, 2)
	want := 10

	if got != want {
		t.Errorf("mineShifts(5, 8, 2) = %d; want %d", got, want)
	}
}

func TestMineShiftsOnlyMaintenance(t *testing.T) {
	got := mineShifts(1, 1, 2)
	want := 0

	if got != want {
		t.Errorf("mineShifts(1, 1, 2) = %d; want %d", got, want)
	}
}

func TestMineShiftsZeroProduction(t *testing.T) {
	got := mineShifts(5, 3, 0)
	want := 0

	if got != want {
		t.Errorf("mineShifts(5, 3, 0) = %d; want %d", got, want)
	}
}

func TestMineShiftsNegativeProduction(t *testing.T) {
	got := mineShifts(5, 3, -2)
	want := 0

	if got != want {
		t.Errorf("mineShifts(5, 3, -2) = %d; want %d", got, want)
	}
}

func TestMineShiftsZeroShifts(t *testing.T) {
	got := mineShifts(0, 3, 2)
	want := 0

	if got != want {
		t.Errorf("mineShifts(0, 3, 2) = %d; want %d", got, want)
	}
}

func TestMineShiftsNegativeShifts(t *testing.T) {
	got := mineShifts(-1, 3, 2)
	want := 0

	if got != want {
		t.Errorf("mineShifts(-1, 3, 2) = %d; want %d", got, want)
	}
}

func TestSummarizeProductionEmpty(t *testing.T) {
	gotBalance, gotShifts := summarizeProduction([]int{})
	want := 0

	if gotBalance != want || gotShifts != want {
		t.Errorf("summarizeProduction([]int{}) = %d; %d shifts, want %d", gotBalance, gotShifts, want)
	}
}

func TestSummarizeProductionMixed(t *testing.T) {
	gotBalance, gotShifts := summarizeProduction([]int{3, 0, 5, 0, 2})
	wantBalance := 10
	wantShifts := 3

	if gotBalance != wantBalance || gotShifts != wantShifts {
		t.Errorf("summarizeProduction([]int{3, 0, 5, 0, 2}) = %d; %d shifts, want balance %d, with %d shifts", gotBalance, gotShifts, wantBalance, wantShifts)
	}
}

func TestBuildProductionMaintenance(t *testing.T) {
	got := buildProduction(5, 3, 2)
	if len(got) != 5 {
		t.Fatalf("len = %d; want 5; production = %v", len(got), got)
	}

	for i := range got {
		want := 2
		if i == 2 {
			want = 0
		}
		if got[i] != want {
			t.Errorf("production[%d] = %d; want %d", i, got[i], want)
		}
	}
}

func TestBuildProductionWithoutMaintenance(t *testing.T) {
	got := buildProduction(3, 5, 2)
	if len(got) != 3 {
		t.Fatalf("len = %d; want 3, production = %v", len(got), got)
	}

	for i := range got {
		want := 2
		if got[i] != want {
			t.Errorf("production[%d] = %d; want %d", i, got[i], want)
		}
	}
}

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

func TestMineRound(t *testing.T) {
	miners := []Miner{
		{Name: "Борис", Energy: 2, Production: 3},
		{Name: "Анна", Energy: 1, Production: 5},
	}
	got := mineRound(miners)
	if got != 8 {
		t.Errorf("production = %d; want %d", got, 8)
	}
	if miners[0].Energy != 1 {
		t.Errorf("Boris energy = %d; want %d", miners[0].Energy, 1)
	}
	if miners[1].Energy != 0 {
		t.Errorf("Anna energy = %d; want %d", miners[1].Energy, 0)
	}
}
