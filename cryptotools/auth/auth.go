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
	"go-dmtor/cryptotools/onion"
	"go-dmtor/interfaces"
	"log"
	"strings"
)

// TODO:
// on server startup
// - generate password
// - after onion is created, encrypt it with a random password
// - format encrypted binary as access key AB3D-E2FA-...
// - display access key to user
// on connection
// - user enters access key
// - ask for password
// - save passwrod for signatures
// - decrypt access key, get onion address
// - connect to onion address
// - add hmac signature to messages

type Auth struct {
	crypter   interfaces.Symmetric
	password  string
	accessKey string
	onion     interfaces.Onioner
}

func New(crypter interfaces.Symmetric, onion interfaces.Onioner) *Auth {
	password := generatePassword()
	accessKey := Encode(onion.PubKey())

	return &Auth{
		crypter:   crypter,
		onion:     onion,
		password:  password,
		accessKey: accessKey,
	}
}

func NewFromKey(crypter interfaces.Symmetric, accessKey, password string) (*Auth, error) {
	// TODO: test
	keyBytes, err := Decode(accessKey)
	if err != nil {
		return nil, err
	}
	onion, err := onion.NewFromPrivKey(keyBytes)
	if err != nil {
		return nil, err
	}
	return &Auth{
		crypter:   crypter,
		accessKey: accessKey,
		password:  password,
		onion:     onion,
	}, nil
}

func (ac *Auth) Encrypt(plaintext []byte) ([]byte, error) {
	return ac.crypter.Encrypt(plaintext, ac.password)
}
func (ac *Auth) Decrypt(ciphertext []byte) ([]byte, error) {
	return ac.crypter.Decrypt(ciphertext, ac.password)
}

func (ac *Auth) String() string {
	return fmt.Sprintf("=====\nAccess key: %s\nPassword: %s\n=====", ac.accessKey, ac.password)
}

func (ac *Auth) OnionAddress() string {
	return ac.onion.Address()
}

// ETC
// len is 9 bytes
// TODO: make a test
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

// TODO: rewrite to be a method of a sctuct
// encrypt onion pub key to hex format
// example of access key format: AB3D-E2FA-...
func Encode(pubKey []byte) string {
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

// decode from custom hex format to bytes
func Decode(key string) ([]byte, error) {
	bHex := strings.ReplaceAll(key, "-", "")
	bHex = strings.ToLower(bHex)
	return hex.DecodeString(bHex)
}
