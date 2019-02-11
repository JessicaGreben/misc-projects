package network

import (
	"fmt"
	"sync"

	b "github.com/jessicagreben/kademlia/pkg/bucket"
	"github.com/jessicagreben/kademlia/pkg/node"
	"github.com/jessicagreben/kademlia/pkg/types"
)

const (
	bitsPerByte = 8                            // How many bits in a byte.
	bucketCount = types.IDLength * bitsPerByte // How many buckets should be created for each route table.
)

type routingTable struct {
	boot        types.Contact
	currentNode types.Contact

	// Keep a mapping of node IDs that currently exist in the buckets.
	currentNodesInBucket map[types.NodeID]struct{}

	// One bucket for each bit in the current node's ID.
	buckets [bucketCount]*b.Bucket

	mu sync.Mutex
}

func newRoutingTable(c types.Contact) *routingTable {
	boot := types.Contact{
		NodeID: types.NodeID{0},
		IP:     "boot",
		Port:   "8080",
	}
	var mu sync.Mutex
	rt := routingTable{
		currentNodesInBucket: map[types.NodeID]struct{}{},
		boot:                 boot,
		mu:                   mu,
	}
	for i := 0; i < types.IDLength; i++ {
		rt.buckets[i] = &b.Bucket{}
	}
	rt.currentNode = c
	return &rt
}

func (rt *routingTable) find(c types.Contact) error {
	return nil
}

func (rt *routingTable) add(c types.Contact) error {
	ind := node.FindBucketIndex(c.NodeID, rt.currentNode.NodeID)
	bucket := rt.buckets[ind]
	if err := bucket.Push(c); err != nil {
		return err
	}

	// Add Contact node ID to the currentNodesInBucket.
	rt.mu.Lock()
	rt.currentNodesInBucket[c.NodeID] = struct{}{}
	rt.mu.Unlock()

	fmt.Println("bucket value:", *bucket)
	return nil
}

func (rt *routingTable) update(c types.Contact) error {
	return nil
}

func (rt *routingTable) remove(c types.Contact) error {
	return nil
}

func (rt *routingTable) ping(c types.Contact) error {
	return nil
}

func (rt *routingTable) findClosestNodes(xorPrefix int) []types.Contact {
	for i := xorPrefix; i < 0; i-- {
		currBucket := rt.buckets[i]
		if len(*currBucket) > 0 {
			return *currBucket
		}
	}
	return []types.Contact{}
}
