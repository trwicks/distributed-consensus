package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/graph/{nodeCount:[0-9]+}", a.getGraph).Methods("GET")
}

type App struct {
	Router *mux.Router
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) getGraph(w http.ResponseWriter, r *http.Request) {

	// need to name variables/files betterer

	params := mux.Vars(r)
	nodeCount, err := strconv.Atoi(params["nodeCount"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid NodeCount")
		return
	}
	fmt.Println(nodeCount)
	nodes = new(graph.Node)
	// graph.CreateRandomGraph(nodeCount)

	// graph.BroadcastNodeInfo()
	// // var results grap
	// results := graph.ConsensusResult()
	respondWithJSON(w, http.StatusOK, 2)
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8000", a.Router))
}

func main() {
	a := App{}
	a.Run(":8000")
}
