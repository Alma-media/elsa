package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("reverse")

	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(time.Second)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to initialize a client: %s", token.Error())
	}

	token := client.Subscribe("/reverse", 0, func(client mqtt.Client, msg mqtt.Message) {
		var (
			in = msg.Payload()
		)

		out := make([]byte, len(in))

		for index, symbol := range in {
			out[len(in)-index-1] = symbol
		}

		token := client.Publish("/reverse-out", 0, false, string(out))
		token.Wait()
	})

	if token.Wait() && token.Error() != nil {
		log.Fatalf("failed to subscribe: %s", token.Error())
	}

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan
}
