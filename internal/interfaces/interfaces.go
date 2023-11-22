package interfaces

// auth
type Symmetric interface {
	Encrypt([]byte, string) ([]byte, error)
	Decrypt([]byte, string) ([]byte, error)
}

// msg crypter
type Asymmetric interface {
	Encrypt([]byte, []byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
	PubKey() []byte
}

// tor onion address
type Onioner interface {
	PubKey() []byte
	PrivKey() []byte
	Address() string
}
