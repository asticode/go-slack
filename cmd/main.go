package main

import (
	"flag"

	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astitools/flag"
	"github.com/asticode/go-slack"
)

// Flags
var (
	channel = flag.String("c", "", "the channel")
	message = flag.String("m", "", "the message")
)

func main() {
	// Get subcommand
	s := astiflag.Subcommand()
	flag.Parse()

	// Init logger
	astilog.FlagInit()

	// Init slack
	sl := slack.New(slack.FlagConfig())

	// Log
	astilog.Debugf("Subcommand is %s", s)

	// Switch on subcommand
	switch s {
	default:
		// Init message
		m := slack.Message{
			Channel: *channel,
			Text:    *message,
		}

		// Slack
		if err := sl.Slack(m); err != nil {
			astilog.Fatal(err)
		}
		break
	}
}
