package node

import (
	"testing"

	"github.com/jessicagreben/kademlia/pkg/types"
)

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
			actualOut := GenerateNodeID(tt.input)
			if len(actualOut) != idLength {
				t.Errorf("Expected %d, Actual %d", idLength, len(actualOut))
			}
		})
	}
}

func TestNodeDistance(t *testing.T) {
	id1 := types.NodeID{}
	a := []byte{255}
	copy(id1[:], a)

	id2 := types.NodeID{}
	b := []byte{0}
	copy(id2[:], b)

	id3 := types.NodeID{}
	d := []byte{5, 8}
	copy(id3[:], d)

	id4 := types.NodeID{}
	e := []byte{8, 5}
	copy(id4[:], e)

	id5 := types.NodeID{}
	f := []byte{13, 13}
	copy(id5[:], f)

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
