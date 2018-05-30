package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"regexp"
)

type TeamPassFile struct {
	Comment string                   `yaml:"comment"`
	Users   map[string]TeamPassUser  `yaml:"users"`
	Groups  map[string]TeamPassGroup `yaml:"groups"`
}

type TeamPassUser struct {
	PublicKey string `yaml:"public_key"`
	Comment   string `yaml:"comment"`
}

type TeamPassGroup struct {
	Keys      map[string]string `yaml:"keys"`
	Values    map[string]string `yaml:"values"`
	Decrypted map[string]string `yaml:"decrypted"`
}

func (teamPassFile *TeamPassFile) ensureMapsExist() {
	if teamPassFile.Users == nil {
		teamPassFile.Users = make(map[string]TeamPassUser)
	}

	if teamPassFile.Groups == nil {
		teamPassFile.Groups = make(map[string]TeamPassGroup)
	}

	for groupName, group := range teamPassFile.Groups {
		if group.Values == nil {
			group.Values = make(map[string]string)
			teamPassFile.Groups[groupName] = group
		}
	}
}

func readFile(filename *string, mustExist bool) (teamPassFile TeamPassFile, err error) {
	var file *os.File

	defer teamPassFile.ensureMapsExist()

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

func checkParam(param string, regex string, message string) error {
	re, err := regexp.Compile(regex)
	if err != nil {
		return err
	}

	if !re.MatchString(param) {
		return errors.New(fmt.Sprintf(message))
	}

	return nil
}
