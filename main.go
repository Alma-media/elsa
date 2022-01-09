package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alma-media/elsa/api"
	"github.com/Alma-media/elsa/config"
	"github.com/Alma-media/elsa/flow"
	"github.com/Alma-media/elsa/storage/database"
	"github.com/Alma-media/elsa/storage/database/sqlite"
	"github.com/Alma-media/elsa/storage/memory"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	conf "github.com/tiny-go/config"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("[DEFAULT HANDLER] TOPIC: %s\n", msg.Topic())
	fmt.Printf("[DEFAULT HANDLER] MSG:   %s\n", msg.Payload())
}

type Storage interface {
	Load(context.Context) (flow.Pipe, error)
	Save(context.Context, flow.Pipe) error
}

func main() {
	var (
		appConfig config.Config
		storage   Storage
		ctx       = context.Background()
	)

	if err := conf.Init(&appConfig, ""); err != nil {
		log.Fatalf("unable to parse the config: %s", err)
	}

	switch appConfig.Storage.Type {
	case "memory":
		storage = new(memory.Storage)
	case "database":
		db, err := sql.Open(appConfig.Storage.Database.Driver, appConfig.Storage.Database.DSN)
		if err != nil {
			log.Fatalf("unable to establish database connection: %s", err)
		}

		if err := sqlite.Init(ctx, db); err != nil {
			log.Fatalf("database migration failure: %s", err)
		}

		defer db.Close()

		storage = database.NewStorage(db, new(sqlite.PipeManager))
	default:
		log.Fatalf("unknown storage type %q", appConfig.Storage.Type)
	}

	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.
		NewClientOptions().
		AddBroker(appConfig.Broker.DSN).
		SetClientID(appConfig.Broker.ClientID)

	opts.SetKeepAlive(60 * time.Second)
	// Message callback handler
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(time.Second)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to initialize a client: %s", token.Error())
	}

	handler, err := api.NewHandler(storage, flow.NewManager(client))
	if err != nil {
		log.Fatalf("failed to create a handler: %s", err)
	}

	go func() {
		log.Printf("Started API on port %d", appConfig.HTTP.Port)

		http.ListenAndServe(
			fmt.Sprintf(":%d", appConfig.HTTP.Port),
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					handler.LoadHandler(w, r)
				case http.MethodPost:
					handler.ApplyHandler(w, r)
				default:
					http.Error(w, fmt.Sprintf("Method %q not allowed", r.Method), http.StatusMethodNotAllowed)
				}
			}),
		)
	}()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	handler.Stop()
}
