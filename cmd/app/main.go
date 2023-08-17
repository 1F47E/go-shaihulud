package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/1F47E/go-shaihulud/client"
	myrsa "github.com/1F47E/go-shaihulud/cryptotools/rsa"

	// "github.com/1F47E/go-shaihulud/gui"
	"github.com/1F47E/go-shaihulud/logger"
	"github.com/1F47E/go-shaihulud/tui"
)

var log = logger.New()

var usage = "Usage: <srv | cli key>\n"

func main() {

	// get input args
	args := os.Args
	if len(args) == 1 {
		log.Fatal(usage)
	}
	arg := args[1]

	ctx, cancel := context.WithCancel(context.Background())

	// TUI init (events channel for tui status updates)
	eventsCh := make(chan tui.Event)
	t := tui.New(ctx, eventsCh)
	t.RenderChat()
	panic("test")

	go t.Listner()
	go t.RenderLoader()

	//
	// time.Sleep(1 * time.Second)
	// eventsCh <- tui.NewEventSpin("loading tor...")
	// time.Sleep(3 * time.Second)
	// eventsCh <- tui.NewEventAccess("key", "password")
	// time.Sleep(3 * time.Second)
	// panic("test")

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
	cli := client.NewClient(ctx, cancel, connType, crypter, eventsCh)

	// TODO: add new session command and connect to old session.
	// or select a previous session from a list
	go func() {
		switch arg {
		case "srv":
			eventsCh <- tui.NewEventSpin("Loading...")
			session := ""
			err := cli.RunServer(session)
			if err != nil {
				log.Fatalf("server start error: %v\n", err)
			}
		case "cli":
			key, password, err := t.RenderAuth()
			if err != nil {
				log.Fatalf("auth error: %v\n", err)
			}
			fmt.Printf("key: %s\n", key)
			fmt.Printf("password: %s\n", password)

			// TODO: validate password first before connecting,
			// expose crypter func to validate password
			eventsCh <- tui.NewEventSpin("Connecting...")
			err = cli.RunClient(key, password)
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
