package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/1F47E/go-shaihulud/client"
	myrsa "github.com/1F47E/go-shaihulud/cryptotools/rsa"

	"github.com/1F47E/go-shaihulud/logger"
	"github.com/1F47E/go-shaihulud/tui"
)

var log = logger.New()

var usage = "Usage: <srv | cli>\n"

func main() {

	// ctx, cancel := context.WithCancel(context.Background())
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// TUI init (events channel for tui status updates)
	eventsCh := make(chan tui.Event)
	t := tui.New(ctx, eventsCh)
	go t.Listner()
	go t.RenderChat()

	<-ctx.Done()
	println("Bye")
	os.Exit(0)

	// get input args
	args := os.Args
	if len(args) == 1 {
		log.Fatal(usage)
	}
	arg := args[1]

	// time.Sleep(3 * time.Second)
	// t.SetMode(t.Mode)
	// panic("test")

	// go t.RenderLoader()

	//
	// time.Sleep(1 * time.Second)
	// eventsCh <- tui.NewEventSpin("loading tor...")
	// time.Sleep(3 * time.Second)
	// eventsCh <- tui.NewEventAccess("key", "password")
	// time.Sleep(10 * time.Second)
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

	// TODO: move connection logic our of main
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
			// TODO: auth disabled for debug
			// key, password, err := t.RenderAuth()
			// if err != nil {
			// 	log.Fatalf("auth error: %v\n", err)
			// }
			// fmt.Printf("key: %s\n", key)
			// fmt.Printf("password: %s\n", password)
			// key, password := "", ""
			// err := cli.AuthVerify(key, password)
			// if err != nil {
			// 	log.Errorf("")
			// 	os.Exit(0)
			// }

			// TODO: validate password first before connecting,
			// expose crypter func to validate password
			eventsCh <- tui.NewEventSpin("Connecting...")
			err = cli.RunClient()
			if err != nil {
				log.Error(err)
				os.Exit(0)
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
