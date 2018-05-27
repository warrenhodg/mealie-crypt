package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
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

func createSymmetricalKey(bits int) (result string, err error) {
	cmd := exec.Command("openssl", "rand", fmt.Sprintf("%d", bits))

	bytes, err := cmd.CombinedOutput()
	if err != nil {
		err = errors.New(fmt.Sprintf("%s : %s", err.Error(), string(bytes)))
		return
	}

	return string(bytes), nil
}

func encryptSymmetricalKey(symKey string, publicKey string) (encSymKey string, err error) {
	cmd := exec.Command("bash", "-c", "openssl rsautl -encrypt -oaep -pubin -inkey <(cat <<< \"$PUB_KEY\")")
	cmd.Stdin = strings.NewReader(symKey)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PUB_KEY=%s", publicKey),
	)

	bytes, err := cmd.CombinedOutput()
	if err != nil {
		err = errors.New(fmt.Sprintf("%s : %s", err.Error(), string(bytes)))
		return
	}

	return string(bytes), nil
}

func decryptSymmetricalKey(encSymKey string, privateKeyFile string) (symKey string, err error) {
	cmd := exec.Command("openssl", "rsautl", "-decrypt", "-oaep", "-inkey", privateKeyFile)
	cmd.Stdin = strings.NewReader(encSymKey)

	bytes, err := cmd.CombinedOutput()
	if err != nil {
		err = errors.New(fmt.Sprintf("%s : %s", err.Error(), string(bytes)))
		return
	}

	return string(bytes), nil
}
