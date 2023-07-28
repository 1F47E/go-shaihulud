package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"log"
)

func Encrypt(msg string, publicKey *rsa.PublicKey) []byte {
	label := []byte("")
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, publicKey, []byte(msg), label)
	if err != nil {
		log.Fatal(err)
	}
	return ciphertext
}

func Decrypt(ciphertext []byte, privateKey *rsa.PrivateKey) string {
	label := []byte("")
	hash := sha256.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, privateKey, ciphertext, label)
	if err != nil {
		log.Fatal(err)
	}
	return string(plaintext)
}

func Keygen() rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
	return *privateKey
}
