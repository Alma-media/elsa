package convert

import (
	"reflect"
	"testing"

	"github.com/Alma-media/elsa/flow"
)

var (
	pipe = flow.Pipe{
		{
			Input: flow.Element{
				Path: "/switch/2ccfd6b1-afdb-4c94-a651-5e112084c360",
			},
			Output: flow.Element{
				Path: "/trig.in.switch/bedroom",
			},
			Options: flow.Options{
				Description: "main swith for bedroom light",
			},
		},
		{
			Input: flow.Element{
				Path: "/switch/5d09399e-8a48-41c5-9f47-83127d8a69e0",
			},
			Output: flow.Element{
				Path: "/trig.in.switch/bedroom",
			},
			Options: flow.Options{
				Description: "additional switch for bedroom light",
			},
		},
		{
			Input: flow.Element{
				Path: "/trig.out.switch/bedroom",
			},
			Output: flow.Element{
				Path: "/led/2ccfd6b1-afdb-4c94-a651-5e112084c360",
			},
			Options: flow.Options{
				Description: "led to indicate main bedroom light switch state",
				Retain:      true,
			},
		},
		{
			Input: flow.Element{
				Path: "/trig.out.switch/bedroom",
			},
			Output: flow.Element{
				Path: "/led/5d09399e-8a48-41c5-9f47-83127d8a69e0",
			},
			Options: flow.Options{
				Description: "led to indicate additional bedroom light switch state",
				Retain:      true,
			},
		},
		{
			Input: flow.Element{
				Path: "/trig.out.switch/bedroom",
			},
			Output: flow.Element{
				Path: "/relay/3f8a79dc-854d-49d6-aa24-522422bf9140",
			},
			Options: flow.Options{
				Description: "relay module bedroom light",
				Retain:      true,
			},
		},
	}

	channels = []*Channel{
		{
			ID: "bedroom",
			Inputs: []*Device{
				{
					Type:        "switch",
					ID:          "2ccfd6b1-afdb-4c94-a651-5e112084c360",
					Description: "main swith for bedroom light",
				},
				{
					Type:        "switch",
					ID:          "5d09399e-8a48-41c5-9f47-83127d8a69e0",
					Description: "additional switch for bedroom light",
				},
			},
			Outputs: []*Device{
				{
					Type:        "led",
					ID:          "2ccfd6b1-afdb-4c94-a651-5e112084c360",
					Description: "led to indicate main bedroom light switch state",
				},
				{
					Type:        "led",
					ID:          "5d09399e-8a48-41c5-9f47-83127d8a69e0",
					Description: "led to indicate additional bedroom light switch state",
				},
				{
					Type:        "relay",
					ID:          "3f8a79dc-854d-49d6-aa24-522422bf9140",
					Description: "relay module bedroom light",
				},
			},
		},
	}
)

func TestPipeToChannels(t *testing.T) {
	actual, err := PipeToChannels(pipe)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if !reflect.DeepEqual(actual, channels) {
		t.Errorf("the result:\n%#v\ndoes not match expected:%#v", actual, channels)
	}
}

func TestChannelsToPipe(t *testing.T) {
	actual := ChannelsToPipe(channels)

	if !reflect.DeepEqual(actual, pipe) {
		t.Errorf("the result:\n%#v\ndoes not match expected:%#v", actual, channels)
	}
}
