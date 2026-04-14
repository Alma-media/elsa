package flow

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/Alma-media/elsa/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Publisher interface {
	Publish(string, byte, bool, interface{}) mqtt.Token
}

type Manager struct {
	mu sync.Mutex

	mqtt.Client

	subscriptions map[string]map[string]model.Options
}

// NewManager creates a new flow manager.
func NewManager(client mqtt.Client) *Manager { return &Manager{Client: client} }

// TODO:
// - detect circular deps
func (m *Manager) Apply(ctx context.Context, elements model.Pipe) (<-chan struct{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	await := make(chan struct{})

	m.subscriptions = make(map[string]map[string]model.Options)

	for _, element := range elements {
		outputs, ok := m.subscriptions[element.Input]
		if !ok {
			outputs = make(map[string]model.Options)
			m.subscriptions[element.Input] = outputs

			token := m.Subscribe(element.Input, 0, createHandler(m.Client, outputs))
			if token.Wait() && token.Error() != nil {
				return nil, token.Error()
			}
		}

		if _, ok := outputs[element.Output]; ok {
			return nil, fmt.Errorf("input %q and output %q already linked", element.Input, element.Output)
		}

		outputs[element.Output] = element.Options
	}

	go func() {
		<-ctx.Done()

		for topic := range m.subscriptions {
			if token := m.Unsubscribe(topic); token.Wait() && token.Error() != nil {
				log.Printf("cannot unsubscribe: %s", token.Error())
			}
		}

		close(await)
	}()

	return await, nil
}

func createHandler(publisher Publisher, outputs map[string]model.Options) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		for output, options := range outputs {
			publisher.Publish(output, 0, options.Retain, string(msg.Payload())).Wait()
		}
	}
}
