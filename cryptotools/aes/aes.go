package caes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/scrypt"
)

type AEScrypter struct{}

func New() *AEScrypter {
	return new(AEScrypter)
}

// using scrypt for brute force resistance
func (a *AEScrypter) Encrypt(data []byte, password string) ([]byte, error) {

	// transform text password into appropriate 32 byte key for AES
	salt, err := aesSaltGen()
	if err != nil {
		return nil, err
	}
	key, err := aesDeriveKey([]byte(password), salt)
	if err != nil {
		return nil, err
	}

	// generate a new aes cipher using our 32 byte long key
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	// if any error generating new GCM
	// handle them
	if err != nil {
		return nil, err
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// encrypt our data using the gsm Seal
	cipher := gcm.Seal(nonce, nonce, data, nil)
	// add salt at the end to get in back on decoding
	cipher = append(cipher, salt...)
	return cipher, nil
}

func (a *AEScrypter) Decrypt(data []byte, password string) ([]byte, error) {

	// data has salt already
	// check input text length
	if len(data) < 32 {
		return nil, errors.New("invalid data len")
	}
	// fmt.Printf("data len: %d\n", len(data))

	// get salt from the end
	salt, ciphertext := data[len(data)-32:], data[:len(data)-32]

	key, err := aesDeriveKey([]byte(password), salt)
	if err != nil {
		return nil, err
	}

	// generate a new aes cipher using our 32 byte long key
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	// get the nonce size
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		fmt.Println(err)
	}
	// extract our nonce from our encrypted text
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plain, nil
}

func aesSaltGen() ([]byte, error) {
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

func aesDeriveKey(password, salt []byte) ([]byte, error) {
	// minimum N is 16384
	// x32 will take about 2 sec
	n := 16384 * 32
	key, err := scrypt.Key(password, salt, n, 8, 1, 32)
	if err != nil {
		return nil, err
	}

	return key, nil
}
