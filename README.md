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

# TODO
- [ ] do onion routing
- [ ] gen chat key for access, hide onion
- [ ] cli on start generate key with pin
- [ ] ack on message received
- [ ] chat gui
- [ ] send files
- [ ] allow multiple users in a chat room
- [x] basic tcp echo server
- [x] basic chat server-client
- [x] custom protocol with header
- [x] chat via custom protocol
- [x] handshake on connection, exchange public keys
- [x] encrypt chat with public key
- [x] basic tor tcp connection
