package flow

import (
	"context"
	"fmt"
	"log"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Pipe is a linear list of onput/output bindings.
type Pipe []Element

type Element struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type Publisher interface {
	Publish(string, byte, bool, interface{}) mqtt.Token
}

type Manager struct {
	mu sync.Mutex

	mqtt.Client

	subscriptions map[string]map[string]struct{}
}

// NewManager creates a new flow manager.
func NewManager(client mqtt.Client) *Manager { return &Manager{Client: client} }

// TODO:
// - detect circular deps
func (m *Manager) Apply(ctx context.Context, elements Pipe) (<-chan struct{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	await := make(chan struct{})

	m.subscriptions = make(map[string]map[string]struct{})

	for _, element := range elements {
		outputs, ok := m.subscriptions[element.Input]
		if !ok {
			outputs = make(map[string]struct{})
			m.subscriptions[element.Input] = outputs

			token := m.Subscribe(element.Input, 0, createHandler(m.Client, outputs))
			if token.Wait() && token.Error() != nil {
				return nil, token.Error()
			}
		}

		if _, ok := outputs[element.Output]; ok {
			return nil, fmt.Errorf("input %q and output %q already linked", element.Input, element.Output)
		}

		outputs[element.Output] = struct{}{}
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

func createHandler(publisher Publisher, outputs map[string]struct{}) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		for output := range outputs {
			publisher.Publish(output, 0, false, string(msg.Payload())).Wait()
		}
	}
}
