package bucket

import (
	"github.com/jessicagreben/kademlia/pkg/types"
)

const k = 20 // Max number of contacts in any one bucket.

// Bucket is container of current contacts the.
// Most recently contacted is at the end, least recently contacted is at beginning.
type Bucket []types.Contact

// Shift removes the contact that is at the beginning of the bucket.
// Most recently contacted is at the end, least recently contacted is at beginning.
func (b Bucket) Shift() (Bucket, types.Contact) {

	// First check that there are any Contacts in this bucket.
	if len(b) < 1 {
		return b, types.Contact{}
	}

	return b[:0], b[0]
}

// IsFull checks if the bucket is currently full.
func (b Bucket) IsFull() bool {
	if len(b) == k {
		return true
	}

	// TODO: what if the bucket length is greater than k?
	return false
}

// Push adds a new contact to the end of the bucket.
// Most recently contacted is at the end, least recently contacted is at beginning.
func (b Bucket) Push(c types.Contact) Bucket {

	// Add a Contact to the end of the bucket list.
	return append(b, c)
}

// Find returns the index of the node in the bucket.
func (b Bucket) Find(id types.NodeID) (int, types.Contact, bool) {
	var index int
	for index, contact := range b {
		if id == contact.NodeID {
			return index, contact, true
		}
	}
	return index, types.Contact{}, false
}

// Remove the contact at the index position from the bucket.
func (b Bucket) Remove(index int) Bucket {

	// TODO: return error if index is beyond bucket size.
	if len(b)-1 < index {
		return b
	}
	return append(b[:index], b[index+1:]...)
}
