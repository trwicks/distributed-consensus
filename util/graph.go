package graph

import (
	"fmt"
	"math/rand"
	"sync"
)

type Node struct {
	id        uint64
	largest   uint64
	nodeEdges []*Node
	messages  chan uint64
	mu        sync.Mutex
}

type Leader struct {
	id    uint64
	count uint64
}

type Graph struct {
	// a basic graph type that contains a list of references to associated nodes.
	nodes []*Node
	done  chan bool
}

func (graph *Graph) CreateRandomGraph(nodeNumber int) {

	r := rand.New(rand.NewSource(22))

	// construct graph - edges unassigned
	for i := 0; i < nodeNumber; i++ {
		randNumber := r.Uint64()
		graph.nodes = append(graph.nodes, &Node{
			id:        randNumber,
			largest:   randNumber,
			nodeEdges: nil,
			messages:  make(chan uint64, nodeNumber/10),
		})
	}

	fmt.Println("Graph size:", len(graph.nodes))

	for i := 0; i < len(graph.nodes); i++ {
		fmt.Println("Node ID:", graph.nodes[i].id)

		// assign edges for each edge until number of edges = 1/10 * number of nodes connected
		for len(graph.nodes[i].nodeEdges) < calcEdgeNumber(nodeNumber) {
			rN := r.Uint64() % uint64(nodeNumber)
			if nodeInSet(rN, graph.nodes[i].getEdgeIds()) == false {
				addEdge(graph.nodes[i], graph.nodes[rN])
			}
		}

		fmt.Println("Node edges:", len(graph.nodes[i].nodeEdges))
		fmt.Println("Highest ID:", graph.nodes[i].largest)
		fmt.Println("================================")
	}
}

func calcEdgeNumber(n int) int {
	// edge case: for smaller graphs < 20 nodes number of edges needs to be higher to ensure
	// that all nodes
	if n < 20 {
		return n / 4
	}
	return n / 10
}

func nodeInSet(x uint64, nodes []uint64) bool {
	for _, n := range nodes {
		if x == n {
			return true
		}
	}
	return false
}

func addEdge(n1, n2 *Node) {
	n1.nodeEdges = append(n1.nodeEdges, n2)
	n2.nodeEdges = append(n2.nodeEdges, n1)
}

func (n *Node) getEdgeIds() []uint64 {
	var edgeIds []uint64
	for _, edgeId := range n.nodeEdges {
		edgeIds = append(edgeIds, edgeId.id)
	}
	return edgeIds
}

var wg sync.WaitGroup

func (graph *Graph) BroadcastNodeInfo() {

	// var graphWg sync.WaitGroup
	for _, node := range graph.nodes {
		// graphWg.Add(1)

		go node.messageEdges()
	}
	// graphWg.Wait()
	wg.Wait()

}

// receive message from queue - adjust largest if necessary
// TODO: to prevent race condition lock value while it is being pulled from the queue
// does this go routine exit prematurely and messages may be pushed on to the channel after it has exited -
// large diameter graphs

func (n *Node) messageEdges() {
	wg.Add(1)
	go n.announceLargest()

	for message := range n.messages {
		if message > n.largest {
			n.mu.Lock()
			n.largest = message // CONSIDER USING MUTEX FOR ALTERING THIS TO ENSURE ATOMICITY
			n.mu.Unlock()
			wg.Add(1)
			go n.announceLargest()
		}
	}
}

// might be closing to early if other nodes are transfer message
func (n *Node) announceLargest() {
	defer wg.Done()
	for _, e := range n.nodeEdges {
		e.messages <- n.largest
	}
}

func (graph *Graph) ConsensusResult() {
	nodeCounts := make(map[uint64]int)
	for _, node := range graph.nodes {
		fmt.Printf("Node Id: %d \t - Largest Node: %d\n", node.id, node.largest)
		// make a map that contains count of largests from all nodes
		// print map for final statistics on consensus
		nodeCounts[node.largest] += 1
	}
	fmt.Printf("Node Consensus Results:\n")
	for k, v := range nodeCounts {
		fmt.Printf("Node ID: %d \t Count: %d\n", k, v)
	}
}
