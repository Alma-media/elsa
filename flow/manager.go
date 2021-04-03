package flow

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/Alma-media/elsa/pipe"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Manager struct {
	mqtt.Client

	// map[input][output][]Processor
	subscriptions map[string]map[string][]pipe.Processor
}

// NewManager creates a new manager.
func NewManager(client mqtt.Client) *Manager { return &Manager{Client: client} }

// TODO:
// - detect circular deps
func (m *Manager) Apply(ctx context.Context, elements pipe.Pipe) (<-chan struct{}, error) {
	var (
		await = make(chan struct{})
		wg    sync.WaitGroup
	)

	m.subscriptions = make(map[string]map[string][]pipe.Processor)

	for _, element := range elements {
		outputs, ok := m.subscriptions[element.Input]
		if !ok {
			outputs = make(map[string][]pipe.Processor)
			m.subscriptions[element.Input] = outputs

			handler := createHandler(m.Client, outputs)

			// TODO: allow multiple handlers for the same topic (use map[topicName]<-chan []byte)
			if token := m.Subscribe(element.Input, 0, handler); token.Wait() && token.Error() != nil {
				return nil, token.Error()
			}

		}

		if _, ok := outputs[element.Output]; ok {
			return nil, fmt.Errorf("input %q and output %q already linked", element.Input, element.Output)
		}

		outputs[element.Output] = element.Processors
	}

	wg.Add(1)

	go func() {
		<-ctx.Done()

		defer wg.Done()

		for topic := range m.subscriptions {
			if token := m.Unsubscribe(topic); token.Wait() && token.Error() != nil {
				log.Printf("cannot unsubscribe: %s", token.Error())
			}
		}
	}()

	go func() {
		wg.Wait()

		close(await)
	}()

	return await, nil
}

type Publisher interface {
	Publish(string, byte, bool, interface{}) mqtt.Token
}

func createHandler(publisher Publisher, outputs map[string][]pipe.Processor) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		var (
			data = msg.Payload()
			err  error
		)

		for output, processors := range outputs {
			go func(output string, processors []pipe.Processor) {
				for _, processor := range processors {
					data, err = processor.Process(data)
					if err != nil {
						return
					}
				}

				token := publisher.Publish(output, 0, false, string(data))
				token.Wait()
			}(output, processors)
		}
	}
}
