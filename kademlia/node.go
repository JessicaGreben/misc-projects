package main

import (
	"crypto/rand"
	"fmt"
	"sync"
)

// Contact stores information about how to contact a node in the network.
type Contact struct {
	Node   nodeID
	IpAddr string
	Port   string
}

type nodeID [idLength]byte

func generateNodeID(idLength int) nodeID {
	id := nodeID{}
	idSlice := make([]byte, idLength)
	_, err := rand.Read(idSlice)
	if err != nil {
		fmt.Println("Error: ", err)
		return id
	}

	copy(id[:], idSlice)
	return id
}

func nodeDistance(node1, node2 nodeID) nodeID {
	xorBytes := nodeID{}

	// Iterate over each byte of node1 and node2 ID and XOR each byte.
	for i := 0; i < idLength; i++ {
		xorBytes[i] = node1[i] ^ node2[i]
	}

	return xorBytes
}

func findLongestPrefix(xorBytes nodeID) int {
	var prefix int

	// Iterate over each byte in xorBytes.
	// As bigendian, the byte at index 0 is the highest order byte.
	for i := 0; i < len(xorBytes); i++ {

		// Iterate over one bit at a time.
		// The lowest order bit (left most) bit is the highest order bit as bigendian.
		for bit := 7; bit >= 0; bit-- {

			// We use bit masking here to check if the value of each bit is zero.
			// As we iterate through bit values, the mask value will be: 128, 64, 32, 16, 8, 4, 2, 1
			mask := byte(1 << uint(bit))

			// We use the bitwise AND operation is performed on the bit mask value and
			// an xorByte to find where the first non-zero bit occurence is.
			b := xorBytes[i] & mask

			// When we encounter the first non-zero occurence that means we found
			// where the longest common prefix ends.
			if b != 0 {
				return prefix
			}
			prefix++
		}
	}

	return prefix
}

type routingTable struct {
	boot        Contact
	currentNode Contact

	// Keep a mapping of node IDs that currently exist in the buckets.
	currentNodesInBucket map[nodeID]struct{}

	// One bucket for each bit in the current node's ID.
	buckets [bucketCount]bucket

	mu sync.Mutex
}

func newRoutingTable(c Contact) *routingTable {
	boot := Contact{
		Node:   nodeID{0},
		IpAddr: "localhost",
		Port:   "8080",
	}
	var mu sync.Mutex
	rt := routingTable{
		boot: boot,
		mu:   mu,
	}
	for i := 0; i < idLength; i++ {
		rt.buckets[i] = bucket{}
	}
	rt.currentNode = c
	return &rt
}

func (rt *routingTable) add(c Contact) {
	bucket := getBucket(c.Node, rt)
	bucket.push(c, rt)
}
