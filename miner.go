package main

type Miner struct {
	Name       string
	Energy     int
	Production int
}

func (miner *Miner) MineOnce() int {
	if miner.Energy <= 0 {
		return 0
	}
	miner.Energy--
	return miner.Production
}
