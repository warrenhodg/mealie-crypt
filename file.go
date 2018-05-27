package main

import (
	"errors"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
)

func setupFileCommand(app *kingpin.Application) {
	fileCommand := app.Command("init-file", "Initialize a new empty file")

	fileCommand.Command("init", "Initialize a new empty file").Default().Default()
}

func handleFileCommand(commands []string) error {
	if len(commands) < 1 {
		return errors.New("No subcommand found for file command")
	}

	switch commands[1] {
	case "init":
		return handleInitFileCommand(commands)
	default:
		return errors.New(fmt.Sprintf("File subcommand not supported : %s", commands[1]))
	}
}

func handleInitFileCommand(commands []string) error {
	var teamPassFile TeamPassFile

	teamPassFile.Comment = *comment

	return writeFile(filename, true, teamPassFile)
}
