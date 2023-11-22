package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/1F47E/go-shaihulud/internal/api"
	"github.com/1F47E/go-shaihulud/internal/client"
	"github.com/1F47E/go-shaihulud/internal/core"
	myrsa "github.com/1F47E/go-shaihulud/internal/cryptotools/rsa"
	"github.com/1F47E/go-shaihulud/internal/tui"

	"github.com/1F47E/go-shaihulud/internal/logger"

	zlog "github.com/rs/zerolog/log"
)

var log = logger.New()

var usage = "Usage: <srv | cli>\n"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

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

	eventsCh := make(chan tui.Event, 10000) // FIXME: delete this
	cli := client.NewClient(ctx, cancel, connType, crypter, eventsCh)

	core := core.NewCore(cli)

	api := api.NewApi(core)

	// gracefull shutdown
	go func() {
		<-ctx.Done()
		zlog.Info().Msg("Shutting down server...")
		if err := api.Shutdown(); err != nil {
			zlog.Error().Err(err).Msg("Error shutting down server")
		}
	}()

	endpoint := "localhost:3003"
	// open link in a new target window to ensure it opens every time the same
	// go func() {
	// 	time.Sleep(1 * time.Second)
	// 	err := open.Run("http://" + endpoint)
	// 	if err != nil {
	// 		log.Errorf("Failed to open browser: %v\n", err)
	// 	}
	// }()

	// start the server
	api.Start(endpoint)
	log.Warn("Bye!")
}
