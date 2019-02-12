package network

import (
	"errors"
	"fmt"
	"log"
	"net/rpc"

	"github.com/jessicagreben/kademlia/pkg/node"
	"github.com/jessicagreben/kademlia/pkg/types"
)

// Network is a x.
type Network struct {
	rt *routingTable
}

// Join creates a node ID, creates a routing table, and populates the routing table for a node.
func (n *Network) Join(currIP string, currPort string) error {

	// Generate an ID for the current node.
	// TODO: If the ID has been created before, use the previous ID instead of createing a new ID.
	id := node.GenerateID(types.IDLength)

	self := types.Contact{
		NodeID: id,
		IP:     currIP,
		Port:   currPort,
	}

	// Create a routing table.
	n.rt = newRoutingTable(self)
	fmt.Println("route table", n.rt)

	if currIP == "boot" {
		return nil
	}

	// Populate the routing table by performing iterative queries to find nodes in the network.
	// Start by adding self to the bootstrap node routing table. Do this by performing a lookup on self.
	listContacts, err := lookup(self.NodeID, n.rt.boot, n.rt.currentNode)
	if err != nil {
		return err
	}
	fmt.Println(listContacts)

	if len(listContacts) == 0 {
		return nil
	}

	// Next, add self to the closest contacts returned by the bootstrap node.
	_, err = recursiveLookup(self.NodeID, listContacts, n.rt)
	if err != nil {
		return err
	}
	return nil
}

func lookup(desiredNodeID types.NodeID, otherNode types.Contact, currentNode types.Contact) ([]types.Contact, error) {
	addr := fmt.Sprintf("%s:%s", otherNode.IP, otherNode.Port)
	client := client(addr)
	l := ListContacts{}
	args := LookupArgs{
		RequestFrom:   currentNode,
		DesiredNodeID: desiredNodeID,
	}
	if err := client.Call("Network.Lookup", args, &l); err != nil {
		return []types.Contact{}, err
	}
	if l.Success {
		fmt.Println("success")
		return l.Contacts, nil
	}
	fmt.Println("error")
	return []types.Contact{}, errors.New(l.ErrMsg)
}

// ListContacts is the response to the Lookup RPC.
type ListContacts struct {
	Success  bool
	Found    bool
	Contacts []types.Contact
	ErrMsg   string
}

// Lookup checks the route table for the node ID and returns a contact if it exists.
// Otherwise, it returns the closest nodes to that ID.
func (n *Network) Lookup(a LookupArgs, reply *ListContacts) error {
	desiredNodeID := a.DesiredNodeID
	requestFrom := a.RequestFrom

	// Look for desired node ID in the route table.
	c, found := n.rt.find(desiredNodeID)
	if found {
		reply.Contacts = []types.Contact{c}
		reply.Success = true
		reply.Found = true

		// Update Contact of the node making the request to the route table.
		n.rt.add(requestFrom)
		return nil
	}

	// If the desiredNodeID is not found in the routing table
	// then find the closest nodes and return those.
	xorBytes := node.Distance(desiredNodeID, n.rt.currentNode.NodeID)
	ind := node.FindLongestPrefix(xorBytes)
	closestNodes := n.rt.findClosestNodes(ind)
	reply.Contacts = closestNodes
	reply.Success = true
	reply.Found = false

	// Update Contact of the node making the request to the route table.
	n.rt.add(requestFrom)
	return nil
}

func client(addr string) *rpc.Client {
	client, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		log.Fatal("Connection error: ", err)
	}
	return client
}

// Pong is the response to the Ping RPC.
type Pong struct {
	Success bool
	ErrMsg  string
}

// Args is x.
type Args struct{}

// Ping is a method to see if a contact is still available.
func Ping(c types.Contact) (bool, error) {
	addr := fmt.Sprintf("%s:%s", c.IP, c.Port)
	client := client(addr)
	p := Pong{}
	if err := client.Call("Network.Pong", Args{}, &p); err != nil {
		fmt.Println("error")
		return false, err
	}
	if p.Success {
		fmt.Println("success")
		return true, nil
	}
	return false, errors.New(p.ErrMsg)
}

// Pong is a x.
func (n *Network) Pong(a Args, reply *Pong) error {
	fmt.Println("\nin pong")
	reply.Success = true
	return nil
}

// LookupArgs are the arguments to the Lookup RPC.
type LookupArgs struct {
	RequestFrom   types.Contact
	DesiredNodeID types.NodeID
}

func recursiveLookup(desiredNodeID types.NodeID, listContacts []types.Contact, rt *routingTable) (types.Contact, error) {
	switch len(listContacts) {
	case 0:
		return types.Contact{}, nil
	case 1:
		currContact := listContacts[0]

		// If the current Contact is the nodeID we're looking for then we are done.
		if currContact.NodeID == desiredNodeID {
			return currContact, nil
		}

		// If the current Contact is not the nodeID we're looking for, then
		// continue with the recursive search.
		rt.add(currContact)
		lcs, err := lookup(desiredNodeID, currContact, rt.currentNode)
		if err != nil {
			return types.Contact{}, err
		}
		recursiveLookup(desiredNodeID, lcs, rt)
	default:
		for _, currContact := range listContacts {
			rt.add(currContact)
			lcs, err := lookup(desiredNodeID, currContact, rt.currentNode)
			if err != nil {
				return types.Contact{}, err
			}
			recursiveLookup(desiredNodeID, lcs, rt)
		}
	}

	return types.Contact{}, nil
}
