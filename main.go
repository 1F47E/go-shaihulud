package main

import (
	"context"
	"crypto/rsa"
	"fmt"
	"go-dmtor/client"
	cfg "go-dmtor/config"
	ct "go-dmtor/cryptotools"
	"go-dmtor/logger"
	"go-dmtor/tor"
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

	ctx, cancel := context.WithCancel(context.Background())
	// start the server or connect
	cli := client.NewClient(ctx, cancel, cfg.ADDR)
	go func() {
		switch arg {
		case "srv":
			// err := cli.ServerStart()
			// if err != nil {
			// 	log.Fatalf("server start error: %v\n", err)
			// }
		case "cli":
			err := cli.ServerConnect()
			if err != nil {
				log.Fatalf("server connect error: %v\n", err)
			}
		case "demo":
			// o := "onion pub key"
			// h := ct.EncryptOnion(o)
			// fmt.Printf("connection key: %s\n", h)
		default:
			log.Fatalf(usage, args[0])
		}
	}()

	// start tor
	if os.Getenv("TOR") == "1" {
		go func() {
			err := tor.Run(ctx)
			if err != nil {
				log.Fatalf("cant start tor: %v\n", err)
			}
		}()
	}

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
	key := ct.Keygen()
	// Get the public key
	publicKey := &key.PublicKey

	// test PEM
	pubPem, err := ct.PubToPem(publicKey)
	if err != nil {
		log.Fatalf("encode error: %v\n", err)
	}
	fmt.Printf("Public key pem: (%d) %x\n", len(pubPem), pubPem)
	pubFromPem, err := ct.PemToPub(pubPem)
	if err != nil {
		log.Fatalf("decode error: %v\n", err)
	}

	// test pub bytes
	pubBytes, err := ct.PubToBytes(publicKey)
	if err != nil {
		log.Fatalf("pub to bytes error: %v\n", err)
	}
	fmt.Printf("Public key bytes: (%d) %x\n", len(pubBytes), pubBytes)
	pubFromBytes, err := ct.BytesToPub(pubBytes)
	if err != nil {
		log.Fatalf("bytes to pub error: %v\n", err)
	}

	crypt_demo_run(&key, pubFromPem, "hello world 1 - pem")
	crypt_demo_run(&key, pubFromBytes, "hello world 2 - bytes")

}

func crypt_demo_run(key *rsa.PrivateKey, pub *rsa.PublicKey, message string) {
	// Encrypt a message
	cipher := ct.Encrypt([]byte(message), pub)
	fmt.Printf("Ciphertext: %x\n", cipher)

	// Decrypt the message
	plain := ct.Decrypt(cipher, key)
	fmt.Printf("Plaintext: %s\n", plain)

}
