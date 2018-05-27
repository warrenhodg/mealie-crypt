package main

import (
	"errors"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
)

var valuesCommand *kingpin.CmdClause
var valuesGroup *string
var listValuesCommand *kingpin.CmdClause

func setupValuesCommand(app *kingpin.Application) {
	valuesCommand = app.Command("values", "Manage values")

	valuesGroup = valuesCommand.Flag("group", "Name of group").Short('g').Default("_").String()

	listValuesCommand = valuesCommand.Command("list", "List values")
}

func handleValuesCommand(commands []string) error {
	if len(commands) < 2 {
		return errors.New("No subcommand found for values command")
	}

	switch commands[1] {
	case "list":
		return handleListValuesCommand(commands)
	default:
		return errors.New(fmt.Sprintf("Values subcommand not supported : %s", commands[1]))
	}
}

func handleListValuesCommand(commands []string) error {
	var group TeamPassGroup

	teamPassFile, err := readFile(filename, true)
	if err != nil {
		return err
	}

	group, found := teamPassFile.Groups[*valuesGroup]
	if !found {
		return errors.New(fmt.Sprintf("Group not found : %s", valuesGroup))
	}

	for valueName, _ := range group.Values {
		fmt.Printf("%s\n", valueName)
	}

	return writeFile(filename, false, teamPassFile)
}
