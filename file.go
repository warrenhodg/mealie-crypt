package main

import (
	"errors"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
)

func setupFileCommand(app *kingpin.Application) {
	fileCommand := app.Command("file", "Initialize a new empty file")

	fileCommand.Command("init", "Initialize a new empty file").Default()

	fileCommand.Command("touch", "Simply load and save the file")
}

func handleFileCommand(commands []string) error {
	if len(commands) < 1 {
		return errors.New("No subcommand found for file command")
	}

	switch commands[1] {
	case "init":
		return handleInitFileCommand(commands)

	case "touch":
		return handleTouchFileCommand(commands)

	default:
		return errors.New(fmt.Sprintf("File subcommand not supported : %s", commands[1]))
	}
}

func handleInitFileCommand(commands []string) error {
	var dioscoreaFile DioscoreaFile

	dioscoreaFile.Comment = *comment

	return writeFile(filename, true, dioscoreaFile)
}

func handleTouchFileCommand(commands []string) error {
	dioscoreaFile, err := readFile(filename, true)
	if err != nil {
		return err
	}

	return writeFile(filename, false, dioscoreaFile)
}
