package main

import (
	"errors"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
)

var groupCommand *kingpin.CmdClause

var addGroupCommand *kingpin.CmdClause
var groupName *string
var groupKeyLen *int

func setupGroupCommand(app *kingpin.Application) {
	groupCommand = app.Command("group", "Add a group to the project")

	groupName = groupCommand.Flag("name", "Name of group").Short('n').Required().String()
	groupKeyLen = groupCommand.Flag("key-len", "Length of encryption key").Short('l').Default("32").Int()

	addGroupCommand = groupCommand.Command("add", "Add a group to the project")
}

func handleGroupCommand(commands []string) error {
	if len(commands) < 1 {
		return errors.New("No subcommand found for group command")
	}

	switch commands[1] {
	case "add":
		return handleAddGroupCommand(commands)
	default:
		return errors.New(fmt.Sprintf("Group subcommand not supported : %s", commands[1]))
	}
}

func handleAddGroupCommand(commands []string) error {
	var group TeamPassGroup

	teamPassFile, err := readFile(filename, true)
	if err != nil {
		return err
	}

	_, found := teamPassFile.Groups[*groupName]
	if found {
		return errors.New(fmt.Sprintf("Group already exists : %s", *groupName))
	}

	key, err := CreateKey(*groupKeyLen)
	if err != nil {
		return err
	}

	group.Keys = make(map[string]string)
	group.Keys[*groupName] = string(key)

	teamPassFile.Groups[*groupName] = group

	return writeFile(filename, false, teamPassFile)
}
