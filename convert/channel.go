package convert

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Alma-media/elsa/flow"
)

const (
	inputPrefix  = "/trig.in.switch/"
	outputPrefix = "/trig.out.switch/"
)

type Device struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	Description string `json:"description"`
}

func (device Device) Path() string {
	return fmt.Sprintf("/%s/%s", device.Type, device.ID)
}

type Channel struct {
	ID      string    `json:"id"`
	Inputs  []*Device `json:"inputs"`
	Outputs []*Device `json:"outputs"`
}

func (channel Channel) Input() string {
	return inputPrefix + channel.ID
}

func (channel Channel) Output() string {
	return outputPrefix + channel.ID
}

type Channels []*Channel

func (channels Channels) Find(id string) (int, bool) {
	for index, channel := range channels {
		if channel.ID == id {
			return index, true
		}
	}

	return 0, false
}

var deviceExp = regexp.MustCompile(`^/(\w+)/(.*)$`)

func deviceFromPath(path string) (*Device, error) {
	matches := deviceExp.FindStringSubmatch(path)

	if len(matches) != 3 {
		return nil, fmt.Errorf("bad path")
	}

	return &Device{
		Type: matches[1],
		ID:   matches[2],
	}, nil
}

func PipeToChannels(pipe flow.Pipe) ([]*Channel, error) {
	var list Channels

	for _, element := range pipe {
		var (
			id      string
			channel *Channel
		)
		// check if it is state manager route
		switch {
		// output route (led, relay, uhf ...)
		case strings.HasPrefix(element.Input.Path, "/trig.out.switch/"):
			id = element.Input.Path[17:]

			if index, ok := list.Find(id); ok {
				channel = list[index]
			} else {
				channel = &Channel{
					ID: id,
				}
				list = append(list, channel)
			}

			device, err := deviceFromPath(element.Output.Path)
			if err != nil {
				return nil, err
			}

			device.Description = element.Description
			channel.Outputs = append(channel.Outputs, device)
		// input route (knob, switch, uhf, rfid ...)
		case strings.HasPrefix(element.Output.Path, inputPrefix):
			id = element.Output.Path[len(inputPrefix):]

			if index, ok := list.Find(id); ok {
				channel = list[index]
			} else {
				channel = &Channel{
					ID: id,
				}
				list = append(list, channel)
			}

			device, err := deviceFromPath(element.Input.Path)
			if err != nil {
				return nil, err
			}

			device.Description = element.Description
			channel.Inputs = append(channel.Inputs, device)
		default:
			return nil, fmt.Errorf("error")
		}
	}

	return list, nil
}

func ChannelsToPipe(list []*Channel) flow.Pipe {
	var pipe flow.Pipe

	for _, channel := range list {
		for _, input := range channel.Inputs {
			route := flow.Route{
				Input: flow.Element{
					Path: input.Path(),
				},
				Output: flow.Element{
					Path: channel.Input(),
				},
				Options: flow.Options{
					Description: input.Description,
				},
			}

			pipe = append(pipe, route)
		}

		for _, output := range channel.Outputs {
			route := flow.Route{
				Input: flow.Element{
					Path: channel.Output(),
				},
				Output: flow.Element{
					Path: output.Path(),
				},
				Options: flow.Options{
					Retain:      true,
					Description: output.Description,
				},
			}

			pipe = append(pipe, route)
		}
	}

	return pipe
}
