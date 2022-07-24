package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ineverbee/wbl0/internal/store"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	Error404NotFound     = fmt.Errorf("error: 404 not found")
	ErrorTimeoutExceeded = fmt.Errorf("db connection failed after timeout")

	SetQuery = `
INSERT INTO wb_data (order_uid,track_number,entry,delivery,payment,items,locale,internal_signature,customer_id,delivery_service,shardkey,sm_id,date_created,oof_shard) 
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING id`
	GetQuery    = "SELECT * FROM wb_data WHERE id=%d"
	GetAllQuery = "SELECT * FROM wb_data"
)

type PoolIface interface {
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

type DBStore struct {
	connPool PoolIface
}

func NewDBStore(ctx context.Context, connStr string, timeout time.Duration) (*DBStore, error) {
	log.Printf("Trying to connect to %s\n", connStr)
	var (
		conn *pgxpool.Pool
		err  error
	)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	timeoutExceeded := time.After(timeout)
LOOP:
	for {
		select {
		case <-timeoutExceeded:
			return nil, ErrorTimeoutExceeded

		case <-ticker.C:
			conn, err = pgxpool.Connect(ctx, connStr)
			if err == nil {
				break LOOP
			}
			log.Println("Failed! Trying to reconnect..")
		}
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	log.Println("Connect success!")

	return &DBStore{conn}, nil
}

func (db *DBStore) Set(id *int, m *store.Model) error {
	err := db.connPool.QueryRow(
		context.Background(), SetQuery,
		m.Order_uid,
		m.Track_number,
		m.Entry,
		m.Delivery,
		m.Payment,
		m.Items,
		m.Locale,
		m.Internal_signature,
		m.Customer_id,
		m.Delivery_service,
		m.Shardkey,
		m.Sm_id,
		m.Date_created,
		m.Oof_shard,
	).Scan(id)
	if err != nil {
		return err
	}
	return nil
}

func (db *DBStore) Get(id int) (*store.Model, error) {
	res := new(store.Model)
	q := fmt.Sprintf(GetQuery, id)
	err := db.connPool.QueryRow(context.Background(), q).Scan(
		&id,
		&res.Order_uid,
		&res.Track_number,
		&res.Entry,
		&res.Delivery,
		&res.Payment,
		&res.Items,
		&res.Locale,
		&res.Internal_signature,
		&res.Customer_id,
		&res.Delivery_service,
		&res.Shardkey,
		&res.Sm_id,
		&res.Date_created,
		&res.Oof_shard,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, Error404NotFound
		}
		return nil, err
	}
	return res, nil
}

func (db *DBStore) GetAll() (map[int]*store.Model, error) {
	m := make(map[int]*store.Model)
	rows, err := db.connPool.Query(context.Background(), GetAllQuery)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, Error404NotFound
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		id, temp := 0, new(store.Model)
		rows.Scan(
			&id,
			&temp.Order_uid,
			&temp.Track_number,
			&temp.Entry,
			&temp.Delivery,
			&temp.Payment,
			&temp.Items,
			&temp.Locale,
			&temp.Internal_signature,
			&temp.Customer_id,
			&temp.Delivery_service,
			&temp.Shardkey,
			&temp.Sm_id,
			&temp.Date_created,
			&temp.Oof_shard,
		)
		m[id] = temp
	}
	return m, nil
}
