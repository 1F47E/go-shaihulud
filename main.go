package main

import (
	"go-dmtor/client"
	cfg "go-dmtor/config"
	"go-dmtor/logger"
	"os"
)

var log = logger.New()

var usage = "Usage: %s <srv|cli>\n"

func main() {
	args := os.Args
	if len(args) == 1 {
		log.Fatalf(usage, args[0])
	}
	arg := args[1]
	if arg != "srv" && arg != "cli" {
		log.Fatalf(usage, args[0])
	}
	cli := client.NewClient(cfg.ADDR)
	if arg == "srv" {
		cli.ServerStart()
		// crypt_demo()
	} else {
		// read user input
		err := cli.ServerConnect()
		if err != nil {
			log.Fatalf("connect error: %v\n", err)
		}
	}

	// block and wait for user input
	for {
		input := make([]byte, cfg.MSG_MAX_SIZE)
		_, err := os.Stdin.Read(input)
		if err != nil {
			log.Fatalf("read error: %v\n", err)
			return
		}
		cli.MsgCh <- string(input)
	}
}

// ====== CLIENT

// func crypt_demo() {

// 	key := crypto.Keygen()
// 	// Get the public key
// 	publicKey := &key.PublicKey

// 	// Encrypt a message
// 	message := "hello, world"
// 	cipher := crypto.Encrypt(message, publicKey)
// 	fmt.Printf("Ciphertext: %x\n", cipher)

// 	plain := crypto.Decrypt(cipher, &key)
// 	fmt.Printf("Plaintext: %s\n", plain)
// }
