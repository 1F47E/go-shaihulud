package onion

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base32"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/sha3"
)

func PrivKeyFileToOnionAddress(privKeyFile string) (string, error) {
	keybytes, err := os.ReadFile(privKeyFile)
	if err != nil {
		return "", err
	}
	return PrivKeyBytesToOnionAddress(keybytes)
}

func PrivKeyBytesToOnionAddress(privKeyBytes []byte) (string, error) {
	if len(privKeyBytes) != 64 {
		return "", fmt.Errorf("invalid priv key bytes length: %d", len(privKeyBytes))
	}
	keyPair := ed25519.PrivateKey(privKeyBytes)
	// get public key bytes
	pubKeyBytes := keyPair.Public().(ed25519.PublicKey)
	// get onion from public key
	addr, err := PubKeyToOnionAddress(pubKeyBytes)
	if err != nil {
		return "", err
	}
	return addr, nil
}

func PubKeyToOnionAddress(pubKeyBytes []byte) (string, error) {
	pubKey := ed25519.PublicKey(pubKeyBytes)

	// checksum = H(".onion checksum" || pubkey || version)
	var checksumBytes bytes.Buffer
	checksumBytes.Write([]byte(".onion checksum"))
	checksumBytes.Write([]byte(pubKey))
	checksumBytes.Write([]byte{0x03})
	checksum := sha3.Sum256(checksumBytes.Bytes())

	// onion_address = base32(pubkey || checksum || version)
	var onionAddressBytes bytes.Buffer
	onionAddressBytes.Write([]byte(pubKey))
	onionAddressBytes.Write([]byte(checksum[:2]))
	onionAddressBytes.Write([]byte{0x03})
	onionAddress := base32.StdEncoding.EncodeToString(onionAddressBytes.Bytes())

	return strings.ToLower(onionAddress), nil
}
