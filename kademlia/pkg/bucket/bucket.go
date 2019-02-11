package bucket

import (
	"fmt"

	"github.com/jessicagreben/kademlia/pkg/types"
)

const k = 20 // Max number of contacts in any one bucket.

// Bucket is container of current contacts the.
// Most recently contacted is at the end, least recently contacted is at beginning.
type Bucket []types.Contact

func (b Bucket) shift() types.Contact {

	// First check that there are any Contacts in this bucket.
	if len(b) < 1 {
		return types.Contact{}
	}

	// Find the first element to return.
	c := b[0]

	// Remove the first Contact from the bucket.
	b = b[1:]

	// TODO: this logic should be moved out to a route table function.
	// Delete this first element from the tracking map.
	// rt.mu.Lock()
	// delete(rt.currentNodesInBucket, c.Node)
	// rt.mu.Unlock()

	return c
}

func (b Bucket) isFull() bool {
	if len(b) == k-1 {
		return true
	}
	return false
}

// Push does x.
func (b *Bucket) Push(c types.Contact) error {

	fmt.Println("1. In push")

	// TODO: first check if bucket is full.
	// If its full, ping first node and remove if unresponsive.
	count := 0
	if b.isFull() {
		fmt.Println("2. In push")

		// Try to make space, by pinging each node currently in the bucket.
		for count < k {
			for _, c := range *b {
				found := b.ping(c)
				if !found {
					// b.cut(ind, rt)
					break
				}
				count++
			}
		}

		// If we get here, then all nodes in the bucket are responsive
		// and we cannot add this new node
		return nil
	}

	fmt.Println("3. In push")

	// Add a Contact to the end of the bucket list.
	*b = append(*b, c)
	fmt.Println("4. In push")
	fmt.Println("In push, new bucket content:", b)
	return nil
}

//TODO: move to net package
func (b Bucket) ping(c types.Contact) bool {

	// 	// Make Ping request.
	// 	if err := Ping(c); err != nil {
	// 		return false
	// 	}

	return true
}

// Return the index of the node in the bucket.
func (b Bucket) find(id types.NodeID) (int, bool) {
	index := 0
	for index, Contact := range b {
		if id == Contact.NodeID {
			return index, true
		}
	}
	return index, false
}

// Remove the Contact from the bucket.
func (b Bucket) cut(index int) types.Contact {
	c := b[index]

	// Delete this first element from the tracking map.
	// rt.mu.Lock()
	// delete(rt.currentNodesInBucket, c.Node)
	// rt.mu.Unlock()

	b = append(b[:index], b[index+1:]...)
	return c
}

// MoveToEnd is x.
func (b *Bucket) MoveToEnd(id types.NodeID) error {
	_, found := b.find(id)
	if found {
		// contact := b.cut(ind, rt)
		// b.push(contact, rt)
		return nil
	}
	return nil
}

func (b *Bucket) update(c types.Contact) error {
	_, found := b.find(c.NodeID)
	if found {

		// If the Contact is in the bucket already, then move it to the end
		// since its recently communicated with.
		b.MoveToEnd(c.NodeID)
		return nil
	}

	// If the contact doesn't already exist in the bucket list,
	// then add it to the end of the bucket list.
	// b.push(c, rt)
	return nil
}
