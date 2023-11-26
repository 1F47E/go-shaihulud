package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/1F47E/go-shaihulud/internal/api"
	"github.com/1F47E/go-shaihulud/internal/client"
	"github.com/1F47E/go-shaihulud/internal/core"
	myaes "github.com/1F47E/go-shaihulud/internal/cryptotools/aes"
	"github.com/1F47E/go-shaihulud/internal/cryptotools/auth"
	myrsa "github.com/1F47E/go-shaihulud/internal/cryptotools/rsa"
	"github.com/1F47E/go-shaihulud/internal/tui"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

const (
	portServerTransport = 3000
	portServerWeb       = 80
)

func main() {
	zlog.Logger = zlog.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// create assym crypter for communication
	crypter, err := myrsa.New()
	if err != nil {
		zlog.Fatal().Err(err).Msg("cant create crypter")
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

	// generate auth with onion domain

	session := ""
	// generate auth key and password
	crypterAes := myaes.New()
	ath, err := auth.New(crypterAes, session)
	if err != nil {
		zlog.Fatal().Err(err).Msg("cant create auth")
	}

	// auth, err := cli.GenerateAuth()
	// if err != nil {
	// 	zlog.Fatal().Err(err).Msg("cant generate auth")
	// }
	// auth creds for the client
	zlog.Debug().Msgf("auth key: \n%s\n", ath.AccessKey())
	zlog.Debug().Msgf("password: %s\n", ath.Password())

	// println()
	zlog.Warn().Msg("🔑 Client auth creds")
	zlog.Warn().Msg("=======================================")
	zlog.Warn().Msgf(" Key: %s\n\n", ath.AccessKey())
	zlog.Warn().Msgf(" Password: %s\n", ath.Password())
	zlog.Warn().Msg("=======================================")
	println()

	// get address
	address := ""
	msgLoading := ""
	msgSuccess := ""
	switch connType {
	case client.Local:
		address = "localhost:3000"
		msgLoading = fmt.Sprintf("Starting local server on %s", address)
		msgSuccess = fmt.Sprintf("Local server started at %s, waiting for incoming connections...", address)
	case client.Tor:
		msgLoading = "Starting TOR..."
		msgSuccess = "Tor server started, waiting for incoming connections..."
		address = ath.OnionAddressTransport()
		zlog.Debug().Msgf("starting tor, onion address: %v\n", address)
	default:
		zlog.Fatal().Msgf("unknown connection type: %v\n", connType)
	}
	zlog.Info().Msg(msgLoading)

	// onionEndpointTransport := ath.OnionAddressTransport()
	// onionEndpointWeb := ath.OnionAddressWeb()
	onionKeys := ath.Onion().PrivKey()
	listnerTransport, err := cli.CreateListner(portServerTransport, onionKeys)
	if err != nil {
		zlog.Fatal().Err(err).Msg("cant create listner")
	}

	cli.RunListner(listnerTransport)
	zlog.Info().Msg(msgSuccess)

	zlog.Debug().Msg("starting api server")
	listnerWeb, err := cli.CreateListner(portServerWeb, onionKeys)
	if err != nil {
		zlog.Fatal().Err(err).Msg("cant create listner")
	}

	core := core.NewCore(cli)
	// auth.Onion().PrivKey()

	// listener, err := cli.connector.RunServer(address, )
	// if err != nil {
	// 	return nil, err
	// }

	api := api.NewApi(core, listnerWeb)

	// gracefull shutdown
	go func() {
		<-ctx.Done()
		zlog.Info().Msg("Shutting down server...")
		if err := api.Shutdown(); err != nil {
			zlog.Error().Err(err).Msg("Error shutting down server")
		}
	}()

	// endpoint := "localhost:8080"
	// open link in a new target window to ensure it opens every time the same
	// go func() {
	// 	time.Sleep(1 * time.Second)
	// 	err := open.Run("http://" + endpoint)
	// 	if err != nil {
	// 		log.Errorf("Failed to open browser: %v\n", err)
	// 	}
	// }()

	// start the server
	api.Start()
	zlog.Info().Msg("Server stopped")
}
