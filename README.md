```
 _____  ___   _   _______ _    _  ______________  ___
/  ___|/ _ \ | \ | |  _  \ |  | ||  _  | ___ \  \/  |
\ `--./ /_\ \|  \| | | | | |  | || | | | |_/ / .  . |
 `--. \  _  || . ` | | | | |/\| || | | |    /| |\/| |
/\__/ / | | || |\  | |/ /\  /\  /\ \_/ / |\ \| |  | |
\____/\_| |_/\_| \_/___/  \/  \/  \___/\_| \_\_|  |_/
                                                     

```
# go-tor-chat-wip
WIP p2p chat via tor network
With keys exchange and RSA encryption

# Connection flow
```
Client A starts the chat server
Client A shares the access key and password with Client B
Client B enters the access key and password to decrypt the onion address
Client B connects to the onion address
Clients exchange public keys
Clients encrypt messages with each other's public keys
```


# Access key and password
```
The access key is a readable binary key in hex format
that resembles 1234-ABCD-EFGH-5678....
This key represents the AES-encrypted onion address.

The password is used to encrypt/decrypt the access key to obtain the onion address,
which takes a form like 1234-ABCD.

The password consists of random bytes converted to upper-case hex format.

It is also used to sign messages via HMAC to verify message integrity.

Workflow:

User A (server), after connecting to Tor and generating an onion address, encrypts this address with a randomly generated password.

User A then shares the access key (AES-encrypted onion address) and password with User B.

The password and access key should be shared via different channels for security.

User B enters the access key and then the password to decrypt the onion address.

                                                                       
+----------------+                                                        +-------------+
|     User A     |                                                        |    User B   |
|    (Server)    |                                                        |   (Client)  |
+-------+--------+                                                        +--------+----+
        |                                                                          |
        |                                                                          |
        |<------------>Connects to Tor Network                                     |
        |                                                                          |
        |--->Generates random password, encrypts onion address with password       |
        |                                                                          |
        |--->Generates access key (AES-encrypted onion address)                    |
        |                                                                          |
        |-------------------------- Shared Access Key ---------------------------->|
        |                                                                          |
        |-------------------- Shares password via Channel 2 ---------------------->|
        |                                                                          |
        |                                                                          |
        |        Enters access key and password to decrypt onion address <---------|
        |                                                                          |
        |                            Decrypts the key with the password  <---------|
        |                                                                          |
        |                                        Connects to Tor Network <---------|
        |                                                                          |
        |                                             Connects to User A <---------|
        |                                                                          |
        |<------------------------- RSA pub key exchange ------------------------->|
        |                                                                          |
        |<---------- Users verify message integrity via HMAC signature ----------->|
        |                                                                          |
+-------+--------+                                                        +--------+----+
|     User A     |                                                        |    User B   |
|    (Server)    |                                                        |   (Client)  |
+----------------+                                                        +-------------+
```


# TODO
- [ ] add timestamps to the messages to prevent replay attacks
- [ ] sign every message with hmac to verify integrity and prevent MITM attacks
- [ ] do onion routing
- [x] gen chat key for access, hide onion
- [ ] cli on start generate key with password
- [ ] ack on handshake received
- [ ] ack on every message
- [ ] chat gui
- [ ] send files
- [ ] allow multiple users in a chat room
- [ ] test coverage for all packages
- [x] test coverage for crypto packages
- [x] basic tcp echo server
- [x] basic chat server-client
- [x] custom protocol with header
- [x] chat via custom protocol
- [x] handshake on connection, exchange public keys
- [x] encrypt chat with public key
- [x] basic tor tcp connection
