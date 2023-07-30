package tor

import (
	"context"
	"go-dmtor/logger"
	"time"

	"github.com/cretz/bine/tor"
)

var log = logger.New()

// address should be in the form of onion:port
// 6hiqm5iky3yow7akjhwgdfqkrd7shsz3smw4shihaksh5n3jc34fodid.onion:80
func Connect(ctx context.Context, addr string) error {
	// Start tor with default config (can set start conf's DebugWriter to os.Stdout for debug logs)
	log.Info("Starting tor...")
	t, err := tor.Start(ctx, nil)
	if err != nil {
		return err
	}
	defer t.Close()
	// Wait at most a minute to start network and get
	dialCtx, dialCancel := context.WithTimeout(context.Background(), time.Minute)
	defer dialCancel()
	// Make connection
	dialer, err := t.Dialer(dialCtx, nil)
	if err != nil {
		return err
	}

	log.Infof("Connecting to %s\n", addr)
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		return err
	}
	log.Info("Connected")

	// test tcp connection
	for {
		bytes := make([]byte, 1024)
		n, err := conn.Read(bytes)
		if err != nil {
			log.Fatalf("Listner: Read error: %v\n", err)
		}
		log.Printf("Received: %d bytes:\n\n", n)
		log.Printf("%s\n", bytes[:n])
	}

	return nil
}
