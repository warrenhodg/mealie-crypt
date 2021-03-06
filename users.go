package main

import (
	"errors"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
)

var usersCommand *kingpin.CmdClause

var addUserCommand *kingpin.CmdClause
var userName *string
var userKeyFile *string

var removeUserCommand *kingpin.CmdClause

var listUsersCommand *kingpin.CmdClause

func setupUsersCommand(app *kingpin.Application) {
	usersCommand = app.Command("users", "Manage users")

	userName = usersCommand.Flag("name", "Name of user").Short('u').Default(configDefaults[keyUsername]).String()
	userKeyFile = usersCommand.Flag("key-file", "Filename of user's public key").Short('k').Default(configDefaults[keyPublicKeyFile]).String()

	addUserCommand = usersCommand.Command("list", "List users")

	addUserCommand = usersCommand.Command("add", "Add a user")

	removeUserCommand = usersCommand.Command("remove", "Remove a user")
}

func handleUsersCommand(commands []string) error {
	if len(commands) < 2 {
		return errors.New("No subcommand found for users command")
	}

	switch commands[1] {
	case "list":
		return handleListUsersCommand(commands)
	case "add":
		return handleAddUserCommand(commands)
	case "remove":
		return handleRemoveUserCommand(commands)
	default:
		return errors.New(fmt.Sprintf("Users subcommand not supported : %s", commands[1]))
	}
}

func handleListUsersCommand(commands []string) error {
	mealieCryptFile, err := readFile(filename, true)
	if err != nil {
		return err
	}

	for username, _ := range mealieCryptFile.Users {
		fmt.Printf("%s\n", username)
	}

	return nil
}

func checkUsername() error {
	return checkParam(*userName, "^.+$", "Username must not be empty")
}

func handleAddUserCommand(commands []string) error {
	var user MealieCryptUser

	mealieCryptFile, err := readFile(filename, true)
	if err != nil {
		return err
	}

	err = checkUsername()
	if err != nil {
		return err
	}

	keyValue, err := readPublicKey(userKeyFile)
	if err != nil {
		return err
	}

	_, found := mealieCryptFile.Users[*userName]
	if found {
		return errors.New(fmt.Sprintf("User already exists : %s", *userName))
	}

	user.PublicKey = keyValue
	user.Comment = *comment

	mealieCryptFile.Users[*userName] = user

	return writeFile(filename, false, mealieCryptFile)
}

func handleRemoveUserCommand(commands []string) error {
	mealieCryptFile, err := readFile(filename, true)
	if err != nil {
		return err
	}

	err = checkUsername()
	if err != nil {
		return err
	}

	//Remove user from users
	delete(mealieCryptFile.Users, *userName)

	//Remove user from groups
	for groupName, _ := range mealieCryptFile.Groups {
		delete(mealieCryptFile.Groups[groupName].Keys, *userName)
		if len(mealieCryptFile.Groups[groupName].Keys) == 0 {
			return errors.New(fmt.Sprintf("Cannot remove last user from group : %s", groupName))
		}
	}

	return writeFile(filename, false, mealieCryptFile)
}
