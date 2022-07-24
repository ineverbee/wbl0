package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ineverbee/wbl0/internal/store"

	stan "github.com/nats-io/stan.go"
)

func subHandler(log *log.Logger, db store.DBIface, cache store.CacheIface) stan.MsgHandler {
	return func(m *stan.Msg) {
		d := m.Data
		if !json.Valid(d) {
			log.Printf("[WORKER] JSON Validation Error\n")
			return
		}
		id, unmarshData := -1, new(store.Model)
		decoder := json.NewDecoder(bytes.NewReader(d))
		err := decoder.Decode(unmarshData)
		if err != nil {
			log.Printf("[WORKER] Decode Error: %s\n", err.Error())
			return
		}
		err = validateFields(unmarshData)
		if err != nil {
			log.Printf("[WORKER] Field Validation Error: %s\n", err.Error())
			return
		}
		err = db.Set(&id, unmarshData)
		if err != nil {
			log.Printf("[WORKER] DB Error: %s\n", err.Error())
			return
		}
		if id != -1 {
			err = cache.Set(&id, unmarshData)
			if err != nil {
				log.Printf("[WORKER] Cache Error: %s\n", err.Error())
				return
			}
		}
	}
}

func Worker(db store.DBIface, cache store.CacheIface, sc stan.Conn, channel, durable string) error {
	// Subscribe with durable name
	sub, err := sc.Subscribe(channel, subHandler(log.Default(), db, cache), stan.DurableName(durable))

	if err != nil {
		log.Printf("[WORKER] Sub Error: %s\n", err.Error())
		return err
	}

	log.Printf("Listening on [%s], durable=[%s]\n", channel, durable)

	// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
	// Run cleanup when signal is received
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT)
	go func() {
		for range signalChan {
			log.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			// Do not unsubscribe a durable on exit, except if asked to.
			if durable == "" {
				sub.Unsubscribe()
			}
			sc.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
	return nil
}

func validateFields(m *store.Model) error {
	if m.Order_uid == "" ||
		m.Track_number == "" ||
		m.Entry == "" ||
		m.Locale == "" ||
		m.Customer_id == "" ||
		m.Delivery_service == "" ||
		m.Shardkey == "" ||
		m.Oof_shard == "" ||
		m.Sm_id <= 0 ||
		m.Date_created == nil ||
		m.Delivery == nil ||
		m.Payment == nil ||
		m.Items == nil {
		return fmt.Errorf("error: all fields should be provided, except optional field 'internal_signature'")
	}
	return nil
}
