package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/kr/pretty"
	"github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

// This is the main class

// Not sure this is where this should go. Storing the global difficulty of the blocks.
const difficulty = 1

// A block chain is just a slice of blocks
var Blockchain []Block

// seems like we need a mutex
var mutex = &sync.Mutex{}

var verboseMode = false

var hostList []string

func main() {

	//TODO get this to work on local vagrant machines with a hostfile

	readHostfile()

	//so the BPMs is simply the data of the block
	currtime := time.Now()
	genesisBlock := Block{0, currtime.String(), 0, "", "", difficulty, ""}
	genesisBlock.Hash = genesisBlock.calculateHash()

	Blockchain = append(Blockchain, genesisBlock)

	//who needs logging

	// We dont most of these command line arguemtns
	listenF := flag.Int("l", 0, "wait for incoming connections")
	target := flag.String("d", "", "target peer to dial")
	verbose := flag.Bool("v", false, "turn on verbose logging")
	flag.Parse()

	if *listenF == 0 {
		log.Fatal("Please provide a port to bind on with -l")
	}

	verboseMode = *verbose

	// Make a host that listens on the given multiaddress
	ha, err := makeNewPeer(*listenF)
	if err != nil {
		log.Fatal(err)
	}

	/**

	//we dont want this first part
	if *target == "" {
		log.Println("listening for connections")
		// Set a stream handler on host A. /p2p/1.0.0 is
		// a user-defined protocol name.
		ha.SetStreamHandler("/p2p/1.0.0", handleStream)

		select {} // hang forever

	} else {

		peerLogic(ha, target)

	}

	*/

	peerLogic(ha, target)

}

func peerLogic(ha host.Host, targ *string) {

	if *targ != "" {
		hostList = append(hostList, *targ)
	}

	for _, target := range hostList {

		verboseLog(target)

		//I need to set this up to work with correct hostnames
		ha.SetStreamHandler("/p2p/1.0.0", handleStream)

		// The following code extracts target's peer ID from the
		// given multiaddress
		ipfsaddr, err := ma.NewMultiaddr(target)
		if err != nil {
			log.Fatalln(err)
		}

		verboseLog("Ipfsaddr:" + ipfsaddr.String())

		pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
		if err != nil {
			log.Fatalln(err)
		}

		verboseLog("pid: " + pid)

		peerid, err := peer.IDB58Decode(pid)
		if err != nil {
			log.Fatalln(err)
		}

		verboseLog("peerid: " + peerid.String())

		if ha.ID() == peerid{
			verboseLog("Skipping self-dial")
			continue
		}

		// Decapsulate the /ipfs/<peerID> part from the target
		// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
		targetPeerAddr, _ := ma.NewMultiaddr(
			fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
		targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

		verboseLog(targetAddr)

		// We have a peer ID and a targetAddr so we add it to the peerstore
		// so LibP2P knows how to contact it
		ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

		log.Println("opening stream")
		// make a new stream from host B to host A
		// it should be handled on host A by the handler we set above because
		// we use the same /p2p/1.0.0 protocol
		s, err := ha.NewStream(context.Background(), peerid, "/p2p/1.0.0")
		if err != nil {
			log.Println(err)
		} else {
			// Create a buffered stream so that read and writes are non blocking.
			rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

			// Create a thread to read and write data.
			go writeData(rw)
			go readData(rw)
		}


	}

	log.Println("listening for connections")
	// Set a stream handler on host A. /p2p/1.0.0 is
	// a user-defined protocol name.
	ha.SetStreamHandler("/p2p/1.0.0", handleStream)

	select {} // hang forever

}

func readHostfile() {

	csvFile, _ := os.Open("hosts.txt")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	linecount := 1
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error Reading from file.")
			log.Fatal(err)
		}

		apeer := "/ip4/" + line[0] + "/tcp/" + line[1] + "/ipfs/" + line[2]

		verboseLog(apeer)
		hostList = append(hostList, apeer)

		linecount++
	}
}

func verboseLog(message interface{}) {
	if verboseMode == true {
		log.Println(pretty.Sprint(message))
	}

}
