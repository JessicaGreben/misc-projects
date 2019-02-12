# Implementation of Kademlia Distributed Hash Table

## Spec

http://xlattice.sourceforge.net/components/protocol/kademlia/specs.html


## Routing Implementation

### A node joins the network

Case 1: the node joining is the bootstrap node.

Case 2: the node joining is the first non-boostrap node.

Case 3: the node joining is the second or more non-boostrap node.

Each node does the following when joining the network:
- Create a unique ID.
- Create a routing table.
- Populate the routing table with contact info of other nodes in the network.

#### Create a unique ID

160 bit unique integer.

#### Create a routing table

The routing table should contain the following information:
- the current node's ID
- a data structure for storing contacts of other nodes in the network.

#### Populate routing table contacts

When a node joins the network, at first the only contact its knows of is the boostrap node, since the bootstrap node info is hardcoded into the network.
The first thing the new node does is send a request to the bootstrap node to look for itself in the bootstrap node's routing table.

This request accomplised two main things:
- the bootstrap node now has the contact info of the new node and adds that data to its routing table.
- since the bootstrap node returns a list of nodes that are the closest to current node. Then the current node can add those closest nodes to its routing table.

## Calculating Distance Between Nodes

In Kademlia, distance is defined by how similar two Node IDs are to each other.

If you XOR two Node IDs, the result is an integer that describes which bits where the same between the two IDs.
In Kademlia, Node IDs are read as bigendian numbers, meaning the lowest order bit is the most significant value. 
The XOR value with the most leading zeros (starting left to right) is the "closest".

1 byte examples:

----
decimal value:                     5                 8

binary value:                  `0000 0101`         `0000 1000`

XOR of binary:                          `0000 1101`

XOR value in decimal:                       13

Distance (XOR 0's prefix from left to right):4


----
decimal value:                     5                 31

binary value:                  `00000101`         `00011111`

XOR of binary:                          `00011010`

XOR value in decimal:                       30

Distance (XOR 0's prefix from left to right):3

----
decimal value:                     5                 100

binary value:                  `00000101`         `01100100`

XOR of binary:                          `01100001`

XOR value in decimal:                       97

Distance (XOR 0's prefix from left to right):1

----
decimal value:                          30               31

binary value:                        `00011010`         `00011011`

XOR of binary:                                `00000001`

XOR value in decimal:                             1

Distance (XOR 0's prefix from left to right):     7

----
80 bit (10 bytes) UID example:

`UID1`:            `00111000 10010000 11100000 11110000 00000011 00100001 00000000 01111111 00000100 01111000`

`UID2`:            `00011000 10010000 11011000 00001011 00000011 00100001 00100010 01000111 00000100 01100011`

XOR of `UID1^UID2`: `00100000 00000000 00111000 11111011 00000000 00000000 00100010 00111000 00000000 00011011`

XOR 0's prefix (left to right): 2

## Other descriptions

#### Contact

A contact contains the following information about a node in the network:
- Node ID
- IP 
- Port

#### k

k is a hard coded value for how many contacts can be in any bucket. Typically this is 20.

#### Bucket

A node organizes the contacts into buckets. There is one bucket for each bit in the Node ID. A bucket contains k contacts.  Contacts in the buckets are sorted by most-recently communication, with least-recently communicated at beginning of list.

You can implement the bucket storage in a number of different ways. The original paper describes it as a tree.
However I implement it as an array, where the index position corresponds to the count of prefix 0s of the node IDs. 
