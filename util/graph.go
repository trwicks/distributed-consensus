package graph

import (
	"fmt"
	"math/rand"
	"sync"
)

type Node struct {
	ID        uint64  `json:"ID"`
	Largest   *Leader `json:"Largest"`
	NodeEdges []*Node `json:"NodeEdges"`
	messages  chan *Leader
	mu        sync.Mutex
}

type Leader struct {
	ID    uint64 `json:"ID"`
	Count uint64 `json:"Count"`
}
type Graph struct {
	// a basic graph type that contains a list of references to associated nodes.
	Nodes []*Node `json:"Nodes"`
	done  chan bool
}

type Results struct {
	Nodes map[uint64][]uint64
	Stats map[uint64]int
}

func (graph *Graph) CreateRandomGraph(nodeNumber int) {

	r := rand.New(rand.NewSource(22))

	// construct graph - edges unassigned
	for i := 0; i < nodeNumber; i++ {
		randNumber := r.Uint64()
		graph.Nodes = append(graph.Nodes, &Node{
			ID: randNumber,
			Largest: &Leader{
				ID:    randNumber,
				Count: 1,
			},
			NodeEdges: nil,
			messages:  make(chan *Leader, nodeNumber/10),
		})
	}

	fmt.Println("Graph size:", len(graph.Nodes))

	for i := 0; i < len(graph.Nodes); i++ {
		fmt.Println("Node ID:", graph.Nodes[i].ID)

		// assign edges for each edge until number of edges = 1/10 * number of nodes connected
		for len(graph.Nodes[i].NodeEdges) < calcEdgeNumber(nodeNumber) {
			rN := r.Uint64() % uint64(nodeNumber)
			if nodeInSet(rN, graph.Nodes[i].getEdgeIDs()) == false {
				addEdge(graph.Nodes[i], graph.Nodes[rN])
			}
		}

		// logging information
		fmt.Println("Node edges:", len(graph.Nodes[i].NodeEdges))
		// fmt.Println("Highest ID:", graph.Nodes[i].largest.ID)
		fmt.Println("================================")
	}

}

func calcEdgeNumber(n int) int {
	// edge case: for smaller graphs < 20 nodes number of edges needs to be higher to ensure
	// that all nodes have enough edges to ensure graph is fully connected.
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
	n1.NodeEdges = append(n1.NodeEdges, n2)
	n2.NodeEdges = append(n2.NodeEdges, n1)
}

func (n *Node) getEdgeIDs() []uint64 {
	var edgeIDs []uint64
	for _, edgeID := range n.NodeEdges {
		edgeIDs = append(edgeIDs, edgeID.ID)
	}
	return edgeIDs
}

var wg sync.WaitGroup

func (graph *Graph) BroadcastNodeInfo() {
	for _, node := range graph.Nodes {
		wg.Add(1)
		go node.messageEdges()
	}
	// wait before exiting to ensure all go-routines terminate cleanly
	wg.Wait()
	fmt.Println("Doing pushup on 100s")
}

// Send out messages to associated nodes. Receive messages from queue - adjust largest known if necessary and resend if updated.
func (n *Node) messageEdges() {
	go n.announceLargest()
	for message := range n.messages {
		if message.ID > n.Largest.ID {

			n.Largest.ID = message.ID
			n.Largest.Count = message.Count + 1

			// Messaging process can continue but the graph is highly connected
			wg.Add(1)
			go n.announceLargest()
		}
	}
	close(n.messages)
}

// TODO: each node that needs to pass messages to other nodes will pass a directed channel via argument in function call
// receiving node will read messages from the channel then close the channel
// once all messages are received then return highest value.
func (n *Node) receiveMessage(messages <-chan *Leader) {
	for m := range messages {
		if m.ID > n.Largest.ID {
			n.messages <- m
		}
	}
}

// make relay message function

func (n *Node) announceLargest() {
	defer wg.Done()
	for _, e := range n.NodeEdges {
		// DATA RACE FOR DISCUSSION AND IMPROVEMENT
		// MUTEX USED TO ENSURE THAT ONLY 1 PROCESSES CAN WRITE TO A NODE'S CHANNEL
		e.messages <- n.Largest
	}
}

func getLargestNode(nodes []*Node) uint64 {
	largest := nodes[0].ID
	for _, n := range nodes {
		if n.ID > largest {
			largest = n.ID
		}
	}
	return largest
}

func (graph *Graph) ConsensusResult() []Node {
	// largestNode := getLargestNode(graph.Nodes)
	var Nodes []Node
	stats := make(map[uint64]int)
	for _, node := range graph.Nodes {
		fmt.Printf("Node ID: %d \t - Largest Known Node: %d with count %d\n", node.ID, node.Largest.ID, node.Largest.Count)
		stats[node.Largest.ID] += 1
		Nodes = append(Nodes, *node)
	}
	// fmt.Printf("Largest node: %d\n", largestNode)
	fmt.Printf("Node Consensus Results:\n")
	for k, v := range stats {
		percentage := v / len(graph.Nodes) * 100
		fmt.Printf("Node ID: %d \t Count: %d. \t%d percent graph consensus.\n", k, v, percentage)
	}
	return Nodes
}
