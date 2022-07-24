package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/ineverbee/wbl0/internal/store"
	"github.com/ineverbee/wbl0/internal/store/db"
	"github.com/ineverbee/wbl0/internal/store/mapstore"
	"github.com/ineverbee/wbl0/internal/worker"
	"github.com/nats-io/stan.go"
)

type App struct {
	server *http.Server
	db     store.DBIface
	cache  store.CacheIface
}

var app *App

func StartApp() error {
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	router := mux.NewRouter()

	router.Handle("/", limit(errorHandler(GetHomePageHandler()))).Methods("GET", "POST")
	router.Handle("/data/{id}", limit(errorHandler(GetDataPageHandler()))).Methods("GET")

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?connect_timeout=5",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	dbStore, err := db.NewDBStore(ctx, connStr, 30*time.Second)
	if err != nil {
		return err
	}

	mapStore := mapstore.NewMapStore(make(map[int]*store.Model, 1))

	app = &App{
		&http.Server{Addr: ":8080", Handler: router},
		dbStore,
		mapStore,
	}

	mp, err := app.db.GetAll()

	if err == nil {
		app.cache = mapstore.NewMapStore(mp)
	} else if err != db.Error404NotFound {
		return err
	}

	sc, err := stan.Connect(
		os.Getenv("NATS_CLUSTER_ID"),
		os.Getenv("NATS_CLIENT_ID"),
		stan.NatsURL(os.Getenv("NATS_URL")))
	if err != nil {
		log.Printf("[WORKER] Error: %s\n", err.Error())
		return err
	}

	go worker.Worker(
		app.db,
		app.cache,
		sc,
		os.Getenv("NATS_CHANNEL"),
		os.Getenv("NATS_DURABLE"),
	)

	log.Println("Starting server on Port 8080")
	err = http.ListenAndServe(":8080", router)
	return err
}
