package graph

import (
	"fmt"
	"math/rand"
	"sync"
)

type Node struct {
	id        uint64
	largest   *Leader
	nodeEdges []*Node
	messages  chan *Leader
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
			id: randNumber,
			largest: &Leader{
				id:    randNumber,
				count: 1,
			},
			nodeEdges: nil,
			messages:  make(chan *Leader, nodeNumber/10),
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
		fmt.Println("Highest ID:", graph.nodes[i].largest.id)
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
	for _, node := range graph.nodes {
		go node.messageEdges()
	}
	// wait before exiting to ensure all go-routines terminate cleanly
	wg.Wait()
}

// receive message from queue - adjust largest if necessary
// TODO: to prevent race condition lock value while it is being pulled from the queue
func (n *Node) messageEdges() {
	wg.Add(1)

	go n.announceLargest()

	for message := range n.messages {

		if message.id > n.largest.id {

			n.largest.id = message.id
			n.largest.count = message.count + 1

			wg.Add(1)
			go n.announceLargest()
		}
	}
}

func (n *Node) announceLargest() {
	defer wg.Done()
	for _, e := range n.nodeEdges {
		// DATA RACE:
		// MUTEX USED TO ENSURE THAT ONLY 1 PROCESSES CAN WRITE TO A NODE'S CHANNEL
		n.mu.Lock()
		e.messages <- n.largest
		n.mu.Unlock()
	}
}

func getLargestNode(nodes []*Node) uint64 {
	largest := nodes[0].id
	for _, n := range nodes {
		if n.id > largest {
			largest = n.id
		}
	}
	return largest
}

func (graph *Graph) ConsensusResult() {
	nodeCounts := make(map[uint64]int)
	largestNode := getLargestNode(graph.nodes)
	for _, node := range graph.nodes {
		fmt.Printf("Node Id: %d \t - Largest Node: %d\n", node.id, node.largest.id)
		// make a map that contains count of largests from all nodes
		// print map for final statistics on consensus
		nodeCounts[node.largest.id] += 1
	}
	fmt.Printf("Largest node: %d\n", largestNode)
	fmt.Printf("Node Consensus Results:\n")
	for k, v := range nodeCounts {
		fmt.Printf("Node ID: %d \t Count: %d\n", k, v)
	}
}
