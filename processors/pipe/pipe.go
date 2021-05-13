package pipe

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/Alma-media/elsa/pipe"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type PipeProcessor struct{
	mqtt.Client
}

func NewPipeProcessor(client 	mqtt.Client) *PipeProcessor {
	return &PipeProcessor{
		client,
	}
}


func Pipe ProcessorFunc = func()
func(in []byte) ([]byte, error) {
	log.Printf("Message: %s", in)

	return in, nil
}