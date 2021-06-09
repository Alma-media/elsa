package pipe

import (
	"bytes"
	"encoding/binary"
	"log"
)

// Processor is a struct/function to process the input data before putting it to the output.
type Processor interface {
	Process(data []byte) ([]byte, error)
}

type ProcessorFunc func([]byte) ([]byte, error)

func (fn ProcessorFunc) Process(data []byte) ([]byte, error) {
	return fn(data)
}

var Count ProcessorFunc = func(in []byte) ([]byte, error) {
	var buff bytes.Buffer

	if err := binary.Write(&buff, binary.LittleEndian, int64(len(in))); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

var Reverse ProcessorFunc = func(in []byte) ([]byte, error) {
	out := make([]byte, len(in))

	for index, symbol := range in {
		out[len(in)-index-1] = symbol
	}

	return out, nil
}

var Print ProcessorFunc = func(in []byte) ([]byte, error) {
	log.Printf("Message: %s", in)

	return in, nil
}

var processors = map[string]Processor{
	"count":   Count,
	"print":   Print,
	"reverse": Reverse,
}
