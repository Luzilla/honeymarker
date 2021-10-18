package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/honeycombio/honeymarker/pkg/client"
	"github.com/honeycombio/honeymarker/pkg/command"
	"github.com/urfave/cli/v2"
)

// BuildID is set by CircleCI
var BuildID string

// Sets the internal version ID and updates libhoney's user-agent
func CreateVersionUserAgent(buildID string) string {
	var version string

	if buildID == "" {
		version = "dev"
	} else {
		version = buildID
	}
	return fmt.Sprintf("honeymarker/%s", version)
}

var httpClient client.HoneyCombClient

func main() {
	app := &cli.App{
		Name: "honeymarker",
		Authors: []*cli.Author{
			{
				Name:  "Honeycomb",
				Email: "solutions@honeycomb.io",
			},
		},
		Usage: "honeymarker is the command line utility for manipulating markers in your Honeycomb dataset.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "writekey",
				Aliases:  []string{"k"},
				EnvVars:  []string{"HONEYCOMB_API_KEY"},
				Usage:    "Honeycomb write key from https://ui.honeycomb.io/account",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "dataset",
				Aliases:  []string{"d"},
				EnvVars:  []string{"HONEYCOMB_DATASET"},
				Usage:    "Honeycomb dataset name from https://ui.honeycomb.io/dashboard",
				Required: true,
			},
			&cli.StringFlag{
				Name:   "api_host",
				Value:  "https://api.honeycomb.io/",
				Hidden: true,
			},
			&cli.StringFlag{
				Name:   "authorization-header",
				Hidden: true,
			},
		},
		Before: func(c *cli.Context) error {
			httpClient = client.NewClient(
				c.String("api_host"),
				c.String("dataset"),
				c.String("writekey"),
				CreateVersionUserAgent(BuildID),
			)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "add",
				Usage: "Add a new marker",
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:    "start_time",
						Aliases: []string{"s"},
						Usage:   "Start time for the marker in unix time (seconds since the epoch)",
						Value:   time.Now().Unix(),
					},
					&cli.Int64Flag{
						Name:    "end_time",
						Usage:   "End time for the marker in unix time (seconds since the epoch)",
						Aliases: []string{"e"},
						Hidden:  true,
					},
					&cli.StringFlag{
						Name:    "msg",
						Aliases: []string{"m"},
						Usage:   "Message describing this specific marker",
					},
					&cli.StringFlag{
						Name:    "url",
						Aliases: []string{"u"},
						Usage:   "URL associated with this marker",
					},
					&cli.StringFlag{
						Name:    "type",
						Aliases: []string{"t"},
						Usage:   "Identifies marker type",
					},
				},
				Action: func(c *cli.Context) error {
					add := command.NewAddCommand(
						c.Int64("start_time"),
						c.Int64("end_time"),
						c.String("message"),
						c.String("url"),
						c.String("type"),
						httpClient,
					)
					return add.Execute()
				},
			},
			{
				Name:  "list",
				Usage: "List all markers",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "json",
						Usage: "Output the list as json instead of in tabular form.",
					},
					&cli.BoolFlag{
						Name:  "unix_time",
						Usage: "In table mode, format times as unit timestamps (seconds since the epoch)",
					},
				},
				Action: func(c *cli.Context) error {
					list := command.NewListCommand(c.Bool("json"), c.Bool("unix_time"), httpClient)
					return list.Execute()
				},
			},
			{
				Name:  "rm",
				Usage: "Delete a marker",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Aliases:  []string{"i"},
						Usage:    "ID of the marker to delete",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					rm := command.NewRmCommand(c.String("id"), httpClient)
					return rm.Execute()
				},
			},
			{
				Name:  "update",
				Usage: "Update a marker",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Aliases:  []string{"i"},
						Usage:    "ID for the marker to update",
						Required: true,
					},
					&cli.Int64Flag{
						Name:    "start_time",
						Aliases: []string{"s"},
						Usage:   "start time for the marker in unix time",
						Value:   time.Now().Unix(),
					},
					&cli.Int64Flag{
						Name:    "end_time",
						Aliases: []string{"e"},
						Hidden:  true,
						Usage:   "end time for the marker in unix time",
					},
					&cli.StringFlag{
						Name:    "msg",
						Aliases: []string{"m"},
						Usage:   "message to attach to the marker",
					},
					&cli.StringFlag{
						Name:    "url",
						Aliases: []string{"u"},
						Usage:   "url to attach to the marker",
					},
					&cli.StringFlag{
						Name:    "type",
						Aliases: []string{"t"},
						Usage:   "type of the marker",
					},
				},
				Action: func(c *cli.Context) error {
					update := command.NewUpdateCommand(
						c.String("id"),
						c.Int64("start_time"),
						c.Int64("end_time"),
						c.String("msg"),
						c.String("url"),
						c.String("type"),
						httpClient,
					)
					return update.Execute()
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	// 	parser.AddCommand("add", "",
	// 		`add creates a new marker with the specified attributes.

	//   All parameters to add are optional.

	//   If start_time is missing, the marker will be assigned the current time.

	//   It is highly recommended that you fill in either message or type.
	//   All markers of the same type will be shown with the same color in the UI.
	//   The message will be visible above an individual marker.

	// 	If a URL is specified along with a message, the message will be shown
	// 	as a link in the UI, and clicking it will take you to the URL.`,
	// 		&command.AddCommand{})

	// 	parser.AddCommand("list", "",
	// 		`List all markers for the specified dataset.

	//   Returned markers will be displayed in tabular format by default,
	// 	ordered by the marker's start time.`,
	// 		&command.ListCommand{})

	// 	parser.AddCommand("rm", "",
	// 		`Delete the marker in the specified dataset, as identified by its ID.

	// 	Marker IDs are available via the 'list' command.`,
	// 		&command.RmCommand{})

	// 	parser.AddCommand("update", "",
	// 		`Update an existing marker in the specified dataset with the specified options.

	// 	The marker ID is required (available via the 'list' command). All other
	// 	parameters are optional, though an 'update' will be a no-op unless a parameter
	// 	is specified with a new value.`,
	// 		&command.UpdateCommand{})

	// 	// run whichever command is chosen
	// 	parser.Usage = usage
	// 	if _, err := parser.Parse(); err != nil {
	// 		if flagErr, ok := err.(*flag.Error); ok {
	// 			if flagErr.Type == flag.ErrHelp {
	// 				// asking for help isn't a failed run.
	// 				os.Exit(0)
	// 			}
	// 			if flagErr.Type == flag.ErrCommandRequired ||
	// 				flagErr.Type == flag.ErrUnknownFlag ||
	// 				flagErr.Type == flag.ErrRequired {
	// 				fmt.Println("  run 'honeymarker --help' for full usage details")
	// 			}
	// 		}
	// 		os.Exit(1)
	// 	}
}
