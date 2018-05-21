package main

import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "os"
)

type TeamPassFile struct {
    Comment string `yaml:"comment"`
    PublicKeys []TeamPassKey `yaml:"public_keys"`
}

type TeamPassKey struct {
    Alias string `yaml:"alias"`
    Value string `yaml:"value"`
    Comment string `yaml:"comment"`
}

func readFile(filename *string) (teamPassFile TeamPassFile, err error) {
    _, err = os.Stat(*filename)
    if os.IsNotExist(err) {
        err = nil
        return
    }

    file, err := os.Open(*filename)
    if err != nil {
        return
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

func writeFile(filename *string, teamPassFile TeamPassFile) (err error) {
    file, err := os.Create(*filename)
    if err != nil {
        return
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
