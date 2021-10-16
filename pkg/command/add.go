package command

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/honeycombio/honeymarker/pkg/client"
	"github.com/honeycombio/honeymarker/pkg/marker"
)

func NewAddCommand(start int64, end int64, msg string, url string, markerType string, httpClient client.HoneyCombClient) *AddCommand {
	return &AddCommand{
		StartTime: start,
		EndTime:   end,
		Message:   msg,
		URL:       url,
		Type:      markerType,
		client:    httpClient,
	}
}

func (a *AddCommand) Execute() error {
	blob, err := json.Marshal(marker.Marker{
		StartTime: a.StartTime,
		EndTime:   a.EndTime,
		Message:   a.Message,
		Type:      a.Type,
		URL:       a.URL,
	})
	if err != nil {
		return err
	}

	req, err := a.client.NewRequest(
		"POST",
		bytes.NewBuffer(blob),
	)
	if err != nil {
		return err
	}

	// if a.AuthorizationHeader != "" {
	// 	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.AuthorizationHeader))
	// }

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
