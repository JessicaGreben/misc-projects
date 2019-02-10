package main

import (
	"errors"
	"fmt"
	"log"
	"net/rpc"
)

// Network is a x.
type Network struct {
	rt *routingTable
}

func (n *Network) join(currIP string, currPort string) error {

	// Generate an ID for the current node.
	// TODO: If the ID has been created before, use the previous ID instead of createing a new ID.
	id := generateNodeID(idLength)

	self := Contact{
		Node:   id,
		IpAddr: currIP,
		Port:   currPort,
	}

	// Create a routing table.
	n.rt = newRoutingTable(self)
	fmt.Println("route table", n.rt)

	if currPort == "8080" {
		return nil
	}

	// Populate the routing table by performing iterative queries to find nodes in the network.
	// Start by adding self to the bootstrap node routing table. Do this by performing a lookup on self.
	listContacts, err := lookup(self.Node, n.rt.boot)
	if err != nil {
		return err
	}
	fmt.Println(listContacts)

	// TODO: if len(listContacts) == 0

	// Next, add self to the closest contacts returned by the bootstrap node.
	// _, err = recursiveLookup(self.node, listContacts, n.rt)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func lookup(desiredNodeID nodeID, otherNode Contact) ([]Contact, error) {
	addr := fmt.Sprintf("%s:%s", otherNode.IpAddr, otherNode.Port)
	fmt.Println("rpc Lookup to server:", addr)
	client := client(addr)
	l := ListContacts{}
	args := LookupArgs{
		DesiredNodeID: desiredNodeID,
	}
	if err := client.Call("Network.Lookup", args, &l); err != nil {
		return []Contact{}, err
	}
	if l.Success {
		fmt.Println(l)
		fmt.Println("success")
		return l.Contacts, nil
	}
	fmt.Println("error")
	return []Contact{}, errors.New(l.ErrMsg)
}

// ListContacts is the response to the Lookup RPC.
type ListContacts struct {
	Success  bool
	Found    bool
	Contacts []Contact
	ErrMsg   string
}

// Lookup is x.
func (n *Network) Lookup(a LookupArgs, reply *ListContacts) error {
	desiredNodeID := a.DesiredNodeID

	// Look in the contact's rt for the current nodeID.
	_, found := n.rt.currentNodesInBucket[desiredNodeID]
	if found {

		// bucket := getBucket(desiredNodeID, n.rt)
		// 	ind, ok := bucket.find(desiredNodeID)
		// 	if !ok {
		// 		// do something
		// 	}
		// 	c := bucket.cut(ind, n.rt)
		// 	reply.Contacts = []Contact{c}
		reply.Success = true
		reply.Found = true
	}

	// If the desiredNodeID is not found in the routing table
	// then find the closest nodes and return those.
	xorBytes := nodeDistance(desiredNodeID, n.rt.currentNode.Node)
	fmt.Println("xorBytes:", xorBytes)
	ind := findLongestPrefix(xorBytes)
	fmt.Println("ind:", ind)
	closestNodes := findClosestNodes(ind, n.rt)
	fmt.Println("closestNodes:", closestNodes)
	reply.Contacts = closestNodes
	reply.Success = true
	reply.Found = false
	return nil
}

func findClosestNodes(xorPrefix int, rt *routingTable) []Contact {
	for i := xorPrefix; i < 0; i-- {
		currBucket := rt.buckets[i]
		if len(currBucket) > 0 {
			return currBucket
		}
	}
	return []Contact{}
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
func Ping() error {
	addr := "localhost:8080"
	client := client(addr)
	p := Pong{}
	if err := client.Call("Network.Pong", Args{}, &p); err != nil {
		return err
	}
	if p.Success {
		fmt.Println("success")
		return nil
	}
	fmt.Println("error")
	return errors.New(p.ErrMsg)
}

// Pong is a x.
func (n *Network) Pong(a Args, reply *Pong) error {
	fmt.Println("\nin pong")
	reply.Success = true
	return nil
}

// LookupArgs is x.
type LookupArgs struct {
	DesiredNodeID nodeID
}

func getBucket(desiredNodeID nodeID, rt *routingTable) bucket {
	xorBytes := nodeDistance(desiredNodeID, rt.currentNode.Node)
	ind := findLongestPrefix(xorBytes)
	bucket := rt.buckets[ind]
	return bucket
}

// func (n *Network) setupRouteTable(currIP [4]byte, currPort [2]byte) error {

// Usage:
// foundContact := recursiveLookup(desiredNodeID, listContacts, rt)
// if foundContact.node == nodeID{} {
// 		 not found
// }
func recursiveLookup(desiredNodeID nodeID, listContacts []Contact, rt *routingTable) (Contact, error) {
	switch len(listContacts) {
	case 0:
		return Contact{}, nil
	case 1:
		currContact := listContacts[0]

		// If the current Contact is the nodeID we're looking for then we are done.
		if currContact.Node == desiredNodeID {
			return currContact, nil
		}

		// If the current Contact is not the nodeID we're looking for, then
		// continue with the recursive search.
		rt.add(currContact)
		lcs, err := lookup(desiredNodeID, currContact)
		if err != nil {
			return Contact{}, err
		}
		recursiveLookup(desiredNodeID, lcs, rt)
	default:
		for _, currContact := range listContacts {
			rt.add(currContact)
			lcs, err := lookup(desiredNodeID, currContact)
			if err != nil {
				return Contact{}, err
			}
			recursiveLookup(desiredNodeID, lcs, rt)
		}
	}

	return Contact{}, nil
}
