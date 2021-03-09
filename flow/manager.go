package flow

import (
	"context"
	"log"
	"sync"

	"github.com/Alma-media/elsa/pipe"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Manager struct{ mqtt.Client }

// NewManager creates a new manager.
func NewManager(client mqtt.Client) *Manager { return &Manager{Client: client} }

func (m *Manager) Apply(ctx context.Context, pipe pipe.Pipe) (<-chan struct{}, error) {
	var (
		await = make(chan struct{})
		wg    sync.WaitGroup
	)

	for _, element := range pipe {
		handler := func(client mqtt.Client, msg mqtt.Message) {
			var (
				data = msg.Payload()
				err  error
			)

			for _, processor := range element.Processors {
				data, err = processor.Process(data)
				if err != nil {
					return
				}
			}

			token := m.Publish(element.Output, 0, false, string(data))
			token.Wait()
		}

		// TODO: allow multiple handlers for the same topic (use map[topicName]<-chan []byte)
		if token := m.Subscribe(element.Input, 0, handler); token.Wait() && token.Error() != nil {
			return await, token.Error()
		}

		wg.Add(1)

		go func(topic string) {
			<-ctx.Done()

			if token := m.Unsubscribe(topic); token.Wait() && token.Error() != nil {
				log.Printf("cannot unsubscribe: %s", token.Error())

				return
			}

			wg.Done()
		}(element.Input)
	}

	go func() {
		wg.Wait()

		close(await)
	}()

	return await, nil
}
