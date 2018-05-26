package main

import (
	"errors"
	"fmt"
	"os/exec"
)

func readPublicKey(filename *string) (keyContent string, err error) {
	cmd := exec.Command("ssh-keygen", "-e", "-f", *filename, "-m", "PKCS8")
	pubkeyBytes, err := cmd.CombinedOutput()
	if err != nil {
		err = errors.New(fmt.Sprintf("%s : %s", err.Error(), string(pubkeyBytes)))
		return
	}

	keyContent = string(pubkeyBytes)
	return
}

func CreateKey(len int) (result string, err error) {
	cmd := exec.Command("openssl", "rand", fmt.Sprintf("%d", len))
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		err = errors.New(fmt.Sprintf("%s : %s", err.Error(), string(bytes)))
		return
	}

	return string(bytes), nil
}
