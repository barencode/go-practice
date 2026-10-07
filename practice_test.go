package main

import "testing"

type RoundExpectation struct {
	Production  int
	BorisEnergy int
	AnnaEnergy  int
}

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

func TestMineRoundUntilExhausted(t *testing.T) {
	team := []Miner{
		{Name: "Борис", Energy: 2, Production: 3},
		{Name: "Анна", Energy: 1, Production: 5},
	}

	expectations := []RoundExpectation{
		{Production: 8, BorisEnergy: 1, AnnaEnergy: 0},
		{Production: 3, BorisEnergy: 0, AnnaEnergy: 0},
		{Production: 0, BorisEnergy: 0, AnnaEnergy: 0},
	}

	for i, want := range expectations {
		got := mineRound(team)
		if got != want.Production {
			t.Errorf("in the cycle %d got %d; want %d", i+1, got, want.Production)
		}
		if team[0].Energy != want.BorisEnergy {
			t.Errorf("in the cycle %d 1st miner energy is %d; want %d", i+1, team[0].Energy, want.BorisEnergy)
		}
		if team[1].Energy != want.AnnaEnergy {
			t.Errorf("in the cycle %d 2nd miner energy is %d; want %d", i+1, team[1].Energy, want.AnnaEnergy)
		}
	}
}

func TestFindMinerIndex(t *testing.T) {
	team := []Miner{
		{Name: "Борис", Energy: 2, Production: 3},
		{Name: "Анна", Energy: 1, Production: 5},
		{Name: "Борис", Energy: 4, Production: 2},
	}
	requestedName := "Борис"
	want := 0
	got := findMinerIndex(team, requestedName)

	if got != want {
		t.Errorf("got index for %q name: %d; want %d", requestedName, got, want)
	}
}

func TestFindMinerIndexNotFound(t *testing.T) {
	team := []Miner{
		{Name: "Борис", Energy: 2, Production: 3},
		{Name: "Анна", Energy: 1, Production: 5},
	}
	requestedName := "Иван"
	want := -1
	got := findMinerIndex(team, requestedName)

	if got != want {
		t.Errorf("got index for %q name: %d; want %d", requestedName, got, want)
	}
}

func TestMineMapRound(t *testing.T) {
	team := map[int]*Miner{
		101: &Miner{Name: "Борис", Energy: 2, Production: 3},
		205: nil,
		309: &Miner{Name: "Анна", Energy: 1, Production: 5},
	}
	got := mineMapRound(team)
	wantAmount := 8
	boris := team[101]
	anna := team[309]
	empty, exists := team[205]

	if got != wantAmount {
		t.Errorf("got %d coal mining amount; want %d", got, wantAmount)
	}
	if boris.Energy != 1 {
		t.Errorf("got %d energy of %q; want %d", boris.Energy, boris.Name, 1)
	}
	if anna.Energy != 0 {
		t.Errorf("got %d energy of %q; want %d", anna.Energy, anna.Name, 0)
	}

	if !exists || empty != nil {
		t.Errorf("entry 205: exists=%t, value=%v; want true, nil", exists, empty)
	}
}
