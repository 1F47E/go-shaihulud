package msgcrypter

import "github.com/1F47E/go-shaihulud/internal/interfaces"

type MessageCrypter struct {
	crypter interfaces.Asymmetric
}

func New(crypter interfaces.Asymmetric) *MessageCrypter {
	return &MessageCrypter{crypter}
}

func (mc *MessageCrypter) Encrypt(plaintext []byte, pubKey []byte) ([]byte, error) {
	return mc.crypter.Encrypt(plaintext, pubKey)
}

func (mc *MessageCrypter) Decrypt(ciphertext []byte) ([]byte, error) {
	return mc.crypter.Decrypt(ciphertext)
}

func (mc *MessageCrypter) PubKey() []byte {
	return mc.crypter.PubKey()
}
