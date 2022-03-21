package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strings"
)

var appName = "mealie-crypt"
var appDescription = "Utility for teams to manage sensitive information"
var version = "1.2.4"
var copyrightYear = 2018
var copyrightHolder = "Warren Hodgkinson"

var filename *string
var comment *string

func addGlobalFlags(app *kingpin.Application) {
	filename = app.Flag("file", "Name of file to manage").Short('f').Default(configDefaults[keyFilename]).String()
	comment = app.Flag("comment", "A comment").Short('c').String()
}

func main() {
	initConfigDefaults()

	app := kingpin.New(appName, appDescription)
	app.Version(version)

	addGlobalFlags(app)
	setupLicenseCommand(app)
	setupFileCommand(app)
	setupUsersCommand(app)
	setupGroupsCommand(app)
	setupValuesCommand(app)
	setupDecryptCommand(app)
	setupEncryptCommand(app)

	err := func() error {
		fullCommand, err := app.Parse(os.Args[1:])
		if err != nil {
			return err
		}

		commands := strings.Split(fullCommand, " ")

		switch commands[0] {
		case "decrypt":
			return handleDecryptCommand(commands)

		case "encrypt":
			return handleEncryptCommand(commands)

		case "file":
			return handleFileCommand(commands)

		case "groups":
			return handleGroupsCommand(commands)

		case "license":
			return handleLicenseCommand(commands)

		case "users":
			return handleUsersCommand(commands)

		case "values":
			return handleValuesCommand(commands)
		}

		return nil
	}()

	if err != nil {
		fmt.Printf("%s", err.Error())
		os.Exit(1)
	}
}
