package main

import (
	"errors"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
)

var addGroupCommand *kingpin.CmdClause
var addGroupName *string
var addGroupKeyLen *int

func addAddGroupCommand(app *kingpin.Application) {
	addGroupCommand = app.Command("add-group", "Add a group to the project")
	addGroupName = addGroupCommand.Flag("name", "Name of group to add").Short('n').Required().String()
	addGroupKeyLen = addGroupCommand.Flag("key-len", "Length of encryption key").Short('l').Default("32").Int()
}

func checkIfGroupExists(teamPassFile *TeamPassFile, name string) error {
	for i := 0; i < len(teamPassFile.Groups); i++ {
		if teamPassFile.Groups[i].Name == name {
			return errors.New(fmt.Sprintf("Group already exists : %s", name))
		}
	}

	return nil
}

func addGroup() error {
	var group TeamPassGroup

	teamPassFile, err := readFile(filename, true)
	if err != nil {
		return err
	}

	err = checkIfGroupExists(&teamPassFile, *addGroupName)
	if err != nil {
		return err
	}

	key, err := CreateKey(*addGroupKeyLen)
	if err != nil {
		return err
	}

	group.Name = *addGroupName
	group.Keys = make(map[string]string)
	group.Keys[*addGroupName] = string(key)

	teamPassFile.Groups = append(teamPassFile.Groups, group)

	return writeFile(filename, false, teamPassFile)
}
