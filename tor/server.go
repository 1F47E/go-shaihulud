package tor

import (
	"context"
	"net"

	"github.com/cretz/bine/tor"
	"github.com/cretz/bine/torutil/ed25519"

	"go-dmtor/interfaces"
)

func Run(ctx context.Context, onioner interfaces.Onioner) (net.Listener, error) {
	// Specify and create the fixed data directory
	// BUG: -- settings tor dir not working properly
	// dataDir := "tor_data"
	// _, err := os.Stat(dataDir)
	// if os.IsNotExist(err) {
	// 	err = os.Mkdir(dataDir, 0700)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }
	// cfg := &tor.StartConf{DataDir: dataDir, NoHush: true}

	// TODO: add cleanup tor forlders on exit

	t, err := tor.Start(ctx, nil)
	if err != nil {
		return nil, err
	}

	keyBytes := onioner.PrivKey()
	keyPair := ed25519.PrivateKey(keyBytes)

	onion, err := t.Listen(ctx, &tor.ListenConf{Key: keyPair, LocalPort: 3000, RemotePorts: []int{80}})
	if err != nil {
		return nil, err
	}
	return onion, nil
}
