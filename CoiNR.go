package CoiNR

import "sync"


// A block chain is just a slice of blocks
var Blockchain []Block

// seems like we need a mutex
var mutex = &sync.Mutex{}



