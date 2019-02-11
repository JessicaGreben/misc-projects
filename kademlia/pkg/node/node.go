package node

import (
	"crypto/rand"
	"fmt"

	"github.com/jessicagreben/kademlia/pkg/types"
)

const (
	idLength  = 20 // Length in bytes of the Node ID.
	keyLength = 20 // Length in bytes of the key for storing data.
	alpha     = 3  // System wide concurrency parameter.
)

// GenerateID does x.
func GenerateID(idLength int) types.NodeID {
	id := types.NodeID{}
	idSlice := make([]byte, idLength)
	_, err := rand.Read(idSlice)
	if err != nil {
		fmt.Println("Error: ", err)
		return id
	}

	copy(id[:], idSlice)
	return id
}

// Distance is x.
func Distance(node1, node2 types.NodeID) types.NodeID {
	xorBytes := types.NodeID{}

	// Iterate over each byte of node1 and node2 ID and XOR each byte.
	for i := 0; i < idLength; i++ {
		xorBytes[i] = node1[i] ^ node2[i]
	}

	return xorBytes
}

// FindLongestPrefix is x.
func FindLongestPrefix(xorBytes types.NodeID) int {
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

// FindBucketIndex is x.
// The bucket that a node contact should be placed in is determined by the
// numbering of leading 0 bits in the XOR of the current node ID with the target
// node ID.
func FindBucketIndex(node1, node2 types.NodeID) int {
	xorBytes := Distance(node1, node2)
	bucketIndex := FindLongestPrefix(xorBytes)
	return bucketIndex
}
