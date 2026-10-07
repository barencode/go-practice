package main

import "fmt"

type Mine struct {
	Miners map[int]*Miner
	Closed bool
}

func (mine *Mine) MineByID(id int) (int, error) {
	if mine.Closed {
		return 0, fmt.Errorf("ошибка шахты: %w", ErrMineClosed)
	}
	miner, ok := mine.Miners[id]
	if !ok || miner == nil {
		return 0, fmt.Errorf("miner %d: %w", id, ErrMinerNotFound)
	}
	return miner.MineOnce(), nil
}
