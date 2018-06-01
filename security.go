package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	encoding_ssh "github.com/ianmcmahon/encoding_ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"io/ioutil"
	"syscall"
)

var keyLenBits = 256

func readPublicKey(filename *string) (keyContent string, err error) {
	pubKeyBytes, err := ioutil.ReadFile(*filename)
	if err != nil {
		return
	}

	pubKeyI, err := encoding_ssh.DecodePublicKey(string(pubKeyBytes))
	if err != nil {
		return
	}
	pubKey, convertOk := pubKeyI.(*rsa.PublicKey)
	if !convertOk {
		err = errors.New("Public key is not an RSA public key")
		return
	}

	derBytes := x509.MarshalPKCS1PublicKey(pubKey)

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
	block, remains := pem.Decode([]byte(publicKey))
	if len(remains) > 0 {
		err = errors.New(fmt.Sprintf("Public key contains extra characters at end"))
		return
	}

	pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return
	}

	encBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, []byte(symKey), nil)
	if err != nil {
		return
	}

	return string(encBytes), nil
}

func getPassword(prompt string) string {
	fmt.Printf("%s", prompt)
	password, _ := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Printf("\n")
	return string(password)
}

func readPrivateKey(privateKeyFile string) (pvtKey *rsa.PrivateKey, err error) {
	privateKeyBytes, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return
	}

	block, remains := pem.Decode([]byte(privateKeyBytes))
	if len(remains) > 0 {
		err = errors.New(fmt.Sprintf("Public key contains extra characters at end"))
		return
	}

	//Decrypt if necessary
	if x509.IsEncryptedPEMBlock(block) {
		var derBytes []byte
		password := getPassword("Enter the password for the private key file : ")
		derBytes, err = x509.DecryptPEMBlock(block, []byte(password))
		if err != nil {
			return
		}

		block = &pem.Block{
			Type:  block.Type,
			Bytes: derBytes,
		}
	}

	pvtKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	return
}

func decryptSymmetricalKey(encSymKey string, pvtKey *rsa.PrivateKey) (symKey string, err error) {
	decBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, pvtKey, []byte(encSymKey), nil)
	if err != nil {
		return
	}

	return string(decBytes), nil
}

func encryptValue(symKey64 string, value string) (encValue string, err error) {
	symKeyBytes, err := base64.StdEncoding.DecodeString(symKey64)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(symKeyBytes)
	if err != nil {
		return
	}

	cipherBytes := make([]byte, aes.BlockSize+len(value))
	iv := cipherBytes[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherBytes[aes.BlockSize:], []byte(value))

	encValue = string(cipherBytes)
	return
}

func decryptValue(symKey64 string, encValue string) (decValue string, err error) {
	symKeyBytes, err := base64.StdEncoding.DecodeString(symKey64)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(symKeyBytes)
	if err != nil {
		return
	}

	encBytes := []byte(encValue)

	if len(encBytes) < aes.BlockSize {
		err = errors.New(fmt.Sprintf("Encrypted value block size is too short"))
		return
	}

	iv := encBytes[:aes.BlockSize]
	encBytes = encBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encBytes, encBytes)

	decValue = string(encBytes)
	return
}
