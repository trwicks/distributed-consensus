package main

import (
	"distributed-consensus/util"
	"fmt"
	"os"
	"strconv"
)

func main() {
	// see if there is a commandline package to make named command line arguments
	args := os.Args[1:]

	nodeNumber, err := strconv.Atoi(args[0])
	// TODO: check type of argument passed in via commandline
	if err != nil {
		fmt.Println("error exiting")
	}

	// need to name variables/files betterer
	var graph graph.Graph

	graph.CreateRandomGraph(nodeNumber)

	graph.BroadcastNodeInfo()

	// CHeating - requires channel to indicate when all messages have finished sending
	// var input string
	// fmt.Scanln(&input)
	graph.ConsensusResult()
}
