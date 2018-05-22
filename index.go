package main

import (
    "fmt"
    "gopkg.in/alecthomas/kingpin.v2"
    "os"
)

var appName = "teampass"
var appDescription = "Utility for teams to manage sensitive information"
var version = "1.0.0"

func main() {
    app := kingpin.New(appName, appDescription)
    app.Version(version)

    filename := app.Flag("file", "Name of file to manage").Short('f').Default("teampass.yaml").String()
    comment := app.Flag("comment", "A comment").Short('c').String()

    app.Command("init-file", "Initialize a new empty file")

    addUserCommand := app.Command("add-user", "Add a user to the project")
    addUserName := addUserCommand.Flag("name", "Name for key to add").Short('n').Required().String()
    addUserFile := addUserCommand.Flag("pub-key", "Filename of user's public key to add").Short('k').Default(os.Getenv("HOME") + "/.ssh/id_rsa.pub").String()

    err := func() error {
        command, err := app.Parse(os.Args[1:])
        if err != nil {
            return err
        }

        switch(command) {
        case "init-file":
            return initFile(filename, comment)

        case "add-user":
            return addUser(filename, addUserName, addUserFile, comment)
        }

        return nil
    }()

    if err != nil {
        fmt.Printf("%s", err.Error())
        os.Exit(1)
    }
}
