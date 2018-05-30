package main

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	encoding_ssh "github.com/ianmcmahon/encoding_ssh"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var keyLenBits = 256

func readPublicKey(filename *string) (keyContent string, err error) {
	pubKeyBytes, err := ioutil.ReadFile(*filename)
	if err != nil {
		return
	}

	pubKey, err := encoding_ssh.DecodePublicKey(string(pubKeyBytes))
	if err != nil {
		return
	}

	derBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return
	}

	pubKeyPem := string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: derBytes,
	}))

	keyContent = string(pubKeyPem)
	return
}

func createSymmetricalKey() (result string, err error) {
	lenBytes := keyLenBits / 8
	bytes := make([]byte, lenBytes)
	_, err = rand.Read(bytes)
	if err != nil {
		return
	}

	result = base64.StdEncoding.EncodeToString(bytes)
	return
}

func encryptSymmetricalKey(symKey string, publicKey string) (encSymKey string, err error) {
	cmd := exec.Command("bash", "-c", "openssl rsautl -encrypt -oaep -pubin -inkey <(cat <<< \"$PUB_KEY\")")
	cmd.Stdin = strings.NewReader(symKey)
	cmd.Env = append(
		os.Environ(),
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
	cmd.Env = os.Environ()
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
	cmd.Env = append(
		os.Environ(),
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
	cmd.Env = append(
		os.Environ(),
		fmt.Sprintf("SYM_KEY=%s", symKey),
	)

	bytes, err := cmd.CombinedOutput()
	if err != nil {
		err = errors.New(fmt.Sprintf("%s : %s", err.Error(), string(bytes)))
		return
	}

	return string(bytes), nil
}
