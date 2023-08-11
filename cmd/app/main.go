package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/1F47E/go-shaihulud/pkg/client"
	myrsa "github.com/1F47E/go-shaihulud/pkg/cryptotools/rsa"
	"github.com/1F47E/go-shaihulud/pkg/gui"
	"github.com/1F47E/go-shaihulud/pkg/logger"

	"golang.org/x/term"
)

var log = logger.New()

var usage = "Usage: <srv | cli key>\n"

func main() {
	if os.Getenv("TUI") == "1" {
		gui.Draw()
	}

	// get input args
	args := os.Args
	if len(args) == 1 {
		log.Fatal(usage)
	}
	arg := args[1]

	ctx, cancel := context.WithCancel(context.Background())

	// create assym crypter for communication
	crypter, err := myrsa.New()
	if err != nil {
		log.Fatalf("cant create crypter: %v\n", err)
	}

	// start the server or connect
	var connType client.ConnectionType
	if os.Getenv("TOR") == "0" {
		connType = client.Local
	} else {
		connType = client.Tor
	}
	cli := client.NewClient(ctx, cancel, connType, crypter)

	// TODO: add new session command and connect to old session.
	// or select a previous session from a list
	go func() {
		switch arg {
		case "srv":
			session := ""
			err := cli.RunServer(session)
			if err != nil {
				log.Fatalf("server start error: %v\n", err)
			}
		case "cli":
			// TODO: allow bypass auth for dev
			// get key as a param
			if len(args) != 3 {
				log.Fatalf("Usage: %s key <key>\n", args[0])
			}
			key := args[2]
			// TODO: validate key

			log.Info("Enter password:")
			password, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				log.Fatalf("Error reading password: %v", err)
			}

			err = cli.RunClient(key, string(password))
			if err != nil {
				log.Fatalf("server connect error: %v\n", err)
			}

		default:
			log.Fatal(usage)
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
