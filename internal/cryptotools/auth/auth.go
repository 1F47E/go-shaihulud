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
	mrand "math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/1F47E/go-shaihulud/internal/config"
	"github.com/1F47E/go-shaihulud/internal/cryptotools/onion"
	"github.com/1F47E/go-shaihulud/internal/cryptotools/symmetric"
)

var SESSION_DIR = config.SESSION_DIR

// NOTE:
// for the access key we encode onion pub key (32 bytes) to hex format
// for our session file we encode onion priv key without encryption

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
	crypter   symmetric.Symmetric
	password  string
	accessKey string
	onioner   onion.Onioner
}

func New(crypter symmetric.Symmetric, session string) (*Auth, error) {
	var err error
	password := generatePassword()
	accessKey := ""
	var onioner onion.Onioner

	// create or load onion key
	if session == "" {
		// create new onion
		o, err := onion.New()
		if err != nil {
			return nil, err
		}
		onioner = o
	} else {
		// load onion from the session file
		o, err := onion.NewFromSession(session)
		if err != nil {
			return nil, err
		}
		onioner = o
	}

	// encrypt onion pub key for the user B
	onionPubKeyCipher, err := crypter.Encrypt(onioner.PubKey(), password)
	if err != nil {
		return nil, err
	}

	accessKey = Encode(onionPubKeyCipher)

	a := Auth{
		crypter:   crypter,
		onioner:   onioner,
		password:  password,
		accessKey: accessKey,
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func NewFromKey(crypter symmetric.Symmetric, accessKey, password string) (*Auth, error) {
	// decode string key to bytes
	keyBytesCipher, err := Decode(accessKey)
	if err != nil {
		return nil, err
	}

	// decrypt key bytes with password
	keyBytes, err := crypter.Decrypt(keyBytesCipher, password)
	if err != nil {
		return nil, err
	}

	// version of onion without priv key, only pub key to connect to
	onion, err := onion.NewFromPubKey(keyBytes)
	if err != nil {
		return nil, err
	}
	a := Auth{
		crypter:   crypter,
		accessKey: accessKey,
		password:  password,
		onioner:   onion,
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (a *Auth) Encrypt(data []byte) ([]byte, error) {
	return a.crypter.Encrypt(data, a.password)
}

func (a *Auth) Decrypt(ciphertext []byte) ([]byte, error) {
	return a.crypter.Decrypt(ciphertext, a.password)
}

// for testing
func (a *Auth) EncryptWithPassword(data []byte, password string) ([]byte, error) {
	return a.crypter.Encrypt(data, password)
}

// for testing
func (a *Auth) DecryptWithPassword(ciphertext []byte, password string) ([]byte, error) {
	return a.crypter.Decrypt(ciphertext, password)
}

func (a *Auth) OnionAddress() string {
	return a.onioner.Address()
}

func (a *Auth) OnionAddressFull() string {
	return fmt.Sprintf("%s.onion:80", a.onioner.Address())
}

func (a *Auth) Onion() onion.Onioner {
	return a.onioner
}

func (a *Auth) AccessKey() string {
	return a.accessKey
}

func (a *Auth) Password() string {
	return a.password
}

func (a *Auth) String() string {
	return fmt.Sprintf("=====\nAccess key:\n%s\nPassword: %s\n=====", a.accessKey, a.password)
}

// save to a session file
func (a *Auth) Save() error {

	data := a.onioner.PrivKey()
	if len(data) == 0 {
		return fmt.Errorf("no data to save")
	}

	// create session dir if not exists
	err := os.MkdirAll(SESSION_DIR, 0700)
	if err != nil {
		return err
	}
	path := filepath.Join(SESSION_DIR, a.accessKey)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// ETC
// len is 9 bytes
// TODO: make a test
// format is AB3D-E2FA
func generatePassword() string {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		// crypto/rand error. Highly unlikely.
		// fall back to math/rand that always works
		r := mrand.Uint32()
		b = []byte(fmt.Sprintf("%x", r))
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
	hex := fmt.Sprintf("%x", pubKey)
	hex = strings.ToUpper(hex)
	return hex
}

// decode from custom hex format to bytes
func Decode(key string) ([]byte, error) {
	key = strings.ToLower(key)
	return hex.DecodeString(key)
}
