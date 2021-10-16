package command

import "github.com/honeycombio/honeymarker/pkg/client"

type Options struct {
}

type AddCommand struct {
	StartTime int64
	EndTime   int64
	Message   string
	URL       string
	Type      string
	client    client.HoneyCombClient
}

type ListCommand struct {
	JSON           bool
	UnixTimestamps bool
	httpClient     client.HoneyCombClient
}

type RmCommand struct {
	MarkerID string
	client   client.HoneyCombClient
}
