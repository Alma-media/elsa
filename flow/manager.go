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
	mu sync.Mutex

	mqtt.Client

	// map[input][output][]Processor
	subscriptions map[string]map[string][]pipe.Processor
}

// NewManager creates a new manager.
func NewManager(client mqtt.Client) *Manager { return &Manager{Client: client} }

// TODO:
// - detect circular deps
func (m *Manager) Apply(ctx context.Context, elements pipe.Pipe) (<-chan struct{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	await := make(chan struct{})

	m.subscriptions = make(map[string]map[string][]pipe.Processor)

	for _, element := range elements {
		outputs, ok := m.subscriptions[element.Input]
		if !ok {
			outputs = make(map[string][]pipe.Processor)
			m.subscriptions[element.Input] = outputs

			token := m.Subscribe(element.Input, 0, createHandler(m.Client, outputs))
			if token.Wait() && token.Error() != nil {
				return nil, token.Error()
			}
		}

		if _, ok := outputs[element.Output]; ok {
			return nil, fmt.Errorf("input %q and output %q already linked", element.Input, element.Output)
		}

		outputs[element.Output] = element.Processors
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
