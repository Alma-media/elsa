package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alma-media/elsa/api"
	"github.com/Alma-media/elsa/flow"
	"github.com/Alma-media/elsa/pipe"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("[DEFAULT HANDLER] TOPIC: %s\n", msg.Topic())
	fmt.Printf("[DEFAULT HANDLER] MSG:   %s\n", msg.Payload())
}

func main() {
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("emqx_test_client")

	opts.SetKeepAlive(60 * time.Second)
	// Message callback handler
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(time.Second)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to initialize a client: %s", token.Error())
	}

	storage := &pipe.InMemoryStorage{
		Pipe: []pipe.Element{
			{
				BaseElement: pipe.BaseElement{
					Input:  "/input",
					Output: "/output",
				},
				Processors: []pipe.Processor{
					pipe.Print,
					pipe.Reverse,
				},
			},
			{
				BaseElement: pipe.BaseElement{
					Input:  "/a",
					Output: "/b",
				},
				Processors: []pipe.Processor{
					pipe.Print,
					pipe.Reverse,
				},
			},
			{
				BaseElement: pipe.BaseElement{
					Input:  "/a",
					Output: "/c",
				},
				Processors: []pipe.Processor{
					pipe.Print,
					pipe.Reverse,
				},
			},
		},
	}

	manager := flow.NewManager(client)

	handler, err := api.NewHandler(storage, manager)
	if err != nil {
		log.Fatalf("failed to create a handler: %s", err)
	}

	go func() {
		http.ListenAndServe(":8888", http.HandlerFunc(handler.ApplyHandler))
	}()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	handler.Stop()
}