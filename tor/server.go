package tor

import (
	"context"
	"net"
	"os"

	"github.com/cretz/bine/tor"
	"github.com/cretz/bine/torutil/ed25519"

	"go-dmtor/interfaces"
)

func Run(ctx context.Context, onioner interfaces.Onioner) (net.Listener, error) {
	// Specify and create the fixed data directory
	// BUG: not working properly
	dataDir := "tor_data"
	_, err := os.Stat(dataDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(dataDir, 0700)
		if err != nil {
			return nil, err
		}
	}
	cfg := &tor.StartConf{RetainTempDataDir: true}

	t, err := tor.Start(ctx, cfg)
	if err != nil {
		return nil, err
	}

	keyBytes := onioner.PrivateKey()
	keyPair := ed25519.PrivateKey(keyBytes)

	onion, err := t.Listen(ctx, &tor.ListenConf{Key: keyPair, LocalPort: 3000, RemotePorts: []int{80}})
	if err != nil {
		return nil, err
	}
	return onion, nil
}
