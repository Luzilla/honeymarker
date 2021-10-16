package command

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/honeycombio/honeymarker/pkg/client"
	"github.com/honeycombio/honeymarker/pkg/marker"
)

func NewListCommand(json bool, ts bool, httpClient client.HoneyCombClient) *ListCommand {
	return &ListCommand{
		JSON:           json,
		UnixTimestamps: ts,
		httpClient:     httpClient,
	}
}

const (
	IdColumnWidth         = 11
	TimeColumnWidthPretty = 15
	TimeColumnWidthUnix   = 10
	TypeColumnWidth       = 12

	MessageColumnMaxWidth = 40
	MessageColumnMinWidth = len("Message")
	URLColumnMaxWidth     = 30
	URLColumnMinWidth     = len("URL")
)

func truncateStr(str string, maxWidth int) string {
	if len(str) > maxWidth {
		return str[:maxWidth-3] + "..."
	}
	return str
}

func (l *ListCommand) formatTime(timestamp int64) string {
	if timestamp == 0 {
		return ""
	}

	if l.UnixTimestamps {
		return strconv.FormatInt(timestamp, 10)
	}

	t := time.Unix(timestamp, 0)
	return t.Format(time.Stamp)
}

func (l *ListCommand) Execute() error {
	req, err := l.httpClient.NewRequest("GET", nil)
	if err != nil {
		return err
	}

	// if l.AuthorizationHeader != "" {
	// 	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", l.AuthorizationHeader))
	// }

	resp, err := client.MakeRequest(req)
	if err != nil {
		return err
	}

	body, err := client.ReadResponse(resp)
	if err != nil {
		return err
	}

	if l.JSON {
		return l.ListAsJSON(body)
	} else {
		return l.ListAsTable(body)
	}
}

func (l *ListCommand) ListAsJSON(body []byte) error {
	// newlineify the JSON for one marker per line
	// TODO json-pretty-print based on a flag or something
	prettyBody := strings.Replace(string(body), "},{", "},\n{", -1)
	fmt.Println(prettyBody)
	return nil
}

func (l *ListCommand) ListAsTable(body []byte) error {
	// Unmarshal string into structs.
	var mkrs []marker.Marker
	if err := json.Unmarshal(body, &mkrs); err != nil {
		return err
	}

	urlColumnWidth := 0
	messageColumnWidth := 0

	for _, m := range mkrs {
		if len(m.Message) > messageColumnWidth {
			messageColumnWidth = len(m.Message)
		}
		if len(m.URL) > urlColumnWidth {
			urlColumnWidth = len(m.URL)
		}
	}

	if messageColumnWidth > MessageColumnMaxWidth {
		messageColumnWidth = MessageColumnMaxWidth
	}
	if messageColumnWidth < MessageColumnMinWidth {
		messageColumnWidth = MessageColumnMinWidth
	}

	if urlColumnWidth > URLColumnMaxWidth {
		urlColumnWidth = URLColumnMaxWidth
	}
	if urlColumnWidth < URLColumnMinWidth {
		urlColumnWidth = URLColumnMinWidth
	}

	var timeColumnWidth int
	if l.UnixTimestamps {
		timeColumnWidth = TimeColumnWidthUnix
	} else {
		timeColumnWidth = TimeColumnWidthPretty
	}

	fmt.Printf("| %-[2]*[1]s | %[4]*[3]s | %[6]*[5]s | %-[8]*[7]s | %-[10]*[9]s | %-[12]*[11]s |\n",
		"ID", IdColumnWidth,
		"Start Time", timeColumnWidth,
		"End Time", timeColumnWidth,
		"Type", TypeColumnWidth,
		"Message", messageColumnWidth,
		"URL", urlColumnWidth,
	)
	fmt.Printf("+-%s-+-%s-+-%s-+-%s-+-%s-+-%s-+\n",
		strings.Repeat("-", IdColumnWidth),
		strings.Repeat("-", timeColumnWidth),
		strings.Repeat("-", timeColumnWidth),
		strings.Repeat("-", TypeColumnWidth),
		strings.Repeat("-", messageColumnWidth),
		strings.Repeat("-", urlColumnWidth),
	)
	for _, m := range mkrs {
		fmt.Printf("| %-[2]*[1]s | %[4]*[3]s | %[6]*[5]s | %-[8]*[7]s | %-[10]*[9]s | %-[12]*[11]s |\n",
			m.ID, IdColumnWidth,
			l.formatTime(m.StartTime), timeColumnWidth,
			l.formatTime(m.EndTime), timeColumnWidth,
			m.Type, TypeColumnWidth,
			truncateStr(m.Message, messageColumnWidth), messageColumnWidth,
			truncateStr(m.URL, urlColumnWidth), urlColumnWidth,
		)
	}

	return nil
}
