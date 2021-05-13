package pipe

import (
	"errors"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// TODO: replace local processors
// processors["count", "reverse"]
// means forward IN to -> count: process: count-out -> reverse: proceess: reverse-out -> OUT

// Processor is a struct/function to process the input data before putting it to the output.
type Processor interface {
	Process(data []byte) ([]byte, error)
}

type ProcessorFunc func([]byte) ([]byte, error)

func (fn ProcessorFunc) Process(data []byte) ([]byte, error) {
	return fn(data)
}

func Reverse(args ...interface{}) (ProcessorFunc, error) {
	if len(args) > 0 {
		return nil, errors.New("no arguments expected")
	}

	return func(in []byte) ([]byte, error) {
		out := make([]byte, len(in))

		for index, symbol := range in {
			out[len(in)-index-1] = symbol
		}

		return out, nil
	}, nil
}

func Print(args ...interface{}) (ProcessorFunc, error) {
	if len(args) > 0 {
		return nil, errors.New("no arguments expected")
	}

	return func(in []byte) ([]byte, error) {
		log.Printf("Message: %s", in)

		return in, nil
	}, nil
}

var processors = map[string]ProcessorFactory{
	"print":   Print,
	"reverse": Reverse,
}

type ProcessorFactory func(args ...interface{}) (ProcessorFunc, error)

type PipeProcessor struct {
	mqtt.Client
}

func NewPipeProcessor(client mqtt.Client) *PipeProcessor {
	return &PipeProcessor{
		client,
	}
}

func (proc *PipeProcessor) Pipe(in []byte) ([]byte, error) {
	log.Printf("Message: %s", in)

	return in, nil
}
