package flow

import (
	"github.com/Alma-media/elsa/pipe"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

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
