// The AUTH package works with an access key and password.
//
// The access key is a readable binary key in hex format
// that resembles 1234-ABCD-EFGH-5678....
// This key represents the AES-encrypted onion address.
//
// The password is used to encrypt/decrypt the access key to obtain the onion address,
// which takes a form like 1234-ABCD.
// The password consists of random bytes converted to upper-case hex format.
// It is also used to sign messages via HMAC to verify message integrity.
//
// Workflow:
// User A (server), after connecting to Tor and generating an onion address,
// encrypts this address with a randomly generated password.
// User A then shares the access key (AES-encrypted onion address) and password with User B.
// The password and access key should be shared via different channels for security.
// User B enters the access key and then the password to decrypt the onion address.

package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"go-dmtor/interfaces"
	"log"
	"strings"
)

type Auth struct {
	crypter interfaces.Symmetric
}

func New(crypter interfaces.Symmetric) *Auth {
	return &Auth{crypter}
}

func (ac *Auth) Encrypt(plaintext []byte, password string) ([]byte, error) {
	return ac.crypter.Encrypt(plaintext, password)
}
func (ac *Auth) Decrypt(ciphertext []byte, password string) ([]byte, error) {
	return ac.crypter.Decrypt(ciphertext, password)
}

// ETC
func generatePassword() string {
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

// ONION STUFF
// encrypt onion pub key to hex format
func AuthKeyFromOnionPubKey(pubKey []byte) string {
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

func AuthKeyToOnionPubKey(hexkey string) ([]byte, error) {
	bHex := strings.ReplaceAll(hexkey, "-", "")
	bHex = strings.ToLower(bHex)
	return hex.DecodeString(bHex)
}
