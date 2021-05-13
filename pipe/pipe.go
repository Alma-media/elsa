package pipe

import (
	"encoding/json"
	"fmt"
)

// 1. Collect all inputs and outputs
// 2. Provide a consumer for every input => goroutine to push raw data into the channel
// 3. Prepare a producer for every output => goroutine to publish raw data to the topic once it is received from the channel
// 4. Link input/output channels (maybe add handlers between input and output to process the data)

// Pipe is a linear list of onput/output bindings.
type Pipe []Element

type BaseElement struct {
	Input  string   `json:"input"`
	Output string   `json:"output"`
	Pipe   []string `json:"pipe"`
}

// Element is a single pipe element.
type Element struct {
	BaseElement

	Processors []Processor
}

func (e *Element) UnmarshalJSON(data []byte) error {
	var base BaseElement

	if err := json.Unmarshal(data, &base); err != nil {
		return err
	}

	e.BaseElement = base

	for _, element := range base.Pipe {

		fn, err := NewFunction(element)
		if err != nil {
			return err
		}

		processorFactory, ok := processors[fn.Name]
		if !ok {
			return fmt.Errorf("unknown processor: %s", element)
		}

		processor, err := processorFactory(fn.Args...)
		if err != nil {
			return err
		}

		e.Processors = append(e.Processors, processor)
	}

	return nil
}
