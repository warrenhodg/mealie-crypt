package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"regexp"
)

type DioscoreaFile struct {
	Comment string                    `yaml:"comment,omitempty"`
	Users   map[string]DioscoreaUser  `yaml:"users,omitempty"`
	Groups  map[string]DioscoreaGroup `yaml:"groups,omitempty"`
}

type DioscoreaUser struct {
	PublicKey string `yaml:"public_key"`
	Comment   string `yaml:"comment,omitempty"`
}

type DioscoreaGroup struct {
	Keys      map[string]string `yaml:"keys"`
	Values    map[string]string `yaml:"values"`
	Decrypted map[string]string `yaml:"decrypted,omitempty"`
}

func (dioscoreaFile *DioscoreaFile) ensureMapsExist() {
	if dioscoreaFile.Users == nil {
		dioscoreaFile.Users = make(map[string]DioscoreaUser)
	}

	if dioscoreaFile.Groups == nil {
		dioscoreaFile.Groups = make(map[string]DioscoreaGroup)
	}

	for groupName, group := range dioscoreaFile.Groups {
		if group.Values == nil {
			group.Values = make(map[string]string)
			dioscoreaFile.Groups[groupName] = group
		}
	}
}

func readFile(filename *string, mustExist bool) (dioscoreaFile DioscoreaFile, err error) {
	var file *os.File

	defer dioscoreaFile.ensureMapsExist()

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

	err = yaml.Unmarshal(bytes, &dioscoreaFile)
	if err != nil {
		return
	}

	return
}

func writeFile(filename *string, mustNotExist bool, dioscoreaFile DioscoreaFile) (err error) {
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

	bytes, err := yaml.Marshal(dioscoreaFile)
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
