package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-host"
	libnet "github.com/libp2p/go-libp2p-net"
	ma "github.com/multiformats/go-multiaddr"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func makeNewPeer(listenPort int) (host.Host, error) {

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
		libp2p.NoSecurity,
		libp2p.RandomIdentity,
		libp2p.DisableRelay(),
		//libp2p.EnableRelay(0),
	}

	verboseLog("My Context: ")
	verboseLog(context.Background())

	basicHost, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:

	for i := 0; i < len(basicHost.Addrs()); i++ {
		verboseLog(basicHost.Addrs()[i])

		addr := basicHost.Addrs()[i]
		fullAddr := addr.Encapsulate(hostAddr)
		verboseLog(fullAddr.String())

	}

	addr := basicHost.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr)
	log.Printf("I am %s\n", fullAddr)

	log.Printf("Now run \"./CoiNR -l %d -d %s\" on a different terminal\n", listenPort+1, fullAddr)

	return basicHost, nil

}

func handleStream(s libnet.Stream) {

	log.Println("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go readData(rw)
	go writeData(rw)
	go mineBlocks(rw)

	// stream 's' will stay open until you close it (or the other side closes it).
}

func writeData(rw *bufio.ReadWriter) {

	go func() {
		for {
			time.Sleep(5 * time.Second)
			mutex.Lock()
			bytes, err := json.Marshal(Blockchain)
			if err != nil {
				log.Println(err)
			}
			mutex.Unlock()

			mutex.Lock()
			rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
			rw.Flush()
			mutex.Unlock()

		}
	}()

	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		sendData = strings.Replace(sendData, "\n", "", -1)

		args := strings.Split(sendData, " ")

		if len(args) != 3 {
			log.Println("not enough arguments for transaction")
		}

		amt, err := strconv.ParseFloat(args[2], 64)

		if err != nil {
			log.Println("invalid amount")
		}

		newTrans := Taction{
			getPrivateKey(args[0]),
			getPrivateKey(args[1]),
			amt,
			time.Now().String(),
		}

		transactionValidator(newTrans)

	}

}

func readData(rw *bufio.ReadWriter) {

	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {

			chain := make([]Block, 0)
			if err := json.Unmarshal([]byte(str), &chain); err != nil {
				log.Fatal(err)
			}

			mutex.Lock()
			if len(chain) > len(Blockchain) {
				Blockchain = chain
				bytes, err := json.MarshalIndent(Blockchain, "", "  ")
				if err != nil {

					log.Fatal(err)
				}
				// Green console color: 	\x1b[32m
				// Reset console color: 	\x1b[0m
				fmt.Printf("\x1b[32m%s\x1b[0m> ", string(bytes))
			}
			mutex.Unlock()
		}
	}
}

/**
func startRelay(ba *basichost.BasicHost) *relay.AutoRelayHost{

	dis := discovery.Discoverer()

	 rh := relay.NewAutoRelayHost(context.Background(), ba, dis)

	 return rh
}

*/
