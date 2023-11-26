package main

import (
	"bufio"
	"context"
	"os"
	"os/signal"
	"strings"

	"github.com/1F47E/go-shaihulud/internal/client"
	myrsa "github.com/1F47E/go-shaihulud/internal/cryptotools/asymmetric/rsa"
	"github.com/1F47E/go-shaihulud/internal/logger"

	"golang.org/x/term"
)

var log = logger.New()

var usage = "Usage: <srv | cli>\n"

func main() {

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
			log.Info("Enter chat key:")
			reader := bufio.NewReader(os.Stdin)
			key, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalf("Error reading chat key: %v", err)
			}
			key = strings.TrimSpace(key)
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
