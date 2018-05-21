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

    filename := app.Flag("file", "Name of file to manage").Short('f').String()
    comment := app.Flag("comment", "A comment").Short('c').String()

    app.Command("init-file", "Initialize a new empty file")

    addPubKeyCommand := app.Command("add-pub-key", "Add a public key to the project")
    addPubKeyAlias := addPubKeyCommand.Flag("alias", "Alias for key to add").Short('a').Required().String()
    addPubKeyFile := addPubKeyCommand.Flag("pub-key", "Filename of public key to add").Short('k').Required().String()

    command, err := app.Parse(os.Args[1:])
    if err != nil {
        fmt.Errorf("%s", err.Error()) 
        os.Exit(1)
    }

    switch(command) {
    case "init-file":
        err = initFile(filename, comment)
        if err != nil {
            fmt.Errorf("%s", err.Error())
            os.Exit(1)
        }

    case "add-pub-key":
        err = addPubKey(filename, addPubKeyAlias, addPubKeyFile, comment)
        if err != nil {
            fmt.Errorf("%s", err.Error())
            os.Exit(1)
        }
    }
}
