package main

import "fmt"

type Mine struct {
	Miners map[int]*Miner
	Closed bool
	Coal   int
}

func (mine *Mine) MineByID(id int) (int, error) {
	if mine.Closed {
		return 0, fmt.Errorf("ошибка шахты: %w", ErrMineClosed)
	}
	miner, err := mine.FindMiner(id)
	if err != nil {
		return 0, err
	}
	mining := miner.MineOnce()
	mine.Coal += mining
	return mining, nil
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

func (mine *Mine) FindMiner(id int) (*Miner, error) {
	miner, ok := mine.Miners[id]
	if !ok || miner == nil {
		return nil, fmt.Errorf("miner %d: %w", id, ErrMinerNotFound)
	}
	return miner, nil
}
