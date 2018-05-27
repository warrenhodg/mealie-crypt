package main

import (
	"errors"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var userCommand *kingpin.CmdClause

var addUserCommand *kingpin.CmdClause
var userName *string
var userKeyFile *string

var removeUserCommand *kingpin.CmdClause

var listUsersCommand *kingpin.CmdClause

func setupUserCommand(app *kingpin.Application) {
	userCommand = app.Command("user", "Manage users")

	userName = userCommand.Flag("name", "Name of user").Short('u').Default(os.Getenv("USER")).String()
	userKeyFile = userCommand.Flag("key-file", "Filename of user's public key").Short('k').Default(os.Getenv("HOME") + "/.ssh/id_rsa.pub").String()

	addUserCommand = userCommand.Command("list", "List users")

	addUserCommand = userCommand.Command("add", "Add a user")

	removeUserCommand = userCommand.Command("remove", "Remove a user")
}

func handleUserCommand(commands []string) error {
	if len(commands) < 1 {
		return errors.New("No subcommand found for user command")
	}

	switch commands[1] {
	case "list":
		return handleListUsersCommand(commands)
	case "add":
		return handleAddUserCommand(commands)
	case "remove":
		return handleRemoveUserCommand(commands)
	default:
		return errors.New(fmt.Sprintf("User subcommand not supported : %s", commands[1]))
	}
}

func handleListUsersCommand(commands []string) error {
	teamPassFile, err := readFile(filename, true)
	if err != nil {
		return err
	}

	for username, _ := range teamPassFile.Users {
		fmt.Printf("%s\n", username)
	}

	return nil
}

func handleAddUserCommand(commands []string) error {
	var user TeamPassUser

	teamPassFile, err := readFile(filename, true)
	if err != nil {
		return err
	}

	keyValue, err := readPublicKey(userKeyFile)
	if err != nil {
		return err
	}

	_, found := teamPassFile.Users[*userName]
	if found {
		return errors.New(fmt.Sprintf("User already exists : %s", *userName))
	}

	user.PublicKey = keyValue
	user.Comment = *comment

	teamPassFile.Users[*userName] = user

	return writeFile(filename, false, teamPassFile)
}

func handleRemoveUserCommand(commands []string) error {
	teamPassFile, err := readFile(filename, true)
	if err != nil {
		return err
	}

	//Remove user from users
	delete(teamPassFile.Users, *userName)

	//Remove user from groups
	for groupName, _ := range teamPassFile.Groups {
		delete(teamPassFile.Groups[groupName].Keys, *userName)
		if len(teamPassFile.Groups[groupName].Keys) == 0 {
			return errors.New(fmt.Sprintf("Cannot remove last user from group : %s", groupName))
		}
	}

	return writeFile(filename, false, teamPassFile)
}
