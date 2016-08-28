package main

import (
	"flag"

	"github.com/asticode/go-slack"
	"github.com/molotovtv/go-logger"
	"github.com/molotovtv/go-toolbox"
	"github.com/rs/xlog"
)

// Flags
var (
	channel = flag.String("c", "", "the channel")
	message = flag.String("m", "", "the message")
)

func main() {
	// Get subcommand
	s := toolbox.Subcommand()
	flag.Parse()

	// Init logger
	l := xlog.New(logger.NewConfig(logger.FlagConfig()))

	// Init slack
	sl := slack.New(slack.FlagConfig())
	sl.Logger = l

	// Log
	l.Debugf("Subcommand is %s", s)

	// Switch on subcommand
	switch s {
	default:
		break
	}
}
