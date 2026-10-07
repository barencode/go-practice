package main

import (
	"errors"
	"testing"
)

func TestMineByIDNotFound(t *testing.T) {
	team := map[int]*Miner{}
	testMine := Mine{
		Miners: team,
		Closed: false,
	}
	got, err := testMine.MineByID(999)

	if got != 0 || !errors.Is(err, ErrMinerNotFound) {
		t.Errorf("got coal=%d, err=%v; want coal=0, err=%v", got, err, ErrMinerNotFound)
	}
}

func TestMineByIDNilMiner(t *testing.T) {
	team := map[int]*Miner{
		205: nil,
	}
	testMine := Mine{
		Miners: team,
		Closed: false,
	}
	got, err := testMine.MineByID(205)

	if got != 0 || !errors.Is(err, ErrMinerNotFound) {
		t.Errorf("got coal=%d, err=%v; want coal=0 and %v error", got, err, ErrMinerNotFound)
	}
}

func TestMineByIDExhausted(t *testing.T) {
	team := map[int]*Miner{
		101: {Name: "Борис", Energy: 0, Production: 3},
	}
	testMine := Mine{
		Miners: team,
		Closed: false,
	}
	got, err := testMine.MineByID(101)

	if !(got == 0 && err == nil) {
		t.Errorf("got coal=%d, err=%v; want coal=0 and nil error", got, err)
	}
}

func TestMineByIDWithEnergy(t *testing.T) {
	team := map[int]*Miner{
		101: {Name: "", Energy: 2, Production: 3},
	}
	testMine := Mine{
		Miners: team,
		Closed: false,
	}
	got, err := testMine.MineByID(101)
	wantProduction := 3 // equals to production and 1 cycle of mining
	wantEnergy := 1

	if got != wantProduction || team[101].Energy != wantEnergy || err != nil {
		t.Errorf("got coal=%d, err=%v and miner's energy=%d; want coal=%d, nil error and energy=%d", got, err, team[101].Energy, wantProduction, wantEnergy)
	}
}

func TestMineByIDClosed(t *testing.T) {
	team := map[int]*Miner{
		101: {Name: "Борис", Energy: 2, Production: 3},
	}
	testMine := Mine{
		Miners: team,
		Closed: true,
	}
	got, err := testMine.MineByID(101)
	wantProduction := 0
	wantEnergy := 2

	if got != wantProduction || team[101].Energy != wantEnergy || !errors.Is(err, ErrMineClosed) {
		t.Errorf("got production=%d, miner energy=%d and err=%v; want production=%d, miner energy=%d and error=%v", got, team[101].Energy, err, wantProduction, wantEnergy, ErrMineClosed)
	}
}

func TestMineByIDClosedNotFound(t *testing.T) {
	team := map[int]*Miner{}
	testMine := Mine{
		Miners: team,
		Closed: true,
	}
	got, err := testMine.MineByID(999)
	wantError := ErrMineClosed

	if got != 0 || !errors.Is(err, wantError) {
		t.Errorf("got coal=%d, err=%v; want coal=0, err=%v",
			got, err, wantError)
	}
}
