package cryptotools

import (
	"bytes"
	// "crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base32"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"hash/fnv"
	"log"
	"strings"

	"github.com/cretz/bine/torutil/ed25519"
	"golang.org/x/crypto/sha3"
)

// ONION ADDR encryption

// encrypt onion pub key to hex format
// TODO: add AES encryption with password
func KeyFromOnionPubKey(pubKey []byte) string {
	// encode to HEX
	hex := fmt.Sprintf("%x", pubKey)
	hex = strings.ToUpper(hex)
	// split to 4 byte parts
	parts := make([]string, 0)
	for i := 0; i < len(hex); i += 4 {
		parts = append(parts, hex[i:i+4])
	}
	return strings.Join(parts, "-")
}
func KeyToOnionPubKey(hexkey string) ([]byte, error) {
	bHex := strings.ReplaceAll(hexkey, "-", "")
	bHex = strings.ToLower(bHex)
	return hex.DecodeString(bHex)
}

func KeyToOnionAddress(hexkey string) (string, error) {
	pubKeyBytes, err := KeyToOnionPubKey(hexkey)
	if err != nil {
		return "", err
	}

	publicKey := ed25519.PublicKey(pubKeyBytes)

	// convert from pub key to onion address

	// checksum = H(".onion checksum" || pubkey || version)
	var checksumBytes bytes.Buffer
	checksumBytes.Write([]byte(".onion checksum"))
	checksumBytes.Write([]byte(publicKey))
	checksumBytes.Write([]byte{0x03})
	checksum := sha3.Sum256(checksumBytes.Bytes())

	// onion_address = base32(pubkey || checksum || version)
	var onionAddressBytes bytes.Buffer
	onionAddressBytes.Write([]byte(publicKey))
	onionAddressBytes.Write([]byte(checksum[:2]))
	onionAddressBytes.Write([]byte{0x03})
	onionAddress := base32.StdEncoding.EncodeToString(onionAddressBytes.Bytes())

	return strings.ToLower(onionAddress), nil
}

// MSG encryption

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

// pub key conversion

func PubToBytes(pub *rsa.PublicKey) ([]byte, error) {
	return x509.MarshalPKIXPublicKey(pub)
}

func BytesToPub(pubBytes []byte) (*rsa.PublicKey, error) {
	publicKeyInterface, err := x509.ParsePKIXPublicKey(pubBytes)
	if err != nil {
		return nil, err
	}
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	return publicKey, nil
}

// PEM

func PubToPem(pub *rsa.PublicKey) ([]byte, error) {
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
func PemToPub(pubBytes []byte) (*rsa.PublicKey, error) {
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

func HashToUInt(msg []byte) uint64 {
	hash := fnv.New64a()
	hash.Write(msg)
	return hash.Sum64()
}
