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

func createSymmetricalKey(lenBits int) (result string, err error) {
	lenBytes := lenBits / 8
	cmd := exec.Command("bash", "-c", fmt.Sprintf("openssl rand %d | base64", lenBytes))

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

func encryptValue(symKey string, value string) (encValue string, err error) {
	cmd := exec.Command("bash", "-c", "openssl aes-256-cbc -e -pass file:<(cat <<< $SYM_KEY | base64 -D)")
	cmd.Stdin = strings.NewReader(value)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("SYM_KEY=%s", symKey),
	)

	bytes, err := cmd.CombinedOutput()
	if err != nil {
		err = errors.New(fmt.Sprintf("%s : %s", err.Error(), string(bytes)))
		return
	}

	return string(bytes), nil
}

func decryptValue(symKey string, value string) (decValue string, err error) {
	cmd := exec.Command("bash", "-c", "openssl aes-256-cbc -d -pass file:<(cat <<< $SYM_KEY | base64 -D)")
	cmd.Stdin = strings.NewReader(value)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("SYM_KEY=%s", symKey),
	)

	bytes, err := cmd.CombinedOutput()
	if err != nil {
		err = errors.New(fmt.Sprintf("%s : %s", err.Error(), string(bytes)))
		return
	}

	return string(bytes), nil
}
