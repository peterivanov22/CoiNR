package main

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/kr/pretty"
	"github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	"log"
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

var publicKey = ""
var privateKey = generateKeys()


func getPrivateKey () *ecdsa.PrivateKey {
	return privateKey
}


//A list of transactions we have yet to process
var pendingTransactions []Taction


func main() {



	publicKey := getPublicKey(privateKey)

	log.Println(privateKey)

	//so the BPMs is simply the data of the block
	currtime := time.Now()
	genesisBlock := Block{0, currtime.String(), []Taction{}, "", "", 0, ""}
	genesisBlock.Hash = genesisBlock.calculateHash()

	Blockchain = append(Blockchain, genesisBlock)

	//who needs logging

	// We dont most of these command line arguemtns
	listenF := flag.Int("l", 0, "wait for incoming connections")
	target := flag.String("d", "", "target peer to dial")
	verbose := flag.Bool("v", false, "turn on verbose logging")
	//pubKey := flag.String("p", "", "public key for the user.")

	flag.Parse()

	if *listenF == 0 {
		log.Fatal("Please provide a port to bind on with -l")
	}
/*
	if *pubKey == "" {
		log.Fatal("Please provide a public key with -p")
	} else {
		publicKey = *pubKey
	}
*/
	verboseMode = *verbose

	// Make a host that listens on the given multiaddress
	ha, err := makeNewPeer(*listenF)
	if err != nil {
		log.Fatal(err)
	}

	//rh := startRelay(ha)

	log.Println("This hosts address is: " + publicKey)

	//we dont want this first part
	if *target == "" {
		log.Println("listening for connections")
		// Set a stream handler on host A. /p2p/1.0.0 is
		// a user-defined protocol name.
		ha.SetStreamHandler("/p2p/1.0.0", handleStream)

		select {} // hang forever
		/**** This is where the listener code ends ****/
	} else {
		//I need to set this up to work with correct hostnames
		ha.SetStreamHandler("/p2p/1.0.0", handleStream)

		// The following code extracts target's peer ID from the
		// given multiaddress
		ipfsaddr, err := ma.NewMultiaddr(*target)
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


		s, err := ha.NewStream(context.Background(), peerid, "/p2p/1.0.0")
		if err != nil {
			log.Fatalln(err)
		}
		// Create a buffered stream so that read and writes are non blocking.
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		// Create a thread to read and write data.
		go writeData(rw)
		go readData(rw)
		go mineBlocks(rw)


		// make a new stream from host B to host A
		// it should be handled on host A by the handler we set above because
		// we use the same /p2p/1.0.0 protocol


		select {} // hang forever

	}

}

func mineBlocks(rw *bufio.ReadWriter){

	for {

		if len(pendingTransactions) > 1 {

			validTrans := filterCommittedTactions(pendingTransactions)

			newBlock := generateBlock(Blockchain[len(Blockchain)-1], validTrans, difficulty)

			if newBlock.validate(&Blockchain[len(Blockchain)-1]) {
				mutex.Lock()
				Blockchain = append(Blockchain, newBlock)
				mutex.Unlock()
			}

			pendingTransactions = filterCommittedTactions(pendingTransactions)

			bytes, err := json.Marshal(Blockchain)
			if err != nil {
				log.Println(err)
			}

			spew.Dump(Blockchain)

			mutex.Lock()
			rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
			rw.Flush()
			mutex.Unlock()

		}

		time.Sleep(time.Second)
	}

}

/*
func getPrivateKey(pubKey string) string{


	return pubKey + "privateKey"
}


func getUserWallet(privKey string) float64{

	var wallet float64

	wallet = 0

	for _, ablock := range Blockchain{
		wallet += ablock.getWalletAmt(privKey)

	}

	return wallet
}
*/

func verboseLog(message interface{}){
	if verboseMode == true{
		log.Println(pretty.Sprint(message))
	}

}

func filterCommittedTactions(tactionList []Taction) []Taction{

	var newList []Taction

	for _, atact := range tactionList{

		found := false

		for _, ablock := range Blockchain{
			if ablock.hasTransaction(atact){
				found = true
			}

		}

		if !found {
			newList = append(newList, atact)
		}

	}

	return newList

}


/*
func transactionValidator(t Taction){

	payerWallet := getUserWallet(t.PrivateKey1)

	if t.isValid(payerWallet){
		pendingTransactions = append(pendingTransactions, t)
	} else {
		log.Println("---   Invalid Transaction.   ---\n" +
			"Payer " + t.PrivateKey1 + " does not have " + strconv.FormatFloat(t.Amount, 'E', -1, 64) + " CoiNR to spend.\n" +
			"That user only has " + strconv.FormatFloat(payerWallet, 'E', -1, 64) + " in their wallet")
	}
}
*/

