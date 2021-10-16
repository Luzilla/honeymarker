package command

import (
	"fmt"

	"github.com/honeycombio/honeymarker/pkg/client"
)

func NewRmCommand(id string, client client.HoneyCombClient) *RmCommand {
	return &RmCommand{
		MarkerID: id,
		client:   client,
	}
}

func (r *RmCommand) Execute() error {
	req, err := r.client.NewRequest(
		"DELETE",
		nil,
	)

	resp, err := client.MakeRequest(req)
	if err != nil {
		return err
	}

	body, err := client.ReadResponse(resp)
	if err != nil {
		return err
	}

	fmt.Println(string(body))
	return nil
}
