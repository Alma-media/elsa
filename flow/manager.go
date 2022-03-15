package flow

import (
	"context"
	"fmt"
	"log"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Publisher interface {
	Publish(string, byte, bool, interface{}) mqtt.Token
}

type Manager struct {
	mu sync.Mutex

	mqtt.Client

	subscriptions map[string]map[string]Element
}

// NewManager creates a new flow manager.
func NewManager(client mqtt.Client) *Manager { return &Manager{Client: client} }

// TODO:
// - detect circular deps
func (m *Manager) Apply(ctx context.Context, elements Pipe) (<-chan struct{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	await := make(chan struct{})

	m.subscriptions = make(map[string]map[string]Element)

	for _, route := range elements {
		outputs, ok := m.subscriptions[route.Input.Path]
		if !ok {
			outputs = make(map[string]Element)
			m.subscriptions[route.Input.Path] = outputs

			token := m.Subscribe(route.Input.Path, 0, createHandler(m.Client, outputs))
			if token.Wait() && token.Error() != nil {
				return nil, token.Error()
			}
		}

		if _, ok := outputs[route.Output.Path]; ok {
			return nil, fmt.Errorf(
				"input %q and output %q already linked",
				route.Input.Path,
				route.Output.Path,
			)
		}

		outputs[route.Output.Path] = route.Output
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

func createHandler(publisher Publisher, outputs map[string]Element) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		for output, options := range outputs {
			publisher.Publish(output, 0, options.Retain, string(msg.Payload())).Wait()
		}
	}
}
