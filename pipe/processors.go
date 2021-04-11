package pipe

import "log"

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
	"print":   Print,
	"reverse": Reverse,
}
