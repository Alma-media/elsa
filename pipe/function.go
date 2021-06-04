package pipe

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var expr = regexp.MustCompile(`^(\w+)\((.*)\)$`)

// TODO: use state machine to parse the function + arguments
// func Pipe(['foo-in', 'foo-out'], ['bar-in', 'bar-out']) means means send input to 'foo' and output of 'bar' to 'baz'
// func Pipe('foo', 'bar') means send input to 'foo' and wait for output from 'bar'

type Function struct {
	Name string
	Args []interface{}
}

func NewFunction(input string) (*Function, error) {
	var (
		fn      Function
		matches = expr.FindStringSubmatch(input)
	)

	fmt.Println(matches)

	if len(matches) != 3 {
		return nil, errors.New("error")
	}

	fn.Name = matches[1]

	if matches[2] == "" {
		return &fn, nil
	}

	args := strings.Split(matches[2], ",")

	for _, arg := range args {
		fn.Args = append(fn.Args, arg)
	}

	return &fn, nil
}
