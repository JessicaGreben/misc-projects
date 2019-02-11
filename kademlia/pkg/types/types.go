package types

// IDLength is how many bytes the Node ID is.
const IDLength = 20

// Contact stores information about how to contact a node in the network.
type Contact struct {
	NodeID NodeID
	IP     string
	Port   string
}

// NodeID is the unique ID of each node in the network.
type NodeID [IDLength]byte
