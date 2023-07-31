package main

import (
	"context"
	"go-dmtor/client"
	cfg "go-dmtor/config"
	myaes "go-dmtor/cryptotools/aes"
	"go-dmtor/cryptotools/auth"
	"go-dmtor/cryptotools/onion"
	myrsa "go-dmtor/cryptotools/rsa"
	"go-dmtor/interfaces"
	"go-dmtor/logger"
	"go-dmtor/tor"
	"os"
	"os/signal"
	"time"
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
			// check session param
			session := ""
			if len(args) == 3 {
				session = args[2]
			}

			// load a session
			// no session - generate new onion and save it
			var onioner interfaces.Onioner
			if session == "" {
				log.Info("no session file provided, generating new session...")
				o, err := onion.New()
				if err != nil {
					log.Fatalf("cant create onion: %v\n", err)
				}
				onioner = o
				log.Infof("session created - %s\n", onioner.Address())
			} else {
				// load onion from the session file
				o, err := onion.NewFromSession(session)
				if err != nil {
					log.Fatalf("cant load session: %v\n", err)
				}
				onioner = o
				log.Infof("session loaded - %s\n", onioner.Address())
			}
			log.Debugf("onioner loaded: %v\n", onioner)

			// run tor with the key
			log.Info("Starting tor, please wait. It can take a few minutes...")
			torconn, err := tor.Run(ctx, onioner)
			if err != nil {
				log.Fatalf("cant start tor: %v\n", err)
			}
			defer torconn.Close()
			log.Infof("Session started - %s\n", onioner.Address())

			// save session to a file if it was created
			if session == "" {
				err = onioner.Save()
				if err != nil {
					log.Fatalf("cant save session: %v\n", err)
				}
			}

			// create auth struct will password
			// and give it to the user
			crypter := myaes.New()
			auth := auth.New(crypter, onioner)
			log.Warnf("%s", auth)

			// test connection
			// listen to tor connection
			for {
				log.Debug("Waiting for new connection")
				conn, err := torconn.Accept()
				if err != nil {
					log.Fatal(err)
				}
				log.Debug("Got new connection")
				ip := conn.RemoteAddr().String()
				// connID := crypto.Hash([]byte(ip))
				log.Debugf("Connection open for %s\n", ip)
				time.Sleep(1 * time.Hour)
			}
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
