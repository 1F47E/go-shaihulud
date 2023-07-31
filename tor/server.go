package tor

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/cretz/bine/tor"
	"github.com/cretz/bine/torutil/ed25519"

	myonion "go-dmtor/cryptotools/onion"
	"go-dmtor/logger"
)

var log = logger.New()

func Run(ctx context.Context, session string) error {
	log.Info("Starting tor, please wait. It can take a few minutes...")
	t, err := tor.Start(ctx, nil)
	if err != nil {
		return err
	}
	defer t.Close()

	// Wait at most a few minutes to publish the service
	listenCtx, listenCancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer listenCancel()

	var keyPair ed25519.KeyPair

	// no session - generate new onion and save it
	if session == "" {
		log.Info("no session file provided, generating new session...")
		// read key from session file
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
		keyPair, err = ed25519.GenerateKey(rnd)
		if err != nil {
			return err
		}
		// TODO: not save untill tor confirms it
	} else {
		// load onion from the session file
		privKeyBytes, err := myonion.PrivKeyRead(session)
		if err != nil {
			return err
		}
		keyPair = ed25519.PrivateKey(privKeyBytes)
		log.Infof("session loaded - %s\n", session)
	}

	onion, err := t.Listen(listenCtx, &tor.ListenConf{Key: keyPair, LocalPort: 3000, RemotePorts: []int{80}})
	if err != nil {
		return err
	}
	defer onion.Close()
	log.Infof("Session started - %s\n", onion.ID)

	// Save priv key to a session file
	if session == "" {
		privKeyBytes := keyPair.PrivateKey()
		addr, err := myonion.PrivKeyBytesToOnionAddress(privKeyBytes)
		if err != nil {
			return fmt.Errorf("failed to get onion address from priv key: %s", err)
		}
		if addr != onion.ID {
			return fmt.Errorf("onion address mismatch: %s != %s", addr, onion.ID)
		}
		// save as session file
		err = myonion.PrivKeySave(privKeyBytes)
		if err != nil {
			return fmt.Errorf("failed to save session: %s", err)
		}
		log.Infof("new session created and saved - %s\n", addr)
	} else {
		log.Infof("new session started - %s\n", session)
	}

	// TODO: return pointer to connection

	// test connection
	for {
		log.Debug("Waiting for new connection")
		conn, err := onion.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Debug("Got new connection")
		ip := conn.RemoteAddr().String()
		// connID := crypto.Hash([]byte(ip))
		log.Debugf("Connection open for %s\n", ip)
		time.Sleep(1 * time.Hour)
	}
}
