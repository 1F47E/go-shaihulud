package asymmetric

// msg crypter
type Asymmetric interface {
	Encrypt([]byte, []byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
	PubKey() []byte
}
