package db

import (
	"context"
	"testing"
	"time"

	"github.com/ineverbee/wbl0/internal/store"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/require"
)

var (
	exampleModel = &store.Model{
		Order_uid:    "b563feb7b2b84b6test",
		Track_number: "WBILMTESTTRACK",
		Entry:        "WBIL",
		Delivery: &store.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: &store.Payment{
			Transaction:   "b563feb7b2b84b6test",
			Request_id:    "",
			Currency:      "USD",
			Provider:      "wbpay",
			Amount:        1817,
			Payment_dt:    1637907727,
			Bank:          "alpha",
			Delivery_cost: 1500,
			Goods_total:   317,
			Custom_fee:    0,
		},
		Items: []*store.Item{
			{
				Chrt_id:      9934930,
				Track_number: "WBILMTESTTRACK",
				Price:        453,
				Rid:          "ab4219087a764ae0btest",
				Name:         "Mascaras",
				Sale:         30,
				Size:         "0",
				Total_price:  317,
				Nm_id:        2389212,
				Brand:        "Vivienne Sabo",
				Status:       202,
			},
		},
		Locale:             "en",
		Internal_signature: "",
		Customer_id:        "test",
		Delivery_service:   "meest",
		Shardkey:           "9",
		Sm_id:              99,
		Date_created:       &time.Time{},
		Oof_shard:          "1",
	}
)

func TestNewDBStore(t *testing.T) {
	res, err := NewDBStore(context.Background(), "", 11*time.Second)
	require.Nil(t, res)
	require.ErrorIs(t, ErrorTimeoutExceeded, err)
}

func TestDBStore(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()
	id, model := 1, exampleModel
	model_map := map[int]*store.Model{id: model}
	dbStore := &DBStore{mock}

	// Testing 'Set', not expecting any error, expecting 1 row
	mock.ExpectQuery("INSERT INTO wb_data").WithArgs(
		model.Order_uid,
		model.Track_number,
		model.Entry,
		model.Delivery,
		model.Payment,
		model.Items,
		model.Locale,
		model.Internal_signature,
		model.Customer_id,
		model.Delivery_service,
		model.Shardkey,
		model.Sm_id,
		model.Date_created,
		model.Oof_shard,
	).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	err = dbStore.Set(&id, model)
	require.NoError(t, err)

	// Testing 'Set', expecting ErrNoRows error
	mock.ExpectQuery("INSERT INTO wb_data").WithArgs(
		model.Order_uid,
		model.Track_number,
		model.Entry,
		model.Delivery,
		model.Payment,
		model.Items,
		model.Locale,
		model.Internal_signature,
		model.Customer_id,
		model.Delivery_service,
		model.Shardkey,
		model.Sm_id,
		model.Date_created,
		model.Oof_shard,
	).WillReturnError(pgx.ErrNoRows)
	err = dbStore.Set(&id, model)
	require.ErrorIs(t, pgx.ErrNoRows, err)

	// Testing 'Get', not expecting any error, expecting 1 row
	mock.ExpectQuery("SELECT (.+) FROM wb_data").WillReturnRows(pgxmock.NewRows(
		[]string{
			"id",
			"order_uid",
			"track_number",
			"entry",
			"delivery",
			"payment",
			"items",
			"locale",
			"internal_signature",
			"customer_id",
			"delivery_service",
			"shardkey",
			"sm_id",
			"date_created",
			"oof_shard",
		}).AddRow(
		id,
		model.Order_uid,
		model.Track_number,
		model.Entry,
		model.Delivery,
		model.Payment,
		model.Items,
		model.Locale,
		model.Internal_signature,
		model.Customer_id,
		model.Delivery_service,
		model.Shardkey,
		model.Sm_id,
		model.Date_created,
		model.Oof_shard,
	))
	res, err := dbStore.Get(id)
	require.NoError(t, err)
	require.IsType(t, &store.Model{}, res)
	require.Equal(t, model, res)

	// Testing 'Get', expecting Error404NotFound error
	mock.ExpectQuery("SELECT (.+) FROM wb_data").WillReturnError(pgx.ErrNoRows)
	res, err = dbStore.Get(id)
	require.Nil(t, res)
	require.ErrorIs(t, Error404NotFound, err)

	// Testing 'Get', expecting any other error except ErrNoRows, Error404NotFound
	mock.ExpectQuery("SELECT (.+) FROM wb_data").WillReturnError(pgx.ErrTxClosed)
	res, err = dbStore.Get(id)
	require.Nil(t, res)
	require.Error(t, err)
	require.NotErrorIs(t, pgx.ErrNoRows, err)
	require.NotErrorIs(t, Error404NotFound, err)

	// Testing 'GetAll', not expecting any error, expecting 1 row
	mock.ExpectQuery("SELECT (.+) FROM wb_data").WillReturnRows(pgxmock.NewRows(
		[]string{
			"id",
			"order_uid",
			"track_number",
			"entry",
			"delivery",
			"payment",
			"items",
			"locale",
			"internal_signature",
			"customer_id",
			"delivery_service",
			"shardkey",
			"sm_id",
			"date_created",
			"oof_shard",
		}).AddRow(
		id,
		model.Order_uid,
		model.Track_number,
		model.Entry,
		model.Delivery,
		model.Payment,
		model.Items,
		model.Locale,
		model.Internal_signature,
		model.Customer_id,
		model.Delivery_service,
		model.Shardkey,
		model.Sm_id,
		model.Date_created,
		model.Oof_shard,
	))
	res_map, err := dbStore.GetAll()
	require.NoError(t, err)
	require.IsType(t, model_map, res_map)
	require.Equal(t, model_map, res_map)

	// Testing 'GetAll', expecting Error404NotFound error
	mock.ExpectQuery("SELECT (.+) FROM wb_data").WillReturnError(pgx.ErrNoRows)
	res_map, err = dbStore.GetAll()
	require.Nil(t, res_map)
	require.ErrorIs(t, Error404NotFound, err)

	// Testing 'GetAll', expecting any other error except ErrNoRows, Error404NotFound
	mock.ExpectQuery("SELECT (.+) FROM wb_data").WillReturnError(pgx.ErrTxClosed)
	res_map, err = dbStore.GetAll()
	require.Nil(t, res_map)
	require.Error(t, err)
	require.NotErrorIs(t, pgx.ErrNoRows, err)
	require.NotErrorIs(t, Error404NotFound, err)
}
