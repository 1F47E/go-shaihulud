package main

import (
	"context"
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
	pin := ct.AccessPinGenerate()
	fmt.Printf("pin: %s\n", pin)
	aes_demo(pin)
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

func aes_demo(pin string) {
	test := []byte("123test123123123123")

	// encode aes with password
	cipher, err := ct.AESencrypt(test, pin)
	if err != nil {
		log.Fatalf("aes encrypt error: %v\n", err)
	}
	fmt.Printf("cipher: %x\n", cipher)

	// decode back
	plain, err := ct.AESdecrypt(cipher, pin)
	if err != nil {
		log.Fatalf("aes decrypt error: %v\n", err)
	}
	fmt.Printf("plain: %s\n", plain)
}

func rsa_demo() {
	key := ct.Keygen()
	// Get the public key
	publicKey := &key.PublicKey

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

	message := "Hello World!"

	// Encrypt a message
	cipher := ct.MessageEncrypt([]byte(message), pubFromBytes)
	fmt.Printf("Ciphertext: %x\n", cipher)

	// Decrypt the message
	plain := ct.MessageDecrypt(cipher, &key)
	fmt.Printf("Plaintext: %s\n", plain)

}
