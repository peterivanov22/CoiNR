CS 7610 Final Project README
Miles Benjamin & Peter Ivanov

How to make prj3:
    Run the make file included with the project.  Golang must be installed on the machine.
    Alternatively run "go build CoiNR.go Block.go network.go Taction.go"

How to run prj3:
    Run from the command line.  You can run multiple nodes on the same machine in seperate terminal windows.
	This is highly recommended for ease of testing.

	The first console can simply be run using the following line:
	./CoiNR -l <yourPort> -p <yourName>

	This first console will give you the information you need to run additional consoles.  Additional peers will be run using:
	./CoiNR -l <yourPort> -p <yourName> -d /ip4/<PUBLIC IP>/tcp/<PortToConnect>/ipfs/<PEER ID>

	You must use different ports for each client if they are running on the same machine.

	Once multiple peers are running simply enter transactions in the following pattern:

	User1 User2 Amt

	You may enter as many as you like, once the peer has enough to get started it will start forming a block. Any additional transactions
	will be stored and put into a future block. 

Command Line Flags:

    -l XXX - * REQUIRED * The port that the process should listen at.
    -p XXX - * REQUIRED * the public key of the user at this node.  Any new coinbase transactions will go to this user.
	-d XXX - An address string provided by the first console. This will let you connect to the network. Required on all further consoles. 
		The string should look like:
		/ip4/127.0.0.1/tcp/10001/ipfs/QmZd7YGcGUUzffMychA7YornZxxN9kJzKujsppsDA2d3yr
    -v true - Verbose logging to the console. Can be enabled on multiple machines.

Known issues:
   
    