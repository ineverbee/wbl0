package store

import (
	"fmt"
	"time"
)

type DBIface interface {
	Set(*int, *Model) error
	Get(int) (*Model, error)
	GetAll() (map[int]*Model, error)
}

type CacheIface interface {
	Set(*int, *Model) error
	Get(int) (*Model, error)
}

type DBMock struct{}

func (dbmock *DBMock) Set(id *int, model *Model) error {
	*id = 1
	if model.Order_uid == "very_wrong_uid_for_db" {
		return fmt.Errorf("error")
	}
	return nil
}

func (dbmock *DBMock) Get(id int) (*Model, error) {
	return nil, nil
}

func (dbmock *DBMock) GetAll() (map[int]*Model, error) {
	return nil, nil
}

type CacheMock struct{}

func (cmock *CacheMock) Set(id *int, model *Model) error {
	if model.Order_uid == "very_wrong_uid_for_cache" {
		return fmt.Errorf("error")
	}
	return nil
}

func (cmock *CacheMock) Get(id int) (*Model, error) {
	if id == -10 {
		return nil, fmt.Errorf("error")
	}
	return nil, nil
}

type Delivery struct {
	Name    string `json:"name" sql:"name"`
	Phone   string `json:"phone" sql:"phone"`
	Zip     string `json:"zip" sql:"zip"`
	City    string `json:"city" sql:"city"`
	Address string `json:"address" sql:"address"`
	Region  string `json:"region" sql:"region"`
	Email   string `json:"email" sql:"email"`
}

type Payment struct {
	Transaction   string `json:"transaction" sql:"transaction"`
	Request_id    string `json:"request_id" sql:"request_id"`
	Currency      string `json:"currency" sql:"currency"`
	Provider      string `json:"provider" sql:"provider"`
	Bank          string `json:"bank" sql:"bank"`
	Amount        uint   `json:"amount" sql:"amount"`
	Payment_dt    uint   `json:"payment_dt" sql:"payment_dt"`
	Delivery_cost uint   `json:"delivery_cost" sql:"delivery_cost"`
	Goods_total   uint   `json:"goods_total" sql:"goods_total"`
	Custom_fee    uint   `json:"custom_fee" sql:"custom_fee"`
}

type Item struct {
	Track_number string `json:"track_number" sql:"track_number"`
	Rid          string `json:"rid" sql:"rid"`
	Name         string `json:"name" sql:"name"`
	Size         string `json:"size" sql:"size"`
	Brand        string `json:"brand" sql:"brand"`
	Chrt_id      uint   `json:"chrt_id" sql:"chrt_id"`
	Price        uint   `json:"price" sql:"price"`
	Sale         uint   `json:"sale" sql:"sale"`
	Total_price  uint   `json:"total_price" sql:"total_price"`
	Nm_id        uint   `json:"nm_id" sql:"nm_id"`
	Status       uint   `json:"status" sql:"status"`
}

type Model struct {
	Order_uid          string     `json:"order_uid" sql:"order_uid"`
	Track_number       string     `json:"track_number" sql:"track_number"`
	Entry              string     `json:"entry" sql:"entry"`
	Locale             string     `json:"locale" sql:"locale"`
	Internal_signature string     `json:"internal_signature" sql:"internal_signature"`
	Customer_id        string     `json:"customer_id" sql:"customer_id"`
	Delivery_service   string     `json:"delivery_service" sql:"delivery_service"`
	Shardkey           string     `json:"shardkey" sql:"shardkey"`
	Oof_shard          string     `json:"oof_shard" sql:"oof_shard"`
	Sm_id              uint       `json:"sm_id" sql:"sm_id"`
	Date_created       *time.Time `json:"date_created" sql:"date_created"`
	Delivery           *Delivery  `json:"delivery" sql:"delivery"`
	Payment            *Payment   `json:"payment" sql:"payment"`
	Items              []*Item    `json:"items" sql:"items"`
}
