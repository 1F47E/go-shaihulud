package cryptotools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/cretz/bine/torutil/ed25519"
	"golang.org/x/crypto/sha3"
)

// ONION encryption

// encrypt onion pub key to hex format
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

// AES

func AESDeriveKey(password string) []byte {
	// hash the password, because we need 256-bit key

	return nil
}

func AESEncrypt(plaintext []byte, password string) ([]byte, error) {
	// hash the password, because we need 256-bit key for encoding
	// TODO: change slow key derivation function
	// scrypt KDF
	// docs https://pkg.go.dev/golang.org/x/crypto/scrypt
	key := sha256.Sum256([]byte(password))

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func AESDecrypt(cipherBytes []byte, password string) ([]byte, error) {
	// Hash the password
	// TODO: change to slow key derivation function
	key := sha256.Sum256([]byte(password))

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// get and check the nonce from the ciphertext
	nonceSize := gcm.NonceSize()
	if len(cipherBytes) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, cipherBytes := cipherBytes[:nonceSize], cipherBytes[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, cipherBytes, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// MSG RSA encryption

func MessageEncrypt(msg []byte, publicKey *rsa.PublicKey) []byte {
	label := []byte("")
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, publicKey, msg, label)
	if err != nil {
		log.Fatal(err)
	}
	return ciphertext
}

func MessageDecrypt(ciphertext []byte, privateKey *rsa.PrivateKey) string {
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

func AccessPinGenerate() string {
	// format is AB3D-E2FA
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		// crypto/rand error. Highly unlikely.
		log.Fatalf("error reading random bytes: %v", err)
	}
	pin := fmt.Sprintf("%x", b)
	pin = strings.ToUpper(pin)
	// split to 4 byte parts
	parts := make([]string, 0)
	for i := 0; i < len(pin); i += 4 {
		parts = append(parts, pin[i:i+4])
	}
	pin = strings.Join(parts, "-")
	return pin
}
