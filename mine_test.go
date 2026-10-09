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
		Coal:   10,
	}
	got, err := testMine.MineByID(101)
	wantProduction := 0
	wantEnergy := 2

	if got != wantProduction || team[101].Energy != wantEnergy || !errors.Is(err, ErrMineClosed) || testMine.Coal != 10 {
		t.Errorf("got production=%d, miner energy=%d, total coal=%d and err=%v; want production=%d, miner energy=%d, total coal=10 and error=%v", got, team[101].Energy, testMine.Coal, err, wantProduction, wantEnergy, ErrMineClosed)
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

func TestAddMinerInitializesMap(t *testing.T) {
	var mine Mine
	boris := &Miner{"Борис", 2, 3}
	mine.AddMiner(101, boris)
	got := mine.Miners[101]
	if len(mine.Miners) != 1 || got != boris {
		t.Errorf("got team length=%d, miner=%p, want: 1 and %p", len(mine.Miners), got, boris)
	}
}

func TestAddMinerPreservesExisting(t *testing.T) {
	var mine Mine
	boris := &Miner{"Борис", 2, 3}
	anna := &Miner{"Анна", 3, 2}
	mine.AddMiner(101, boris)
	mine.AddMiner(205, anna)
	got := mine.Miners
	if len(got) != 2 || got[101] != boris || got[205] != anna {
		t.Errorf("got team length=%d, 1st miner=%p, 2nd miner=%p; want: 2 and %p and %p", len(got), got[101], got[205], boris, anna)
	}
}

func TestAddMinerReplacesExisting(t *testing.T) {
	var mine Mine
	boris := &Miner{"Борис", 2, 3}
	anna := &Miner{"Анна", 3, 2}
	mine.AddMiner(101, boris)
	mine.AddMiner(101, anna)
	got := mine.Miners

	if len(got) != 1 || got[101] != anna {
		t.Errorf("got team length=%d, miner=%p; want: 1 and %p", len(got), got[101], anna)
	}
}

func TestRemoveMinerExisting(t *testing.T) {
	var mine Mine
	boris := &Miner{"Борис", 2, 3}
	anna := &Miner{"Анна", 3, 2}
	mine.AddMiner(101, boris)
	mine.AddMiner(205, anna)
	err := mine.RemoveMiner(101)
	_, exists := mine.Miners[101]

	if err != nil || len(mine.Miners) != 1 || mine.Miners[205] != anna || exists {
		t.Errorf("got error=%v, team length=%d, miner=%p, deleted key exists=%t; want nil, 1, %p, false",
			err, len(mine.Miners), mine.Miners[205], exists, anna)
	}
}

func TestRemoveMinerNilValue(t *testing.T) {
	var mine Mine
	mine.AddMiner(101, nil)
	err := mine.RemoveMiner(101)
	_, exists := mine.Miners[101]

	if err != nil || exists {
		t.Errorf("got error=%v, deleted miner exists=%t; want nil, false", err, exists)
	}
}

func TestRemoveMinerNotFound(t *testing.T) {
	var mine Mine
	err := mine.RemoveMiner(999)

	if !errors.Is(err, ErrMinerNotFound) {
		t.Errorf("got error=%v; want error=%v", err, ErrMinerNotFound)
	}
}

func TestFindMinerExisting(t *testing.T) {
	var mine Mine
	boris := &Miner{"Борис", 2, 3}
	mine.AddMiner(101, boris)
	miner, err := mine.FindMiner(101)

	if miner != boris || err != nil {
		t.Errorf("got miner=%p error=%v; want miner=%p and nil", miner, err, boris)
	}
}

func TestFindMinerNotFound(t *testing.T) {
	var mine Mine
	miner, err := mine.FindMiner(999)

	if miner != nil || !errors.Is(err, ErrMinerNotFound) {
		t.Errorf("got miner=%p error=%v; want miner=nil and error=%v", miner, err, ErrMinerNotFound)
	}
}

func TestFindMinerNilValue(t *testing.T) {
	var mine Mine
	mine.AddMiner(101, nil)
	miner, err := mine.FindMiner(101)

	if miner != nil || !errors.Is(err, ErrMinerNotFound) {
		t.Errorf("got miner=%p error=%v; want miner=nil and error=%v", miner, err, ErrMinerNotFound)
	}
}

func TestFindMinerWhenClosed(t *testing.T) {
	var mine Mine
	mine.Closed = true
	boris := &Miner{"Борис", 2, 3}
	mine.AddMiner(101, boris)
	miner, err := mine.FindMiner(101)

	if miner != boris || err != nil {
		t.Errorf("got miner=%p error=%v; want miner=%p and nil", miner, err, boris)
	}
}

func TestMineByIDAccumulatesCoal(t *testing.T) {
	mine := Mine{Coal: 10}
	mine.AddMiner(101, &Miner{"Борис", 2, 3})

	shiftProduction, err := mine.MineByID(101)
	if shiftProduction != 3 || mine.Coal != 13 || err != nil {
		t.Errorf("after 1st shift: got production=%d, total coal=%d, err=%v; want 3, 13, nil", shiftProduction, mine.Coal, err)
	}
	shiftProduction, err = mine.MineByID(101)
	if shiftProduction != 3 || mine.Coal != 16 || err != nil {
		t.Errorf("after 2nd shift: got production=%d, total coal=%d, err=%v; want 3, 16, nil", shiftProduction, mine.Coal, err)
	}
	shiftProduction, err = mine.MineByID(101)
	if shiftProduction != 0 || mine.Coal != 16 || err != nil {
		t.Errorf("after 3rd shift: got production=%d, total coal=%d, err=%v; want 0, 16, nil", shiftProduction, mine.Coal, err)
	}
}
