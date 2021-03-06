CS 7610 Final Project README
Miles Benjamin & Peter Ivanov

How to make prj3:
    Run the make file included with the project.  Golang must be installed on the machine.
    Alternatively run "go build CoiNR.go Block.go network.go Taction.go"

How to run prj3:
    Run from the command line.  You can run multiple nodes on the same machine in seperate terminal windows.
	This is highly recommended for ease of testing.

	The first console can simply be run using the following line:
	./CoiNR -l <yourPort>

	This first console will give you the information you need to run additional consoles.  Additional peers will be run using:
	./CoiNR -l <yourPort> -d /ip4/<PUBLIC IP>/tcp/<PortToConnect>/ipfs/<PEER ID>

	You must use different ports for each client if they are running on the same machine.

How to make a block:

	The simplest way to make a new block is simply to type 'm".  This creates a transactionless block both as a demo and to seed
	the user with a coin for future transactions.

	You can type b to see the current balance of the user.

	Once multiple peers are running simply enter transactions in the following pattern:

	RecipientAddress Amount

	A transaction will be created, sending the amount from the host to the recipient.
	Not extensively tested, recommended to just use 1 for Amt.

	You may enter as many as you like, once the peer has enough to get started it will start forming a block. Any additional transactions
	will be stored and put into a future block. 

	The address of each node will be printed when first launched.


Command Line Flags:

    -l XXX - * REQUIRED * The port that the process should listen at.
	-d XXX - An address string provided by the first console. This will let you connect to the network. Required on all further consoles.
		The string should look like:
		/ip4/127.0.0.1/tcp/10001/ipfs/QmZd7YGcGUUzffMychA7YornZxxN9kJzKujsppsDA2d3yr
    -v true - Verbose logging to the console. Can be enabled on multiple machines.

   
    