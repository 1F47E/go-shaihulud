package main

import (
	"fmt"
	"go-dmtor/crypto"
	"log"
	"net"
	"sync"
)

func startSever() {
	// open tcp connection to port 3000 and listen to incoming connections.
	// on connection print hello
	conn, err := net.ResolveTCPAddr("tcp4", ":3000")
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", conn)
	if err != nil {
		panic(err)
	}
	log.Println("Listening on port 3000")
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		w, err := conn.Write([]byte("hello"))
		if err != nil {
			panic(err)
		}
		log.Printf("Wrote %d bytes\n", w)
		for {
			// read 32 bytes from the connection
			bytes := make([]byte, 32)
			n, err := conn.Read(bytes)
			if err != nil {
				log.Printf("read error: %v\n", err)
			}
			log.Printf("Received: %d bytes:\n%s\n", n, bytes)
			// response with echo
			w, err = conn.Write(bytes)
			if err != nil {
				log.Printf("write error: %v\n", err)
			}
			log.Printf("Replied: %d bytes\n", w)
		}
		conn.Write([]byte("hello"))
		// read 32 bytes from the connection
		bytes := make([]byte, 32)
		conn.Read(bytes)
		log.Printf("Received: %s\n", bytes)
	}()
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	startSever()
	// crypt_demo()
	wg.Wait()
}

func crypt_demo() {

	key := crypto.Keygen()
	// Get the public key
	publicKey := &key.PublicKey

	// Encrypt a message
	message := "hello, world"
	cipher := crypto.Encrypt(message, publicKey)
	fmt.Printf("Ciphertext: %x\n", cipher)

	plain := crypto.Decrypt(cipher, &key)
	fmt.Printf("Plaintext: %s\n", plain)
}
