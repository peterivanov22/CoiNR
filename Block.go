package CoiNR

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
<<<<<<< HEAD
	"io"
	"io/ioutil"
	"log"
	mrand "math/rand"
	"net"
	"os"
=======
>>>>>>> 1bb2d1d76577d2f46f7df899a0dee7e668d7819f
	"strconv"
	"strings"
	"time"
)

const difficulty = 1

// Block struct

type Block struct {
<<<<<<< HEAD
	Index        int
	Timestamp    string
	BPM          int
	Hash         string
	PreviousHash string
	Difficulty   int
	//in real bitcoin, nonce is 4 bytes
	//string in golang is pointer (of size 8 bytes)
	Nonce string
=======
	Index      int
	Timestamp  string
	BPM        int
	Hash       string
	PrevHash   string
	Difficulty int
	Nonce      string
>>>>>>> 1bb2d1d76577d2f46f7df899a0dee7e668d7819f
}

// validates the block.  Usage is block.validate()

func (b *Block) validate(previous *Block) bool {

	if previous.Index+1 != b.Index {
		return false
	}

	if previous.Hash != b.PrevHash {
		return false
	}

	if b.calculateHash() != b.Hash {
		return false
	}

	return true

}

//calculates the hash for the block.  Usage is block.calculateHash()

func (b *Block) calculateHash() string {

<<<<<<< HEAD
	hashString := strconv.Itoa(b.Index) + b.Timestamp + strconv.Itoa(b.BPM) + b.PreviousHash + b.Nonce
=======
	hashString := strconv.Itoa(b.Index) + b.Timestamp + strconv.Itoa(b.BPM) + b.PrevHash
>>>>>>> 1bb2d1d76577d2f46f7df899a0dee7e668d7819f

	// I looked it up, this is one of two built-in hash functions in the go standard crypto lib.
	// This one seems good.
	hasher := sha256.New()

	hasher.Write([]byte(hashString))

	hashed := hasher.Sum(nil)

	return hex.EncodeToString(hashed)

}

// generates a new block from the previous block and the BPM.  Usage is var newBlock = generateBlock(*oldBlock, bpm)

/*func generateBlock(prev *Block, BPM int) Block {

	var newBlock Block

	t := time.Now()

	newBlock.Index = prev.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = prev.Hash
	newBlock.Hash = newBlock.calculateHash()

	return newBlock

}
*/

// create a new block using previous block's hash
//implement proof of work
//taken from https://medium.com/@mycoralhealth/code-your-own-blockchain-mining-algorithm-in-go-82c6a71aba1f
//ill expand on this

func generateBlock(oldBlock Block, BPM int, difficulty int) Block {
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PreviousHash = oldBlock.Hash
	newBlock.Difficulty = difficulty

	for i := 0; ; i++ {
<<<<<<< HEAD
		hex := fmt.Sprintf("%d", i)
		newBlock.Nonce = hex
		if !isHashValid(calculateHash(newBlock), newBlock.Difficulty) {
			fmt.Println(calculateHash(newBlock), " do more work!")
			//might be unneccessary
=======
		hexVal := fmt.Sprintf("%x", i)
		newBlock.Nonce = hexVal

		hash := newBlock.calculateHash()
		if !isHashValid(hash, newBlock.Difficulty) {
			fmt.Println(hash, " do more work!")
>>>>>>> 1bb2d1d76577d2f46f7df899a0dee7e668d7819f
			time.Sleep(time.Second)
			continue
		} else {
			fmt.Println(hash, " work done!")
			newBlock.Hash = hash
			break
		}

	}
	return newBlock
}

//check if hash has correct number of zeros
func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

<<<<<<< HEAD
func makeNewPeer(listenPort int, secio bool, randseed int64) (host.Host, error) {
	var rdr io.Reader

	//not sure if we need, this check when we start running things
	if randseed == 0 {
		rdr = rand.Reader
	} else {
		rdr = mrand.New(mrand.NewSource(randseed))
	}

	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 1024, rdr)
	if err != nil {
		return nil, err
	}

	//get host names
	var name1 string = "vdi-linux-030.ccs."
	var name2 string = "vdi-linux-031.ccs."
	addr1, err := net.LookupHost(name1)
	addr2, err := net.LookupHost(name1)

	//need to figure out how to set up hosts

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
		libp2p.Identity(priv),
	}

	basicHost, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addr := basicHost.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr)
	log.Printf("I am %s\n", fullAddr)
	if secio {
		log.Printf("Now run \"go run main.go -l %d -d %s -secio\" on a different terminal\n", listenPort+1, fullAddr)
	} else {
		log.Printf("Now run \"go run main.go -l %d -d %s\" on a different terminal\n", listenPort+1, fullAddr)
	}

	return basicHost, nil

}

func handleStream(s net.Stream) {

	log.Println("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go readData(rw)
	go writeData(rw)

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
		bpm, err := strconv.Atoi(sendData)
		if err != nil {
			log.Fatal(err)
		}
		newBlock := generateBlock(Blockchain[len(Blockchain)-1], bpm)

		//how does validate arguments work
		if validate(newBlock, Blockchain[len(Blockchain)-1]) {
			mutex.Lock()
			Blockchain = append(Blockchain, newBlock)
			mutex.Unlock()
		}

		bytes, err := json.Marshal(Blockchain)
		if err != nil {
			log.Println(err)
		}

		//spew.Dump(Blockchain)

		mutex.Lock()
		rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
		rw.Flush()
		mutex.Unlock()
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

=======
>>>>>>> 1bb2d1d76577d2f46f7df899a0dee7e668d7819f
func proofOfWork() {

}

func main() {

	//get hostnames from hosts.txt
	//we use vdi-030, 031, 032 for now
	var host_names [10]string
	var host_count int = 0

	file, err := os.Open("hosts.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		host_names[host_count] = scanner.Text()
		host_count++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	//so the BPMs is simply the data of the block
	currtime := time.Now()
	genesisBlock := Block{0, currtime.String(), 0, calculateHash(genesisBlock), ""}

	Blockchain = append(Blockchain, genesisBlock)

	//who needs logging

	// We dont most of these command line arguemtns
	listenF := flag.Int("l", 0, "wait for incoming connections")
	target := flag.String("d", "", "target peer to dial")
	secio := flag.Bool("secio", false, "enable secio")
	seed := flag.Int64("seed", 0, "set random seed for id generation")
	flag.Parse()

	if *listenF == 0 {
		log.Fatal("Please provide a port to bind on with -l")
	}

	// Make a host that listens on the given multiaddress
	ha, err := makeNewPeer(*listenF, *secio, *seed)
	if err != nil {
		log.Fatal(err)
	}

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

		pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
		if err != nil {
			log.Fatalln(err)
		}

		peerid, err := peer.IDB58Decode(pid)
		if err != nil {
			log.Fatalln(err)
		}

		// Decapsulate the /ipfs/<peerID> part from the target
		// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
		targetPeerAddr, _ := ma.NewMultiaddr(
			fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
		targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

		// We have a peer ID and a targetAddr so we add it to the peerstore
		// so LibP2P knows how to contact it
		ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

		log.Println("opening stream")
		// make a new stream from host B to host A
		// it should be handled on host A by the handler we set above because
		// we use the same /p2p/1.0.0 protocol
		s, err := ha.NewStream(context.Background(), peerid, "/p2p/1.0.0")
		if err != nil {
			log.Fatalln(err)
		}
		// Create a buffered stream so that read and writes are non blocking.
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		// Create a thread to read and write data.
		go writeData(rw)
		go readData(rw)

		select {} // hang forever

	}

}
