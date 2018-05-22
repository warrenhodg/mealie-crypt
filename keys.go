package main

import (
    sshEncoding "github.com/ianmcmahon/encoding_ssh"
    "io/ioutil"
    "os"
)

func readPublicKey(filename *string) (keyContent string, key interface{}, err error) {
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

    key, err = sshEncoding.DecodePublicKey(keyContent)
    if err != nil {
        return
    }

    return
}
