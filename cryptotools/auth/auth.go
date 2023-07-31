package auth

import (
	"crypto/rand"
	"fmt"
	"go-dmtor/interfaces"
	"log"
	"strings"
)

type AccessKeyCrypter struct {
	crypter interfaces.Symmetric
}

func New(crypter interfaces.Symmetric) *AccessKeyCrypter {
	return &AccessKeyCrypter{crypter}
}

func (ac *AccessKeyCrypter) Encrypt(plaintext []byte, password string) ([]byte, error) {
	return ac.crypter.Encrypt(plaintext, password)
}
func (ac *AccessKeyCrypter) Decrypt(ciphertext []byte, password string) ([]byte, error) {
	return ac.crypter.Decrypt(ciphertext, password)
}

func newPin() string {
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
