package msgcrypter

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockAsymmetric struct {
	mock.Mock
}

func (m *mockAsymmetric) Encrypt(plaintext []byte, pubKey []byte) ([]byte, error) {
	args := m.Called(plaintext, pubKey)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockAsymmetric) Decrypt(ciphertext []byte) ([]byte, error) {
	args := m.Called(ciphertext)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockAsymmetric) PubKey() []byte {
	args := m.Called()
	return args.Get(0).([]byte)
}

func TestMessageCrypter(t *testing.T) {
	m := new(mockAsymmetric)
	mc := New(m)

	t.Run("Test Encrypt", func(t *testing.T) {
		m.On("Encrypt", []byte("plaintext"), []byte("pubKey")).Return([]byte("encryptedtext"), nil)
		ciphertext, err := mc.Encrypt([]byte("plaintext"), []byte("pubKey"))
		require.NoError(t, err)
		require.Equal(t, []byte("encryptedtext"), ciphertext)
	})

	t.Run("Test Decrypt", func(t *testing.T) {
		m.On("Decrypt", []byte("encryptedtext")).Return([]byte("decryptedtext"), nil)
		plaintext, err := mc.Decrypt([]byte("encryptedtext"))
		require.NoError(t, err)
		require.Equal(t, []byte("decryptedtext"), plaintext)
	})

	t.Run("Test PubKey", func(t *testing.T) {
		m.On("PubKey").Return([]byte("pubKey"))
		pubKey := mc.PubKey()
		require.Equal(t, []byte("pubKey"), pubKey)
	})
}
