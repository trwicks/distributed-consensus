## Distributed Graph Consensus 

### Project Description
This project aimed to perform the following functions to illustrate consensus of the largest number within the network (assigned as a node ID). This is work in progress and is achieved so far using the following processes:
1. Create a random graph where nodes are connected to one tenth of all other nodes within the graph.
2. Nodes communicate to other nodes via message queues.
3. Each node communicates within its own thread.
4. Nodes broadcast to other nodes, the highest known node ID the node has received.
5. Nodes compare the highest value from associated nodes (edges) with the highest value the node has received previously.
6. Once all broadcasts have ceased (to be completed) nodes vote on the highest value they have received via associated nodes. A consensus is then reached on the highest value node ID.

### Assumptions
- Nodes do not need to know the whole or partial graph structure, only their associated edges.
- Nodes only pass on highest node Id if it is higher than their current highest Id. An edge case might exist where two or more nodes have the same highest value (although very unlikely). The end result, however, will still likely be a consensus between all nodes of the correct highest value, shared between multiple nodes. 
- Nodes pass announce highest ID to all connected nodes initially, and once they have recieved a higher value than their current highest value. This means the sender of the highest value to node the receiver will receive a subsequent message from the received indicating that its highest known value has been updated. While efficient this does not matter as the sender of the value will compare it to it's current value and will not pass on the value to its edge nodes as the value is only equal, not higher than its current known value.  

## Use

To run the code once dependencies are installed simply run the following command:
```$bash
$ make build
$ ./distributed-consensus {number of nodes} // e.g. 67
```

## Dependencies
- go v1.10

## TODO
- voting mechanism after all broadcasts have been announced. 