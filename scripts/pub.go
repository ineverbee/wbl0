package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/nats-io/stan.go"
)

func main() {
	cluster, client, subj := "test-cluster", "test-pub", "foo"

	sc, err := stan.Connect(cluster, client, stan.NatsURL("http://localhost:4222"))
	if err != nil {
		log.Fatalf("Can't connect: %v.\n", err)
	}
	defer sc.Close()
	input := 0
	fmt.Print("How many rows to add? ")
	fmt.Scanf("%d", &input)
	for ; input > 0; input-- {
		msg := []byte(newJSON())

		err = sc.Publish(subj, msg)
		if err != nil {
			log.Fatalf("Error during publish: %v\n", err)
		}
		log.Printf("Published [%s] : '%s'\n", subj, msg)
	}

}

func newJSON() string {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(5)
	items := ""
	for i := 0; i < n; i++ {
		items += fmt.Sprintf(`{
			"chrt_id": %d,
			"track_number": "%s",
			"price": %d,
			"rid": "%s",
			"name": "%s",
			"sale": %d,
			"size": "%d",
			"total_price": %d,
			"nm_id": %d,
			"brand": "%s",
			"status": %d
		  }`,
			gofakeit.Uint32(),
			gofakeit.LetterN(10),
			gofakeit.Uint16(),
			gofakeit.BitcoinAddress(),
			gofakeit.Word(),
			gofakeit.Uint8(),
			gofakeit.Uint8(),
			gofakeit.Uint16(),
			gofakeit.Uint32(),
			gofakeit.Company(),
			gofakeit.Uint8(),
		)
		if i != n-1 {
			items += ","
		}
	}
	return fmt.Sprintf(`{
		"order_uid": "%s",
		"track_number": "%s",
		"entry": "%s",
		"delivery": {
		  "name": "%s",
		  "phone": "%s",
		  "zip": "%s",
		  "city": "%s",
		  "address": "%s",
		  "region": "%s",
		  "email": "%s"
		},
		"payment": {
		  "transaction": "%s",
		  "request_id": "%d",
		  "currency": "%s",
		  "provider": "%s",
		  "amount": %d,
		  "payment_dt": %d,
		  "bank": "%s",
		  "delivery_cost": %d,
		  "goods_total": %d,
		  "custom_fee": %d
		},
		"items": [
		  %s
		],
		"locale": "%s",
		"internal_signature": "%s",
		"customer_id": "%s",
		"delivery_service": "%s",
		"shardkey": "%d",
		"sm_id": %d,
		"date_created": "%s",
		"oof_shard": "%d"
	  }`,
		gofakeit.BitcoinAddress(),
		gofakeit.LetterN(10),
		gofakeit.Word(),
		gofakeit.Name(),
		gofakeit.Phone(),
		gofakeit.Zip(),
		gofakeit.City(),
		gofakeit.Street(),
		gofakeit.State(),
		gofakeit.Email(),
		gofakeit.BitcoinAddress(),
		gofakeit.Uint32(),
		gofakeit.Currency().Short,
		gofakeit.Company(),
		gofakeit.Uint32(),
		gofakeit.Uint32(),
		gofakeit.Company(),
		gofakeit.Uint16(),
		gofakeit.Uint16(),
		gofakeit.Uint8(),
		items,
		gofakeit.CountryAbr(),
		gofakeit.UUID(),
		gofakeit.LetterN(20),
		gofakeit.Company(),
		gofakeit.Uint8(),
		gofakeit.Uint8(),
		gofakeit.Date().Format(time.RFC3339),
		gofakeit.Uint8(),
	)
}
