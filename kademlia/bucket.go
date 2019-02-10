package main

type bucket []Contact

// type bucket [k]Contact

func (b bucket) pop(rt *routingTable) Contact {

	// First check that there are any Contacts in this bucket.
	if len(b) < 1 {
		return Contact{}
	}

	// Find the first element to return.
	Contact := b[0]

	// Delete this first element from the tracking map.
	rt.mu.Lock()
	delete(rt.currentNodesInBucket, Contact.Node)
	rt.mu.Unlock()

	// Remove the first Contact from the bucket.
	b = b[1:]
	return Contact
}

func (b bucket) push(c Contact, rt *routingTable) bucket {

	// Add Contact node ID to the currentNodesInBucket
	rt.mu.Lock()
	rt.currentNodesInBucket[c.Node] = struct{}{}
	rt.mu.Unlock()

	// TODO: first check if bucket is full.
	// If its full, ping first node and remove if unresponsive.
	// if b.isFull() {
	// 	b.ping()
	// }

	// Add a Contact to the end of the bucket list.
	return append(b, c)
}

// Return the index of the node in the bucket.
func (b bucket) find(id nodeID) (int, bool) {
	index := 0
	for index, Contact := range b {
		if id == Contact.Node {
			return index, true
		}
	}
	return index, false
}

// Remove the Contact from the bucket.
func (b bucket) cut(index int, rt *routingTable) Contact {
	c := b[index]

	// Delete this first element from the tracking map.
	rt.mu.Lock()
	delete(rt.currentNodesInBucket, c.Node)
	rt.mu.Unlock()

	b = append(b[:index], b[index+1:]...)
	return c
}

func (b bucket) moveToEnd(id nodeID, rt *routingTable) bucket {
	_, found := b.find(id)
	if found {
		return b
	}
	// contact := b.cut(id, rt)
	// b.push(contact, rt)
	return b
}

func (b bucket) isFull() bool {
	if len(b) == k {
		return true
	}
	return false
}

func (b bucket) update(c Contact, rt *routingTable) bucket {
	_, found := b.find(c.Node)
	if found {
		if !b.isFull() {
			b.push(c, rt)
			return b
		}
		if err := Ping(); err != nil {
			return b
		}
		b.pop(rt)
		b.push(c, rt)
		return b
	}
	b.moveToEnd(c.Node, rt)
	return b
}

// The bucket that a node contact should be placed in is determined by the
// numbering of leading 0 bits in the XOR of the current node ID with the target
// node ID.
func findBucket(node1, node2 nodeID) int {
	xorBytes := nodeDistance(node1, node2)
	bucket := findLongestPrefix(xorBytes)
	return bucket
}
