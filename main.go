package main

import (
	"context"
	"crypto/rsa"
	"fmt"
	"go-dmtor/client"
	cfg "go-dmtor/config"
	"go-dmtor/crypto"
	"go-dmtor/logger"
	"os"
	"os/signal"
)

var log = logger.New()

var usage = "Usage: %s <srv|cli>\n"

func main() {
	// crypt_demo()

	// get input args
	args := os.Args
	if len(args) == 1 {
		log.Fatalf(usage, args[0])
	}
	arg := args[1]
	if arg != "srv" && arg != "cli" {
		log.Fatalf(usage, args[0])
	}

	ctx, cancel := context.WithCancel(context.Background())
	// start the server or connect
	cli := client.NewClient(ctx, cancel, cfg.ADDR)
	go func() {
		if arg == "srv" {
			err := cli.ServerStart()
			if err != nil {
				log.Fatalf("server start error: %v\n", err)
			}
		} else {
			err := cli.ServerConnect()
			if err != nil {
				log.Fatalf("server connect error: %v\n", err)
			}
		}
	}()

	// block and wait for user input
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Warnf("context done: %v\n", ctx.Err())
				return
			default:
				input := make([]byte, cfg.MSG_MAX_SIZE)
				n, err := os.Stdin.Read(input)
				if err != nil {
					log.Fatalf("read error: %v\n", err)
					return
				}
				input = input[:n]
				cli.SendMessage(input)
			}
		}
	}()

	// graceful shutdown
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt)
		<-stop
		cancel()
	}()

	<-ctx.Done()
	log.Warn("Bye!")
}

func crypt_demo() {
	key := crypto.Keygen()
	// Get the public key
	publicKey := &key.PublicKey

	// test PEM
	pubPem, err := crypto.PubToPem(publicKey)
	if err != nil {
		log.Fatalf("encode error: %v\n", err)
	}
	fmt.Printf("Public key pem: (%d) %x\n", len(pubPem), pubPem)
	pubFromPem, err := crypto.PemToPub(pubPem)
	if err != nil {
		log.Fatalf("decode error: %v\n", err)
	}

	// test pub bytes
	pubBytes, err := crypto.PubToBytes(publicKey)
	if err != nil {
		log.Fatalf("pub to bytes error: %v\n", err)
	}
	fmt.Printf("Public key bytes: (%d) %x\n", len(pubBytes), pubBytes)
	pubFromBytes, err := crypto.BytesToPub(pubBytes)
	if err != nil {
		log.Fatalf("bytes to pub error: %v\n", err)
	}

	crypt_demo_run(&key, pubFromPem, "hello world 1 - pem")
	crypt_demo_run(&key, pubFromBytes, "hello world 2 - bytes")

}

func crypt_demo_run(key *rsa.PrivateKey, pub *rsa.PublicKey, message string) {
	// Encrypt a message
	cipher := crypto.Encrypt([]byte(message), pub)
	fmt.Printf("Ciphertext: %x\n", cipher)

	// Decrypt the message
	plain := crypto.Decrypt(cipher, key)
	fmt.Printf("Plaintext: %s\n", plain)

}
