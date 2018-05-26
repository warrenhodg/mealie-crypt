package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type TeamPassFile struct {
	Comment string          `yaml:"comment"`
	Users   []TeamPassUser  `yaml:"users"`
	Groups  []TeamPassGroup `yaml:"groups"`
}

type TeamPassUser struct {
	Name    string `yaml:"name"`
	Value   string `yaml:"value"`
	Comment string `yaml:"comment"`
}

type TeamPassGroup struct {
	Name string            `yaml:"name"`
	Keys map[string]string `yaml:"keys"`
}

func readFile(filename *string, mustExist bool) (teamPassFile TeamPassFile, err error) {
	var file *os.File

	_, err = os.Stat(*filename)
	if os.IsNotExist(err) {
		err = nil
		return
	}

	if *filename == "-" {
		file = os.Stdin
	} else {
		if mustExist {
			_, err = os.Stat(*filename)
			if os.IsNotExist(err) {
				err = errors.New(fmt.Sprintf("File does not exist : %s", *filename))
				return
			}
		}

		file, err = os.Open(*filename)
		if err != nil {
			return
		}
	}

	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(bytes, &teamPassFile)
	if err != nil {
		return
	}

	return
}

func writeFile(filename *string, mustNotExist bool, teamPassFile TeamPassFile) (err error) {
	var file *os.File

	if *filename == "-" {
		file = os.Stdout
	} else {
		if mustNotExist {
			_, err = os.Stat(*filename)
			if !os.IsNotExist(err) {
				err = errors.New(fmt.Sprintf("File already exists : %s", *filename))
				return
			}
		}

		file, err = os.Create(*filename)
		if err != nil {
			return
		}
	}

	defer file.Close()

	bytes, err := yaml.Marshal(teamPassFile)
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)
	if err != nil {
		return err
	}

	return
}
