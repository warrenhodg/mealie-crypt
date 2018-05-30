package main

import (
	"errors"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path/filepath"
)

var groupsCommand *kingpin.CmdClause

var listGroupsCommand *kingpin.CmdClause

var addGroupCommand *kingpin.CmdClause
var groupName *string
var groupUserNames *[]string

var removeGroupCommand *kingpin.CmdClause

var groupAddUserCommand *kingpin.CmdClause
var groupsUsername *string
var groupsPrivateKeyFile *string

func setupGroupsCommand(app *kingpin.Application) {
	groupsCommand = app.Command("groups", "Manage groups")

	groupName = groupsCommand.Flag("group-name", "Name of group").Short('g').Default("_").String()
	groupUserNames = groupsCommand.Flag("users", "Names of users").Short('U').Default(os.Getenv("USER")).Strings()

	listGroupsCommand = groupsCommand.Command("list", "List groups")

	addGroupCommand = groupsCommand.Command("add", "Add a group")

	removeGroupCommand = groupsCommand.Command("remove", "Remove a group")

	groupAddUserCommand = groupsCommand.Command("add-user", "Add user to group")
	groupsUsername = groupsCommand.Flag("user", "Name of user").Short('u').Default(os.Getenv("USER")).String()
	groupsPrivateKeyFile = groupsCommand.Flag("pvt-key", "Filename of private key").Short('k').Default(filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa")).String()
}

func handleGroupsCommand(commands []string) error {
	if len(commands) < 2 {
		return errors.New("No subcommand found for groups command")
	}

	switch commands[1] {
	case "list":
		return handleListGroupsCommand(commands)

	case "add":
		return handleAddGroupCommand(commands)

	case "remove":
		return handleRemoveGroupCommand(commands)

	case "add-user":
		return handleGroupAddUserCommand(commands)

	default:
		return errors.New(fmt.Sprintf("Groups subcommand not supported : %s", commands[1]))
	}
}

func checkGroupname() error {
	return checkParam(*groupName, "^.+$", "Groupname must not be empty")
}

func checkGroupsUsername() error {
	return checkParam(*groupsUsername, "^.+$", "Username must not be empty")
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

	err = checkGroupname()
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

	err = checkGroupname()
	if err != nil {
		return err
	}

	//Remove group from groups
	delete(teamPassFile.Groups, *groupName)

	return writeFile(filename, false, teamPassFile)
}

func handleGroupAddUserCommand(commands []string) error {
	teamPassFile, err := readFile(filename, true)
	if err != nil {
		return err
	}

	err = checkGroupsUsername()
	if err != nil {
		return err
	}

	_, found := teamPassFile.Users[*groupsUsername]
	if !found {
		return errors.New(fmt.Sprintf("User not found : %s", *groupsUsername))
	}

	err = checkGroupname()
	if err != nil {
		return err
	}

	group, found := teamPassFile.Groups[*groupName]
	if !found {
		return errors.New(fmt.Sprintf("Group does not exist : %s", *groupName))
	}

	encSymKey, found := group.Keys[*groupsUsername]
	if !found {
		return errors.New(fmt.Sprintf("User not part of group : %s", *groupsUsername))
	}

	symKey, err := decryptSymmetricalKey(encSymKey, *groupsPrivateKeyFile)
	if err != nil {
		return err
	}

	for i := 0; i < len(*groupUserNames); i++ {
		username := (*groupUserNames)[i]
		user, found := teamPassFile.Users[username]
		if !found {
			return errors.New(fmt.Sprintf("User was not found : %s", username))
		}

		_, found = group.Keys[username]
		if !found {
			err := addEncryptedSymmetricKey(&group, symKey, username, user.PublicKey)
			if err != nil {
				return err
			}
		}
	}

	teamPassFile.Groups[*groupName] = group

	return writeFile(filename, false, teamPassFile)
}
