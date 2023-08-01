package onion

import (
	"bytes"
	"encoding/base32"
	"go-dmtor/config"
	"os"
	"path/filepath"
	"strings"
	"time"

	// "crypto/rand"

	"math/rand"

	"github.com/cretz/bine/torutil/ed25519"

	"golang.org/x/crypto/sha3"
)

var SESSION_DIR = config.SESSION_DIR

// Functions
// Read onion struct from session file
// Create onion struct from priv key bytes

type Onion struct {
	keyPair     *ed25519.PrivateKey
	pubKey      *ed25519.PublicKey
	pubKeyBytes []byte
	address     string
}

// new tor session
func New() (*Onion, error) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	keyPair, err := ed25519.GenerateKey(rnd)
	if err != nil {
		return nil, err
	}
	return NewFromPrivKey(keyPair.PrivateKey())
}

// new tor session from session file (priv key)
func NewFromSession(filename string) (*Onion, error) {
	path := filepath.Join(SESSION_DIR, filename)
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	privKey := ed25519.PrivateKey(bytes)
	return NewFromPrivKey(privKey)
}

func NewFromPrivKey(privKeyBytes []byte) (*Onion, error) {
	privKey := ed25519.PrivateKey(privKeyBytes)
	pubKey := privKey.Public().(ed25519.PublicKey)
	pubKeyBytes := []byte(pubKey)
	address, err := pubKeyToAddress(pubKeyBytes)
	if err != nil {
		return nil, err
	}
	return &Onion{&privKey, &pubKey, pubKeyBytes, address}, nil
}

func NewFromPubKey(pubKeyBytes []byte) (*Onion, error) {
	pubKey := ed25519.PublicKey(pubKeyBytes)
	address, err := pubKeyToAddress(pubKeyBytes)
	if err != nil {
		return nil, err
	}
	return &Onion{nil, &pubKey, pubKeyBytes, address}, nil
}

func (o *Onion) PubKey() []byte {
	return o.pubKeyBytes
}

func (o *Onion) PrivKey() []byte {
	return o.keyPair.PrivateKey()
}

func (o *Onion) Address() string {
	return o.address
}

// TODO: save as session ID not onion address
// func (o *Onion) Save() error {
// 	// create session dir if not exists
// 	err := os.MkdirAll(SESSION_DIR, 0700)
// 	if err != nil {
// 		return err
// 	}
// 	path := filepath.Join(SESSION_DIR, o.session)
// 	f, err := os.Create(path)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()
// 	_, err = f.Write(o.keyPair.PrivateKey())
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// TODO: move to utils
// func sessionName(address string) string {
// 	hash := sha3.Sum256([]byte(address))
// 	hex := strings.ToUpper(fmt.Sprintf("%x", hash))
// 	// split into 4 parts
// 	parts := []string{}
// 	p := 4
// 	for i := 0; i < p; i++ {
// 		parts = append(parts, hex[i*p:i*p+p])
// 		// keep only 2 parts
// 		if len(parts) == 2 {
// 			break
// 		}
// 	}
// 	return strings.Join(parts, "-")
// }

func pubKeyToAddress(pubKeyBytes []byte) (string, error) {
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
