package main

import (
	"errors"
	"fmt"
	"github.com/gobwas/glob"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var valuesCommand *kingpin.CmdClause
var valuesGroup *string
var valuesUsername *string
var valuesPrivateKeyFile *string
var valuesName *string
var valuesValue *string
var listValuesCommand *kingpin.CmdClause
var setValueCommand *kingpin.CmdClause
var getValueCommand *kingpin.CmdClause

func setupValuesCommand(app *kingpin.Application) {
	valuesCommand = app.Command("values", "Manage values")

	valuesGroup = valuesCommand.Flag("group", "Name of group").Short('g').Default("_").String()
	valuesUsername = valuesCommand.Flag("user", "Name of user").Short('u').Default(os.Getenv("USER")).String()
	valuesPrivateKeyFile = valuesCommand.Flag("pvt-key", "Filename of private key").Short('k').Default(os.Getenv("HOME") + "/.ssh/id_rsa").String()
	valuesName = valuesCommand.Flag("name", "Name of value").Short('n').String()
	valuesValue = valuesCommand.Flag("value", "Value").Short('v').String()

	listValuesCommand = valuesCommand.Command("list", "List values")

	setValueCommand = valuesCommand.Command("set", "Set value")

	getValueCommand = valuesCommand.Command("get", "Get value")
}

func handleValuesCommand(commands []string) error {
	if len(commands) < 2 {
		return errors.New("No subcommand found for values command")
	}

	switch commands[1] {
	case "set":
		return handleSetValueCommand(commands)

	case "get":
		return handleGetValueCommand(commands)

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

	return nil
}

func handleGetValueCommand(commands []string) error {
	teamPassFile, err := readFile(filename, true)
	if err != nil {
		return err
	}

	_, found := teamPassFile.Users[*valuesUsername]
	if !found {
		return errors.New(fmt.Sprintf("User not found : %s", valuesUsername))
	}

	group, found := teamPassFile.Groups[*valuesGroup]
	if !found {
		return errors.New(fmt.Sprintf("Group not found : %s", valuesGroup))
	}

	encSymKey, found := group.Keys[*valuesUsername]
	if !found {
		return errors.New(fmt.Sprintf("User not part of group : %s", *valuesUsername))
	}

	symKey, err := decryptSymmetricalKey(encSymKey, *valuesPrivateKeyFile)
	if err != nil {
		return err
	}

	if *valuesName == "" {
		*valuesName = "*"
	}

	g, err := glob.Compile(*valuesName)
	if err != nil {
		return err
	}

	for valueName, encValue := range group.Values {
		if g.Match(valueName) {
			decValue, err := decryptValue(symKey, encValue)
			if err != nil {
				return err
			}

			fmt.Printf("%s : %s\n", valueName, decValue)
		}
	}

	return nil
}

func handleSetValueCommand(commands []string) error {
	teamPassFile, err := readFile(filename, true)
	if err != nil {
		return err
	}

	_, found := teamPassFile.Users[*valuesUsername]
	if !found {
		return errors.New(fmt.Sprintf("User not found : %s", valuesUsername))
	}

	group, found := teamPassFile.Groups[*valuesGroup]
	if !found {
		return errors.New(fmt.Sprintf("Group not found : %s", valuesGroup))
	}

	encSymKey, found := group.Keys[*valuesUsername]
	if !found {
		return errors.New(fmt.Sprintf("User not part of group : %s", *valuesUsername))
	}

	symKey, err := decryptSymmetricalKey(encSymKey, *valuesPrivateKeyFile)
	if err != nil {
		return err
	}

	encValue, err := encryptValue(symKey, *valuesValue)
	if err != nil {
		return err
	}

	group.Values[*valuesName] = encValue

	teamPassFile.Groups[*valuesGroup] = group

	return writeFile(filename, false, teamPassFile)
}
