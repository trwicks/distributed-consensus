// package main

// import (
// 	// "github.com/wickst/distributed-consensus/util"
// 	"fmt"
// 	"log"
//     "net/http"
// )

// func handler(w http.ResponseWriter, r *http.Request) {
//     fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
// }

// func main() {
//     http.HandleFunc("/", handler)
//     log.Fatal(http.ListenAndServe(":8080", nil))
// }
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

	fmt.Println(string(resultJson))
}

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/api/{nodeCount}", GetGraph).Methods("GET")
    log.Fatal(http.ListenAndServe(":8000", router))
}


// func main() {
// 	// see if there is a commandline package to make named command line arguments
// 	if len(os.Args) != 2 {
// 		fmt.Println("Incorrect number of arguments given. Please use a valid argument for the number of nodes.")
// 		os.Exit(1)
// 	}
// 	args := os.Args[1:]

// 	nodeNumber, err := strconv.Atoi(args[0])
// 	// TODO: check type of argument passed in via commandline
// 	if err != nil {
// 		fmt.Println("Incorrect arguments given. Please use the following command ./")
// 		os.Exit(1)
// 	}

// }
