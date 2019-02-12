package node

import (
	"testing"

	"github.com/jessicagreben/kademlia/pkg/types"
)

func setupNodeID(value []byte) types.NodeID {
	n := types.NodeID{}
	copy(n[:], value)
	return n
}

func TestGenerateNodeID(t *testing.T) {
	var testCases = []struct {
		name  string
		input int
	}{
		{"zero", 0},
		{"twenty", 20},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actualOut := GenerateID(tt.input)
			if len(actualOut) != idLength {
				t.Errorf("Expected %d, Actual %d", idLength, len(actualOut))
			}
		})
	}
}

func TestNodeDistance(t *testing.T) {
	id1 := setupNodeID([]byte{255})
	id2 := setupNodeID([]byte{0})
	id3 := setupNodeID([]byte{5, 8})
	id4 := setupNodeID([]byte{8, 5})
	id5 := setupNodeID([]byte{13, 13})

	var testCases = []struct {
		name        string
		id1         types.NodeID
		id2         types.NodeID
		expectedOut types.NodeID
	}{
		{"all different", id1, id2, id1},
		{"same value", id1, id1, id2},
		{"two bytes", id3, id4, id5},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actualOut := Distance(tt.id1, tt.id2)
			if actualOut != tt.expectedOut {
				t.Errorf("Expected %08b, Actual %08b", tt.expectedOut, actualOut)
			}
		})
	}
}

func TestFindLongestPrefix(t *testing.T) {
	var testCases = []struct {
		name        string
		in          types.NodeID
		expectedOut int
	}{
		{"all prefix zeros", types.NodeID{}, 160},
		{"no prefix zeros", types.NodeID{255}, 0},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actualOut := FindLongestPrefix(tt.in)
			if actualOut != tt.expectedOut {
				t.Errorf("Expected %d, Actual %d", tt.expectedOut, actualOut)
			}
		})
	}
}
