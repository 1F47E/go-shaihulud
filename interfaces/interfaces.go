package interfaces

type Symmetric interface {
	Encrypt([]byte, string) ([]byte, error)
	Decrypt([]byte, string) ([]byte, error)
}

type Asymmetric interface {
	Encrypt([]byte, []byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
	PubKey() []byte
}

type Onioner interface {
	PubKey() []byte
	PrivateKey() []byte
	Address() string
	Save() error
}
