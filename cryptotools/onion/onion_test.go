package onion

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestOnionEncryptDecrypt(t *testing.T) {
	// demo data
	keyHex := "207996b0f5b70b25e0a94fe1cbb365a4a26e6ade6d5e66818ea25467e2005550be4a38a3eefb4c53427a22a4c208937e6c0cbe47cb62d2a477cb0820d40b6ecf"
	expected_onion := "xzfdri7o7ngfgqt2eksmecetpzwazpshznrnfjdxzmecbvaln3h2dbqd"

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
