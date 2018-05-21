package main

import (
    "io/ioutil"
    "os"
)

func readKey(filename *string) (keyContent string, err error) {
    var file *os.File

    file, err = os.Open(*filename)
    if err != nil {
        return
    }

    defer file.Close()

    bytes, err := ioutil.ReadAll(file)
    if err != nil {
        return
    }

    keyContent = string(bytes)
    return
}
