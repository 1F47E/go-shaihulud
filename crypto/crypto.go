package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
)

func Encrypt(msg []byte, publicKey *rsa.PublicKey) []byte {
	label := []byte("")
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, publicKey, msg, label)
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

// Function to encode an rsa.PublicKey to bytes
func EncodePublicKeyToBytes(pub *rsa.PublicKey) ([]byte, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes, nil
}

// Function to decode an rsa.PublicKey from bytes
func DecodePublicKeyFromBytes(pubBytes []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubBytes)
	if block == nil {
		return nil, errors.New("failed to parse PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("not RSA public key")
	}
}
