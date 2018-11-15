package CoiNR

import "sync"


// This is the main class


// Not sure this is where this should go. Storing the global difficulty of the blocks.
var difficulty int


// A block chain is just a slice of blocks
var Blockchain []Block

// seems like we need a mutex
var mutex = &sync.Mutex{}



