package main

import (
	"errors"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var groupsCommand *kingpin.CmdClause

var listGroupsCommand *kingpin.CmdClause

var addGroupCommand *kingpin.CmdClause
var groupName *string
var groupUserNames *[]string

var removeGroupCommand *kingpin.CmdClause

func setupGroupsCommand(app *kingpin.Application) {
	groupsCommand = app.Command("groups", "Manage groups")

	groupName = groupCommand.Flag("group-name", "Name of group").Short('g').Default("_").String()
	groupUserNames = groupCommand.Flag("users", "Names of users").Short('u').Default(os.Getenv("USER")).Strings()

	listGroupsCommand = groupCommand.Command("list", "List groups")

	addGroupCommand = groupCommand.Command("add", "Add a group")

	removeGroupCommand = groupCommand.Command("remove", "Remove a group")
}

func handleGroupCommand(commands []string) error {
	if len(commands) < 1 {
		return errors.New("No subcommand found for groups command")
	}

	switch commands[1] {
	case "list":
		return handleListGroupsCommand(commands)
	case "add":
		return handleAddGroupCommand(commands)
	case "remove":
		return handleRemoveGroupCommand(commands)
	default:
		return errors.New(fmt.Sprintf("Groups subcommand not supported : %s", commands[1]))
	}
}

func handleListGroupsCommand(commands []string) error {
	teamPassFile, err := readFile(filename, true)
	if err != nil {
		return err
	}

	for groupname, _ := range teamPassFile.Groups {
		fmt.Printf("%s\n", groupname)
	}

	return nil
}

func addEncryptedSymmetricKey(group *TeamPassGroup, symKey string, userName string, publicKey string) error {
	encSymKey, err := encryptSymmetricalKey(symKey, publicKey)
	if err != nil {
		return err
	}

	group.Keys[userName] = encSymKey
	return nil
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

	symKey, err := createSymmetricalKey()
	if err != nil {
		return err
	}

	//Create a key for each user in the list
	group.Keys = make(map[string]string)
	for i := 0; i < len(*groupUserNames); i++ {
		username := (*groupUserNames)[i]
		user, found := teamPassFile.Users[username]
		if !found {
			return errors.New(fmt.Sprintf("User was not found : %s", username))
		}

		err = addEncryptedSymmetricKey(&group, symKey, username, user.PublicKey)
		if err != nil {
			return err
		}
	}

	teamPassFile.Groups[*groupName] = group

	return writeFile(filename, false, teamPassFile)
}

func handleRemoveGroupCommand(commands []string) error {
	teamPassFile, err := readFile(filename, true)
	if err != nil {
		return err
	}

	//Remove group from groups
	delete(teamPassFile.Groups, *groupName)

	return writeFile(filename, false, teamPassFile)
}
