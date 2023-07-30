package tor

import (
	"context"
	"math/rand"
	"os"
	"time"

	"github.com/cretz/bine/tor"
	"github.com/cretz/bine/torutil/ed25519"

	// "golang.org/x/crypto/ed25519"
	ct "go-dmtor/cryptotools"
)

func Run(ctx context.Context) error {
	// Start tor with default config (can set start conf's DebugWriter to os.Stdout for debug logs)
	log.Info("Starting tor...")
	t, err := tor.Start(ctx, nil)
	if err != nil {
		return err
	}
	defer t.Close()

	// Wait at most a few minutes to publish the service
	listenCtx, listenCancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer listenCancel()

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	keyPair, err := ed25519.GenerateKey(rnd)
	if err != nil {
		return err
	}

	onion, err := t.Listen(listenCtx, &tor.ListenConf{Key: keyPair, LocalPort: 3000, RemotePorts: []int{80}})
	if err != nil {
		return err
	}
	defer onion.Close()
	log.Infof("onion created - %s.onion\n", onion.ID)
	// get public key bytes
	pubKeyBytes := keyPair.Public().(ed25519.PublicKey)
	// save to file to debug
	_ = os.WriteFile("debug_onion.pub", pubKeyBytes, 0644)
	hexkey := ct.KeyFromOnionPubKey(pubKeyBytes)
	log.Warnf("connection key: %s\n", hexkey)
	// decode back and compare
	onionDecoded, err := ct.KeyToOnionAddress(hexkey)
	if err != nil {
		log.Fatalf("onion decode error: %v\n", err)
	}
	log.Warnf("onion decoded: %s\n", onionDecoded)
	if onionDecoded != onion.ID {
		log.Fatalf("onion decode error: %v != %v\n", onionDecoded, onion.ID)
	}

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
