package flow

// Pipe is a linear list of input/output bindings.
type Pipe []Route

type Options struct {
	Retain      bool   `json:"retain"`
	Description string `json:"description"`
}

type Route struct {
	Input  Element `json:"input"`
	Output Element `json:"output"`
	Options
}

type Element struct {
	Path string `json:"path"`
	// Arguments []Argument `json:"arguments"`
}

type Argument struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}
