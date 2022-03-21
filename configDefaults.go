package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

var configDefaults map[string]string
var mealierc = ".mealierc"

const keyFilename = "filename"
const keyUsername = "username"
const keyPublicKeyFile = "public_key_file"
const keyPrivateKeyFile = "private_key_file"
const keyGroupName = "group_name"

func initConfigDefaults() {
	configDefaults = make(map[string]string)

	configDefaults[keyFilename] = "mealie-crypt.yaml"
	configDefaults[keyUsername] = os.Getenv(userVar)
	configDefaults[keyPublicKeyFile] = filepath.Join(os.Getenv(homeVar), ".ssh", "id_rsa.pub")
	configDefaults[keyPrivateKeyFile] = filepath.Join(os.Getenv(homeVar), ".ssh", "id_rsa")
	configDefaults[keyGroupName] = "_"

	initMealierc()
}

func initMealierc() {
	mealiercFilename := mealierc
	_, err := os.Stat(mealiercFilename)
	if os.IsNotExist(err) {
		mealiercFilename = filepath.Join(os.Getenv(homeVar), mealierc)
		_, err := os.Stat(mealiercFilename)
		if os.IsNotExist(err) {
			return
		}
	}

	file, err := os.Open(mealiercFilename)
	if err != nil {
		return
	}

	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	var rcValues map[string]string
	err = yaml.Unmarshal(bytes, &rcValues)
	if err != nil {
		return
	}

	for key, value := range rcValues {
		configDefaults[key] = value
	}
}
