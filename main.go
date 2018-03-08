package main

import (
	"encoding/json"
	"github.com/wickst/distributed-consensus/util"
	"strconv"
	"fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

func GetGraph(w http.ResponseWriter, r *http.Request) {

	// need to name variables/files betterer

	params := mux.Vars(r)
	fmt.Println(params)
	var graph graph.Graph
	nodeCount, err := strconv.Atoi(params["nodeCount"])
	if err != nil {
		fmt.Println("Write code to handle this")
	}
	graph.CreateRandomGraph(nodeCount)

	graph.BroadcastNodeInfo()
	
	results := graph.ConsensusResult()

	resultJson, err := json.Marshal(results)
	if err != nil {
        fmt.Printf("Error: %s", err)
        return;
	}
	for _, nodes := range results.Nodes {
		fmt.Println(nodes.Id)
	}

	fmt.Println(string(resultJson))
}

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/api/{nodeCount}", GetGraph).Methods("GET")
    log.Fatal(http.ListenAndServe(":8000", router))
}