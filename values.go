package main

import (
	"errors"
	"fmt"
	"github.com/gobwas/glob"
	"gopkg.in/alecthomas/kingpin.v2"
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
var removeValueCommand *kingpin.CmdClause

func setupValuesCommand(app *kingpin.Application) {
	valuesCommand = app.Command("values", "Manage values")

	valuesGroup = valuesCommand.Flag("group", "Name of group").Short('g').Default(configDefaults[keyGroupName]).String()
	valuesUsername = valuesCommand.Flag("user", "Name of user").Short('u').Default(configDefaults[keyUsername]).String()
	valuesPrivateKeyFile = valuesCommand.Flag("pvt-key", "Filename of private key").Short('k').Default(configDefaults[keyPrivateKeyFile]).String()
	valuesName = valuesCommand.Flag("name", "Name of value").Short('n').String()
	valuesValue = valuesCommand.Flag("value", "Value").Short('v').String()

	listValuesCommand = valuesCommand.Command("list", "List values")

	setValueCommand = valuesCommand.Command("set", "Set value")

	getValueCommand = valuesCommand.Command("get", "Get value")

	removeValueCommand = valuesCommand.Command("remove", "Remove value")
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

	case "remove":
		return handleRemoveValueCommand(commands)

	case "list":
		return handleListValuesCommand(commands)
	default:
		return errors.New(fmt.Sprintf("Values subcommand not supported : %s", commands[1]))
	}
}

func handleListValuesCommand(commands []string) error {
	var group MealieCryptGroup

	mealieCryptFile, err := readFile(filename, true)
	if err != nil {
		return err
	}

	group, found := mealieCryptFile.Groups[*valuesGroup]
	if !found {
		return errors.New(fmt.Sprintf("Group not found : %s", valuesGroup))
	}

	for valueName, _ := range group.Values {
		fmt.Printf("%s\n", valueName)
	}

	return nil
}

func handleGetValueCommand(commands []string) error {
	mealieCryptFile, err := readFile(filename, true)
	if err != nil {
		return err
	}

	_, found := mealieCryptFile.Users[*valuesUsername]
	if !found {
		return errors.New(fmt.Sprintf("User not found : %s", *valuesUsername))
	}

	if *valuesName == "" {
		*valuesName = "*"
	}

	for groupName, group := range mealieCryptFile.Groups {
		encSymKey, found := group.Keys[*valuesUsername]
		if !found {
			continue
		}

		pvtKey, err := readPrivateKey(*valuesPrivateKeyFile)
		if err != nil {
			return err
		}

		symKey, err := decryptSymmetricalKey(encSymKey, pvtKey)
		if err != nil {
			return err
		}

		g, err := glob.Compile(*valuesName)
		if err != nil {
			return err
		}

		for encValueName, encValue := range group.Values {
			valueName, err := decryptValue(symKey, encValueName)
			if err != nil {
				return err
			}

			if g.Match(valueName) {
				decValue, err := decryptValue(symKey, encValue)
				if err != nil {
					return err
				}

				fmt.Printf("%s : %s = %s\n", groupName, valueName, decValue)
			}
		}
	}

	return nil
}

func handleSetValueCommand(commands []string) error {
	mealieCryptFile, err := readFile(filename, true)
	if err != nil {
		return err
	}

	_, found := mealieCryptFile.Users[*valuesUsername]
	if !found {
		return errors.New(fmt.Sprintf("User not found : %s", valuesUsername))
	}

	group, found := mealieCryptFile.Groups[*valuesGroup]
	if !found {
		return errors.New(fmt.Sprintf("Group not found : %s", valuesGroup))
	}

	encSymKey, found := group.Keys[*valuesUsername]
	if !found {
		return errors.New(fmt.Sprintf("User not part of group : %s", *valuesUsername))
	}

	pvtKey, err := readPrivateKey(*valuesPrivateKeyFile)
	if err != nil {
		return err
	}

	symKey, err := decryptSymmetricalKey(encSymKey, pvtKey)
	if err != nil {
		return err
	}

	encValueName, encValue, err := findEncValue(symKey, &group, *valuesName)
	if err != nil {
		return err
	}

	doAdd := true

	if encValueName != "" {
		decValue, err := decryptValue(symKey, encValue)
		if err != nil {
			return err
		}

		if decValue == *valuesValue {
			doAdd = false
		}
	} else {
		encValueName, err = encryptValue(symKey, *valuesName)
		if err != nil {
			return err
		}
	}

	if doAdd {
		encValue, err := encryptValue(symKey, *valuesValue)
		if err != nil {
			return err
		}

		group.Values[encValueName] = encValue

		mealieCryptFile.Groups[*valuesGroup] = group
	}

	return writeFile(filename, false, mealieCryptFile)
}

func handleRemoveValueCommand(commands []string) error {
	mealieCryptFile, err := readFile(filename, true)
	if err != nil {
		return err
	}

	_, found := mealieCryptFile.Users[*valuesUsername]
	if !found {
		return errors.New(fmt.Sprintf("User not found : %s", valuesUsername))
	}

	group, found := mealieCryptFile.Groups[*valuesGroup]
	if !found {
		return errors.New(fmt.Sprintf("Group not found : %s", valuesGroup))
	}

	encSymKey, found := group.Keys[*valuesUsername]
	if !found {
		return errors.New(fmt.Sprintf("User not part of group : %s", *valuesUsername))
	}

	pvtKey, err := readPrivateKey(*valuesPrivateKeyFile)
	if err != nil {
		return err
	}

	symKey, err := decryptSymmetricalKey(encSymKey, pvtKey)
	if err != nil {
		return err
	}

	encValueName, _, err := findEncValue(symKey, &group, *valuesName)
	if err != nil {
		return err
	}

	if encValueName != "" {
		delete(group.Values, encValueName)
		mealieCryptFile.Groups[*valuesGroup] = group
	}

	return writeFile(filename, false, mealieCryptFile)
}

func findEncValue(symKey string, group *MealieCryptGroup, valueName string) (encValueName string, encValue string, err error) {
	for evn, ev := range group.Values {
		var decValueName string

		decValueName, err = decryptValue(symKey, evn)
		if err != nil {
			return
		}

		if decValueName == valueName {
			encValueName = evn
			encValue = ev
			return
		}
	}

	return
}
