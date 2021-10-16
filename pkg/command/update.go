package command

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/honeycombio/honeymarker/pkg/client"
	"github.com/honeycombio/honeymarker/pkg/marker"
)

type UpdateCommand struct {
	MarkerID  string
	StartTime int64
	EndTime   int64
	Message   string
	URL       string
	Type      string
	client    client.HoneyCombClient
}

func NewUpdateCommand(id string, start int64, end int64, msg string, url string, markerType string, client client.HoneyCombClient) *UpdateCommand {
	return &UpdateCommand{
		MarkerID:  id,
		StartTime: start,
		EndTime:   end,
		Message:   msg,
		URL:       url,
		Type:      markerType,
		client:    client,
	}
}

func (u *UpdateCommand) Execute() error {
	blob, err := json.Marshal(marker.Marker{
		StartTime: u.StartTime,
		EndTime:   u.EndTime,
		Message:   u.Message,
		Type:      u.Type,
		URL:       u.URL,
	})
	if err != nil {
		return err
	}

	// WEIRD
	// u.APIHost = fmt.Sprintf("%s/%s", u.APIHost, u.MarkerID)

	// c := client.NewClient(u.APIHost, u.Dataset, u.WriteKey, UserAgent)
	req, err := u.client.NewRequest(
		"PUT",
		bytes.NewBuffer(blob),
	)
	if err != nil {
		return err
	}

	// if u.AuthorizationHeader != "" {
	// 	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", u.AuthorizationHeader))
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
