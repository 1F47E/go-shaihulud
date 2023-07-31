package onion

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestOnionEncryptDecrypt(t *testing.T) {
	// demo data
	keyHex := "f86af341ed3a612ff0754c77b33b60eb1cd40ed204603134217bc857fc411867f2ff4a987ffefcdb3a8c1b5af8f22aff8bcdeb4420426b491747f73fc5327d24"
	expected_onion := "gbislcwjbx2h3pkdavsaqku3mlx4tcnfhmpq2gji5nyegrbrsqcvv6qd"

	keyBytes, err := hex.DecodeString(keyHex)
	if err != nil {
		t.Fatalf("key decode error: %v\n", err)
	}
	fmt.Printf("keyBytes: %v\n", keyBytes)
	onion, err := PrivKeyBytesToOnionAddress(keyBytes)
	if err != nil {
		t.Fatalf("onion decode error: %v\n", err)
	}
	if onion != expected_onion {
		t.Fatalf("onion mismatch: %s != %s\n", onion, expected_onion)
	}
}
