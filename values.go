package main

import (
	"errors"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
)

var valuesCommand *kingpin.CmdClause

func setupValuesCommand(app *kingpin.Application) {
	valuesCommand = app.Command("values", "Manage values")
}

func handleValuesCommand(commands []string) error {
	if len(commands) < 2 {
		return errors.New("No subcommand found for values command")
	}

	switch commands[1] {
	default:
		return errors.New(fmt.Sprintf("Values subcommand not supported : %s", commands[1]))
	}
}
