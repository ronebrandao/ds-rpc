package db

import (
	"context"
	"ds-rpc/server/model"
	"github.com/jackc/pgx/v4"
)

type DB struct {
	conn *pgx.Conn
	ctx context.Context
}

func (db* DB) OpenConnection(ctx context.Context) error {
	conn, err := pgx.Connect(ctx, "postgresql://localhost/postgres?user=rone&password=123456")
	if err != nil {
		return err
	}

	db.conn = conn
	db.ctx = ctx
	return nil
}

func (db *DB) CloseConnection() {
	db.conn.Close(db.ctx)
}

func (db *DB) SavePerformanceStats(status model.Status) error {
	_, err := db.conn.Exec(db.ctx, `insert into 
		performance_info(cpu,memory_used,memory_avaliable,disk_used,disk_avaliable)
		values($1,$2,$3,$4,$5)`, status.CPU, status.UsedRAM, status.AvaliableRAM, status.UsedDisk, status.AvaliableDisk)

		return err
}
