package main

import (
	"errors"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path/filepath"
)

var encryptCommand *kingpin.CmdClause
var encryptUsername *string
var encryptPrivateKeyFile *string

func setupEncryptCommand(app *kingpin.Application) {
	encryptCommand := app.Command("encrypt", "Encrypt the encryptable parts of the file")

	encryptUsername = encryptCommand.Flag("user", "Name of user").Short('u').Default(os.Getenv(userVar)).String()
	encryptPrivateKeyFile = encryptCommand.Flag("pvt-key", "Filename of private key").Short('k').Default(filepath.Join(os.Getenv(homeVar), ".ssh", "id_rsa")).String()
}

func handleEncryptCommand(commands []string) error {
	mealieCryptFile, err := readFile(filename, true)
	if err != nil {
		return err
	}

	_, found := mealieCryptFile.Users[*encryptUsername]
	if !found {
		return errors.New(fmt.Sprintf("User not found : %s", *encryptUsername))
	}

	pvtKey, err := readPrivateKey(*encryptPrivateKeyFile)
	if err != nil {
		return err
	}

	for groupName, group := range mealieCryptFile.Groups {
		encSymKey, found := group.Keys[*encryptUsername]
		if found {
			if group.Decrypted == nil {
				group.Decrypted = make(map[string]string)
			}

			symKey, err := decryptSymmetricalKey(encSymKey, pvtKey)
			if err != nil {
				return err
			}

			for valueName, value := range group.Decrypted {
				mustAdd := true

				encValueName, encValue, err := findEncValue(symKey, &group, valueName)
				if err != nil {
					return err
				}

				if encValueName != "" {
					decValue, err := decryptValue(symKey, encValue)
					if err != nil {
						return err
					}

					if decValue == value {
						mustAdd = false
					}
				}

				if mustAdd {
					if encValueName == "" {
						encValueName, err = encryptValue(symKey, valueName)
						if err != nil {
							return err
						}
					}

					newEncValue, err := encryptValue(symKey, value)
					if err != nil {
						return err
					}

					group.Values[encValueName] = newEncValue
				}
			}

			group.Decrypted = nil

			mealieCryptFile.Groups[groupName] = group
		} else {
			if group.Decrypted != nil {
				if len(group.Decrypted) > 0 {
					return errors.New(fmt.Sprintf("There are plain-text values in a group (%s) of which you are not a member", groupName))
				} else {
					group.Decrypted = nil

					mealieCryptFile.Groups[groupName] = group
				}
			}
		}
	}

	return writeFile(filename, false, mealieCryptFile)
}
