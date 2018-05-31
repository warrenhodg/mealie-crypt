package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"regexp"
)

var currentFileVersion = 1

type MealieCryptFile struct {
	Version int                         `yaml:"version"`
	Comment string                      `yaml:"comment,omitempty"`
	Users   map[string]MealieCryptUser  `yaml:"users,omitempty"`
	Groups  map[string]MealieCryptGroup `yaml:"groups,omitempty"`
}

type MealieCryptUser struct {
	PublicKey string `yaml:"public_key"`
	Comment   string `yaml:"comment,omitempty"`
}

type MealieCryptGroup struct {
	Keys      map[string]string `yaml:"keys"`
	Values    map[string]string `yaml:"values"`
	Decrypted map[string]string `yaml:"decrypted,omitempty"`
}

func (mealieCryptFile *MealieCryptFile) ensureMapsExist() {
	if mealieCryptFile.Users == nil {
		mealieCryptFile.Users = make(map[string]MealieCryptUser)
	}

	if mealieCryptFile.Groups == nil {
		mealieCryptFile.Groups = make(map[string]MealieCryptGroup)
	}

	for groupName, group := range mealieCryptFile.Groups {
		if group.Values == nil {
			group.Values = make(map[string]string)
			mealieCryptFile.Groups[groupName] = group
		}
	}
}

func readFile(filename *string, mustExist bool) (mealieCryptFile MealieCryptFile, err error) {
	var file *os.File

	defer mealieCryptFile.ensureMapsExist()

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

	err = yaml.Unmarshal(bytes, &mealieCryptFile)
	if err != nil {
		return
	}

	return
}

func writeFile(filename *string, mustNotExist bool, mealieCryptFile MealieCryptFile) (err error) {
	var file *os.File

	mealieCryptFile.Version = currentFileVersion

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

	bytes, err := yaml.Marshal(mealieCryptFile)
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
