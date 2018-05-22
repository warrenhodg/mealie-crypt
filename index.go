package main

import (
    "fmt"
    "gopkg.in/alecthomas/kingpin.v2"
    "os"
)

var appName = "teampass"
var appDescription = "Utility for teams to manage sensitive information"
var version = "1.0.0"

var filename *string
var comment *string

func addGlobalFlags(app *kingpin.Application) {
    filename = app.Flag("file", "Name of file to manage").Short('f').Default("teampass.yaml").String()
    comment = app.Flag("comment", "A comment").Short('c').String()
}

func main() {
    app := kingpin.New(appName, appDescription)
    app.Version(version)

    addGlobalFlags(app)
    addInitFileCommand(app)
    addAddUserCommand(app)

    err := func() error {
        command, err := app.Parse(os.Args[1:])
        if err != nil {
            return err
        }

        switch(command) {
        case "init-file":
            return initFile()

        case "add-user":
            return addUser()
        }

        return nil
    }()

    if err != nil {
        fmt.Printf("%s", err.Error())
        os.Exit(1)
    }
}
