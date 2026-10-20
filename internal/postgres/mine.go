package postgres

import (
	"context"

	"github.com/barencode/go-practice/internal/mining"
	"github.com/jackc/pgx/v5/pgxpool"
)

func LoadMine(ctx context.Context, pool *pgxpool.Pool, id int) (*mining.Mine, error) {
	var mine mining.Mine

	err := pool.QueryRow(ctx,
		"SELECT coal, reserves, closed FROM mines WHERE id = $1",
		id).Scan(&mine.Coal, &mine.Reserves, &mine.Closed)
	if err != nil {
		return nil, err
	}

	rows, err := pool.Query(ctx,
		"SELECT id, name, energy, production FROM miners WHERE mine_id = $1",
		id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var minerID int
		var miner mining.Miner
		err := rows.Scan(&minerID, &miner.Name, &miner.Energy, &miner.Production)
		if err != nil {
			return nil, err
		}
		mine.AddMiner(minerID, &miner)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &mine, nil
}
