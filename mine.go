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

func (mine *Mine) AddMiner(id int, miner *Miner) {
	if mine.Miners == nil {
		mine.Miners = make(map[int]*Miner)
	}
	mine.Miners[id] = miner
}

func (mine *Mine) RemoveMiner(id int) error {
	_, ok := mine.Miners[id]
	if !ok {
		return fmt.Errorf("miner %d: %w", id, ErrMinerNotFound)
	}
	delete(mine.Miners, id)
	return nil
}
