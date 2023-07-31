package onion

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/sha3"
)

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
