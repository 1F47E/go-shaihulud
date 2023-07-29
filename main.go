package main

import (
	"context"
	"go-dmtor/client"
	"go-dmtor/client/message"
	cfg "go-dmtor/config"
	"go-dmtor/logger"
	"os"
	"os/signal"
)

var log = logger.New()

var usage = "Usage: %s <srv|cli>\n"

func main() {
	// get input args
	args := os.Args
	if len(args) == 1 {
		log.Fatalf(usage, args[0])
	}
	arg := args[1]
	if arg != "srv" && arg != "cli" {
		log.Fatalf(usage, args[0])
	}

	// context for graceful shutdown and exit
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt)
		<-stop
		cancel()
	}()

	// start the client cli or server
	cli := client.NewClient(ctx, cancel, cfg.ADDR)
	if arg == "srv" {
		err := cli.ServerStart()
		if err != nil {
			log.Fatalf("server start error: %v\n", err)
		}
		m := message.NewMessageHello()
		cli.MsgCh <- *m
	} else {
		err := cli.ServerConnect()
		if err != nil {
			log.Fatalf("server connect error: %v\n", err)
		}
		m := message.NewMessageHello()
		cli.MsgCh <- *m
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
				m := message.NewMessageText(string(input[:n]))
				cli.MsgCh <- *m
			}
		}
	}()
	<-ctx.Done()
	log.Warn("Bye!")
}

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
