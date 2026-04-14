package model

type Options struct {
	Retain bool `json:"retain"`
}

type Element struct {
	Alias  string `json:"alias"`
	Input  string `json:"input"`
	Output string `json:"output"`

	Options
}
