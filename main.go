package main

import (
	"context"
	"go-dmtor/client"
	myrsa "go-dmtor/cryptotools/rsa"
	"go-dmtor/logger"
	"os"
	"os/signal"

	"golang.org/x/term"
)

var log = logger.New()

var usage = "Usage: <srv|cli>\n"

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
			// TODO: refactor to work with custom tor connection
			session := ""
			err := cli.RunServer(session)
			if err != nil {
				log.Fatalf("server start error: %v\n", err)
			}
		case "cli":
			// TODO: allow bypass auth for dev
			// get key as a param
			if len(args) != 3 {
				log.Fatalf("Usage: %s key <key>\n", args[0])
			}
			key := args[2]
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
		// // test auth key decoding
		// case "key":
		// 	// get key as a param
		// 	if len(args) != 3 {
		// 		log.Fatalf("Usage: %s key <key>\n", args[0])
		// 	}
		// 	key := args[2]

		// 	// ask for password from user input
		// 	log.Info("Enter password:")

		// 	password, err := term.ReadPassword(int(os.Stdin.Fd()))
		// 	if err != nil {
		// 		log.Fatalf("Error reading password: %v", err)
		// 	}

		// 	aes := myaes.New()
		// 	ath, err := auth.NewFromKey(aes, key, string(password))
		// 	if err != nil {
		// 		if strings.Contains(err.Error(), "authentication failed") {
		// 			log.Fatal("wrong password")
		// 		}
		// 		log.Fatalf("cant create auth: %v\n", err)
		// 	}
		// 	err = cli.TorConnectAsClient(ctx, ath.OnionAddressFull())
		// 	// will block
		// 	// err = tor.Connect(ctx, ath.OnionAddressFull())
		// 	if err != nil {
		// 		log.Fatalf("cant connect via tor to %s: %v\n", ath.OnionAddress(), err)
		// 	}
		// 	log.Infof("connected to %s\n", ath.OnionAddress())

		// 	// go listenInput(ctx)

		// // test tor connection, onion generator and auth
		// // TODO: move to client, start server
		// case "tor":
		// 	// check session param
		// 	session := ""
		// 	if len(args) == 3 {
		// 		session = args[2]
		// 	}

		// 	if session != "" {
		// 		log.Infof("loading session %s...\n", session)
		// 	} else {
		// 		log.Info("creating a new session...")
		// 	}

		// 	// create auth struct will password
		// 	// and give it to the user
		// 	crypter := myaes.New()
		// 	auth, err := auth.New(crypter, session)
		// 	if err != nil {
		// 		log.Fatalf("cant create auth: %v\n", err)
		// 	}
		// 	fmt.Printf("%s", auth)
		// 	log.Debugf("onion: %s\n", auth.OnionAddress())

		// 	// start tor with the onion key
		// 	log.Info("Starting tor, please wait. It can take a few minutes...")

		// 	err = cli.TorConnectAsServer(auth.Onion())
		// 	if err != nil {
		// 		log.Fatalf("connection reading error: %v\n", err)
		// 	}
		// 	// test connection
		// 	// listen to tor connection
		// 	// for {
		// 	// 	log.Debug("Waiting for new connection")
		// 	// 	conn, err := torconn.Accept()
		// 	// 	if err != nil {
		// 	// 		log.Fatal(err)
		// 	// 	}
		// 	// 	log.Debug("Got new connection")
		// 	// 	ip := conn.RemoteAddr().String()
		// 	// 	// connID := crypto.Hash([]byte(ip))
		// 	// 	log.Debugf("Connection open for %s\n", ip)
		// 	// 	time.Sleep(1 * time.Hour)
		// 	// }
		default:
			log.Fatal(usage)
		}
	}()

	// start tor
	// if os.Getenv("TOR") == "1" {
	// 	go func() {
	// 		err := tor.Run(ctx)
	// 		if err != nil {
	// 			log.Fatalf("cant start tor: %v\n", err)
	// 		}
	// 	}()
	// }

	// block and wait for user input
	// TODO: move to function. to not run if asking password
	// go func() {
	// 	for {
	// 		select {
	// 		case <-ctx.Done():
	// 			log.Warnf("context done: %v\n", ctx.Err())
	// 			return
	// 		default:
	// 			input := make([]byte, cfg.MSG_MAX_SIZE)
	// 			n, err := os.Stdin.Read(input)
	// 			if err != nil {
	// 				log.Fatalf("read error: %v\n", err)
	// 				return
	// 			}
	// 			input = input[:n]
	// 			err = cli.SendMessage(input)
	// 			if err != nil {
	// 				log.Errorf("can't send a message: %v\n", err)
	// 			}
	// 		}
	// 	}
	// }()

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
