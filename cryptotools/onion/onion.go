package onion

import (
	"bytes"
	"encoding/base32"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cretz/bine/torutil/ed25519"

	"golang.org/x/crypto/sha3"
)

const SESSION_DIR = "sessions"

func PrivKeySave(privKey []byte) error {
	// create session dir if not exists
	err := os.MkdirAll(SESSION_DIR, 0700)
	if err != nil {
		return err
	}
	onion, err := PrivKeyBytesToOnionAddress(privKey)
	if err != nil {
		return err
	}
	path := filepath.Join(SESSION_DIR, fmt.Sprintf("%s.onion", onion))
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(privKey)
	if err != nil {
		return err
	}
	return nil
}

func PrivKeyRead(onion string) ([]byte, error) {
	path := filepath.Join(SESSION_DIR, fmt.Sprintf("%s.onion", onion))
	return os.ReadFile(path)
}

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
	pubKeyBytes := keyPair.Public().(ed25519.PublicKey)
	fmt.Printf("privKey hex: %x\n", privKeyBytes)
	fmt.Printf("pubKey hex: %x\n", pubKeyBytes)
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
