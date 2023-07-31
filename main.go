package main

import (
	"context"
	"go-dmtor/client"
	cfg "go-dmtor/config"
	myrsa "go-dmtor/cryptotools/rsa"
	"go-dmtor/logger"
	"go-dmtor/tor"
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

	ctx, cancel := context.WithCancel(context.Background())

	// create assym crypter for communication
	crypter, err := myrsa.New()
	if err != nil {
		log.Fatalf("cant create crypter: %v\n", err)
	}

	// start the server or connect
	cli := client.NewClient(ctx, cancel, cfg.ADDR, crypter)

	// TODO: add new session command and connect to old session.
	// or select a previous session from a list
	go func() {
		switch arg {
		case "srv":
			err := cli.ServerStart()
			if err != nil {
				log.Fatalf("server start error: %v\n", err)
			}
		case "cli":
			err := cli.ServerConnect()
			if err != nil {
				log.Fatalf("server connect error: %v\n", err)
			}
		case "tor":
			go func() {
				err := tor.Run(ctx)
				if err != nil {
					log.Fatalf("cant start tor: %v\n", err)
				}
			}()
		default:
			log.Fatalf(usage, args[0])
		}
	}()

	// start tor
	// if os.Getenv("TOR") == "1" {
	// 	go func() {
	// 		err := tor.Run(ctx)
	// 		if err != nil {
	// 			log.Fatalf("cant start tor: %v\n", err)
	// 		}
	// 	}()
	// }

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
				err = cli.SendMessage(input)
				if err != nil {
					log.Errorf("can't send a message: %v\n", err)
				}
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

// func aes_demo(pin string) {
// 	test := []byte("123test123123123123")
//
// 	// encode aes with password
// 	cipher, err := ct.AESencrypt(test, pin)
// 	if err != nil {
// 		log.Fatalf("aes encrypt error: %v\n", err)
// 	}
// 	fmt.Printf("cipher: %x\n", cipher)
//
// 	// decode back
// 	plain, err := ct.AESdecrypt(cipher, pin)
// 	if err != nil {
// 		log.Fatalf("aes decrypt error: %v\n", err)
// 	}
// 	fmt.Printf("plain: %s\n", plain)
// }
