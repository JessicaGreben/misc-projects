package network

import (
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
	buckets [bucketCount]b.Bucket

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
		rt.buckets[i] = b.Bucket{}
	}
	rt.currentNode = c
	return &rt
}

func (rt *routingTable) find(id types.NodeID) (types.Contact, bool) {
	_, found := rt.currentNodesInBucket[id]
	ind := node.FindBucketIndex(id, rt.currentNode.NodeID)
	bucket := rt.buckets[ind]
	_, c, found := bucket.Find(id)
	if found {
		return c, true
	}
	return types.Contact{}, false
}

func (rt *routingTable) add(newContact types.Contact) b.Bucket {
	ind := node.FindBucketIndex(newContact.NodeID, rt.currentNode.NodeID)
	bucket := rt.buckets[ind]

	// If the bucket is currently full, then ping each contact to see if it is
	// still responsive. If it isn't then remove it and add the new contact.
	if bucket.IsFull() {

		// Ping each contact in the bucket.
		for ind, existingContact := range bucket {
			responsive, _ := Ping(existingContact)

			// If an existing contact is not responsive,
			// then remove it and add the new contact.
			if !responsive {
				bucket = bucket.Remove(ind)
				break
			}

			// If all existing contacts are responsive then ignore
			// the new contact and do not add it.
			if ind == len(bucket)-1 {
				return bucket
			}
		}
	}
	bucket = bucket.Push(newContact)
	rt.mu.Lock()
	rt.currentNodesInBucket[newContact.NodeID] = struct{}{}
	rt.mu.Unlock()

	rt.buckets[ind] = bucket
	return bucket
}

func (rt *routingTable) update(c types.Contact) error {
	return nil
}

func (rt *routingTable) remove(c types.Contact) b.Bucket {
	return nil
}

func (rt *routingTable) findClosestNodes(xorPrefix int) []types.Contact {

	// First look for the closest nodes.
	for i := xorPrefix; i >= 0; i-- {
		currBucket := rt.buckets[i]
		if len(currBucket) > 0 {
			return currBucket
		}
	}

	// If there aren't any "close" nodes, then return
	// any nodes.
	for i := xorPrefix + 1; i < bucketCount; i++ {
		currBucket := rt.buckets[i]
		if len(currBucket) > 0 {
			return currBucket
		}
	}
	return []types.Contact{}
}
