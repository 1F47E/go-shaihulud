package myrsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
)

// MSG RSA encryption
// implementation of interfaces.Asymmetric interface

type RsaCrypter struct {
	privKey *rsa.PrivateKey
	pubKey  *rsa.PublicKey
}

func New() (*RsaCrypter, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	// log.Printf("rsa private key: %x\n", privateKey)
	if err != nil {
		return nil, err
	}
	pubKey := &privKey.PublicKey
	return &RsaCrypter{privKey: privKey, pubKey: pubKey}, nil
}

// encrypt our message with user B public key
func (r *RsaCrypter) Encrypt(data []byte, pubKeyBytes []byte) ([]byte, error) {
	pubKey, err := x509.ParsePKCS1PublicKey(pubKeyBytes)
	if err != nil {
		return nil, err
	}
	label := []byte("")
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pubKey, data, label)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

func (r *RsaCrypter) Decrypt(data []byte) ([]byte, error) {
	label := []byte("")
	hash := sha256.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, r.privKey, data, label)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// get pub key as bytes to send over network
func (r *RsaCrypter) PubKey() []byte {
	return x509.MarshalPKCS1PublicKey(r.pubKey)
}
