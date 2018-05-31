package main

import (
	"errors"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path/filepath"
)

var decryptCommand *kingpin.CmdClause
var decryptUsername *string
var decryptPrivateKeyFile *string

func setupDecryptCommand(app *kingpin.Application) {
	decryptCommand := app.Command("decrypt", "Decrypt the decryptable parts of the file")

	decryptUsername = decryptCommand.Flag("user", "Name of user").Short('u').Default(os.Getenv("USER")).String()
	decryptPrivateKeyFile = decryptCommand.Flag("pvt-key", "Filename of private key").Short('k').Default(filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa")).String()
}

func handleDecryptCommand(commands []string) error {
	mealieCryptFile, err := readFile(filename, true)
	if err != nil {
		return err
	}

	_, found := mealieCryptFile.Users[*decryptUsername]
	if !found {
		return errors.New(fmt.Sprintf("User not found : %s", *decryptUsername))
	}

	for groupName, group := range mealieCryptFile.Groups {
		encSymKey, found := group.Keys[*decryptUsername]
		if found {
			if group.Decrypted == nil {
				group.Decrypted = make(map[string]string)
			}

			symKey, err := decryptSymmetricalKey(encSymKey, *decryptPrivateKeyFile)
			if err != nil {
				return err
			}

			for valueName, encValue := range group.Values {
				decValue, err := decryptValue(symKey, encValue)
				if err != nil {
					return err
				}

				_, found := group.Decrypted[valueName]
				if !found {
					group.Decrypted[valueName] = decValue
				}
			}

			mealieCryptFile.Groups[groupName] = group
		}
	}

	return writeFile(filename, false, mealieCryptFile)
}
