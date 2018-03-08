package graph

import (
	"fmt"
	"math/rand"
	"sync"
)

type Node struct {
	Id        uint64		`json:"Id"`
	largest   *Leader		`json:"Largest"`
	nodeEdges []*Node		`json:"NodeEdges"`
	messages  chan *Leader
	mu        sync.Mutex
}

type Leader struct {
	Id    uint64
	count uint64
}
type Graph struct {
	// a basic graph type that contains a list of references to associated nodes.
	nodes []*Node 			
	done  chan bool 		
	
}

type Results struct {
	Nodes []*Node 			`json:"Nodes"`
	Stats map[uint64]int	`json:"Stats"`
}


func (graph *Graph) CreateRandomGraph(nodeNumber int) {

	r := rand.New(rand.NewSource(22))

	// construct graph - edges unassigned
	for i := 0; i < nodeNumber; i++ {
		randNumber := r.Uint64()
		graph.nodes = append(graph.nodes, &Node{
			Id: randNumber,
			largest: &Leader{
				Id:    randNumber,
				count: 1,
			},
			nodeEdges: nil,
			messages:  make(chan *Leader, nodeNumber/10),
		})
	}

	fmt.Println("Graph size:", len(graph.nodes))

	for i := 0; i < len(graph.nodes); i++ {
		fmt.Println("Node Id:", graph.nodes[i].Id)

		// assign edges for each edge until number of edges = 1/10 * number of nodes connected
		for len(graph.nodes[i].nodeEdges) < calcEdgeNumber(nodeNumber) {
			rN := r.Uint64() % uint64(nodeNumber)
			if nodeInSet(rN, graph.nodes[i].getEdgeIds()) == false {
				addEdge(graph.nodes[i], graph.nodes[rN])
			}
		}

		// logging information
		fmt.Println("Node edges:", len(graph.nodes[i].nodeEdges))
		// fmt.Println("Highest Id:", graph.nodes[i].largest.Id)
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
	n1.nodeEdges = append(n1.nodeEdges, n2)
	n2.nodeEdges = append(n2.nodeEdges, n1)
}

func (n *Node) getEdgeIds() []uint64 {
	var edgeIds []uint64
	for _, edgeId := range n.nodeEdges {
		edgeIds = append(edgeIds, edgeId.Id)
	}
	return edgeIds
}

var wg sync.WaitGroup

func (graph *Graph) BroadcastNodeInfo() {
	for _, node := range graph.nodes {
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
		if message.Id > n.largest.Id {

			n.largest.Id = message.Id
			n.largest.count = message.count + 1

			// Messaging process can continue but the graph is highly connected
			wg.Add(1)
			go n.announceLargest()
		}
	}
}

// TODO: each node that needs to pass messages to other nodes will pass a directed channel via argument in function call
// receiving node will read messages from the channel then close the channel
// once all messages are received then return highest value.
func (n *Node) receiveMessage(messages <-chan *Leader) {
	for m := range messages {
		if m.Id > n.largest.Id {
			n.messages <- m
		}
	}
}

// make relay message function

func (n *Node) announceLargest() {
	defer wg.Done()
	for _, e := range n.nodeEdges {
		// DATA RACE FOR DISCUSSION AND IMPROVEMENT
		// MUTEX USED TO ENSURE THAT ONLY 1 PROCESSES CAN WRITE TO A NODE'S CHANNEL
		e.messages <- n.largest
	}
}

func getLargestNode(nodes []*Node) uint64 {
	largest := nodes[0].Id
	for _, n := range nodes {
		if n.Id > largest {
			largest = n.Id
		}
	}
	return largest
}

func (graph *Graph) ConsensusResult() Results {
	largestNode := getLargestNode(graph.nodes)

	results := Results{
		Nodes: graph.nodes,
		Stats: make(map[uint64]int),
	}
	// for _, node := range graph.nodes {
	// 	results.Nodes = append(results.Nodes, &Node{
	// 		Id: node.Id,
	// 		largest: node.largest,
	// 		nodeEdges: node.nodeEdges,
	// 		messages:  nil,
	// 	})
	// }
	// results.Nodes = graph.nodes
	// for _, node := range results.Nodes {
	// 	fmt.Println(node)
	// }

	for _, node := range graph.nodes {
		// logging info:
		fmt.Printf("Node Id: %d \t - Largest Known Node: %d with count %d\n", node.Id, node.largest.Id, node.largest.count)
		// make a map that contains count of largests from all nodes
		// print map for final statistics on consensus
		results.Stats[node.largest.Id] += 1
	}
	fmt.Printf("Largest node: %d\n", largestNode)
	fmt.Printf("Node Consensus Results:\n")
	for k, v := range results.Stats {
		percentage := v / len(graph.nodes) * 100
	
		fmt.Printf("Node Id: %d \t Count: %d. \t%d percent graph consensus.\n", k, v, percentage)
	}
	fmt.Println(results.Stats)
	return results
}
