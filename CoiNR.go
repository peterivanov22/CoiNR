package CoiNR

import (
	"fmt"
	"sync"
)

// This is the main class

// Not sure this is where this should go. Storing the global difficulty of the blocks.
var difficulty int

// A block chain is just a slice of blocks
var Blockchain []Block

// seems like we need a mutex
var mutex = &sync.Mutex{}

func main() {

	// presumably kick stuff off

	fmt.Println("Hoooray!  You ran CoiNR!  Full functionality not supported.  Shutting down.")

}
