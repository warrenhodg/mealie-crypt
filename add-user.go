package main

import (
    "gopkg.in/alecthomas/kingpin.v2"
    "os"
)

var addUserCommand *kingpin.CmdClause
var addUserName *string
var addUserKeyFile *string

func addAddUserCommand(app *kingpin.Application) {
    addUserCommand = app.Command("add-user", "Add a user to the project")
    addUserName = addUserCommand.Flag("name", "Name of user to add").Short('n').Default(os.Getenv("USER")).String()
    addUserKeyFile = addUserCommand.Flag("key-file", "Filename of user's public key to add").Short('k').Default(os.Getenv("HOME") + "/.ssh/id_rsa.pub").String()
}

func addUser() error {
    var user TeamPassUser

    teamPassFile, err := readFile(filename, true)
    if err != nil {
        return err
    }

    keyValue, err := readKey(addUserKeyFile)
    if err != nil {
        return err
    }

    user.Name = *addUserName
    user.Value = keyValue
    user.Comment = *comment

    teamPassFile.Users = append(teamPassFile.Users, user)

    return writeFile(filename, false, teamPassFile)
}
