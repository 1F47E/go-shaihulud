package onion

// tor onion address
type Onioner interface {
	PubKey() []byte
	PrivKey() []byte
	Address() string
}
