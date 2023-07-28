package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"log"
)

func encrypt(msg string, publicKey *rsa.PublicKey) []byte {
	label := []byte("")
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, publicKey, []byte(msg), label)
	if err != nil {
		log.Fatal(err)
	}
	return ciphertext
}

func decrypt(ciphertext []byte, privateKey *rsa.PrivateKey) string {
	label := []byte("")
	hash := sha256.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, privateKey, ciphertext, label)
	if err != nil {
		log.Fatal(err)
	}
	return string(plaintext)
}

func keygen() rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
	return *privateKey
}

func main() {

	key := keygen()
	// Get the public key
	publicKey := &key.PublicKey

	// Encrypt a message
	message := "hello, world"
	cipher := encrypt(message, publicKey)
	fmt.Printf("Ciphertext: %x\n", cipher)

	plain := decrypt(cipher, &key)
	fmt.Printf("Plaintext: %s\n", plain)
}
