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
