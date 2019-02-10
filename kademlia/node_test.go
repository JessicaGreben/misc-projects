package main

import (
	"testing"
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
			actualOut := generateNodeID(tt.input)
			if len(actualOut) != idLength {
				t.Errorf("Expected %d, Actual %d", idLength, len(actualOut))
			}
		})
	}
}

func TestNodeDistance(t *testing.T) {
	id1 := nodeID{}
	a := []byte{255}
	copy(id1[:], a)

	id2 := nodeID{}
	b := []byte{0}
	copy(id2[:], b)

	id3 := nodeID{}
	d := []byte{5, 8}
	copy(id3[:], d)

	id4 := nodeID{}
	e := []byte{8, 5}
	copy(id4[:], e)

	id5 := nodeID{}
	f := []byte{13, 13}
	copy(id5[:], f)

	var testCases = []struct {
		name        string
		id1         nodeID
		id2         nodeID
		expectedOut nodeID
	}{
		{"all different", id1, id2, id1},
		{"same value", id1, id1, id2},
		{"two bytes", id3, id4, id5},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actualOut := nodeDistance(tt.id1, tt.id2)
			if actualOut != tt.expectedOut {
				t.Errorf("Expected %08b, Actual %08b", tt.expectedOut, actualOut)
			}
		})
	}
}

func TestFindLongestPrefix(t *testing.T) {
	var testCases = []struct {
		name        string
		in          nodeID
		expectedOut int
	}{
		{"all prefix zeros", nodeID{}, 160},
		{"no prefix zeros", nodeID{255}, 0},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actualOut := findLongestPrefix(tt.in)
			if actualOut != tt.expectedOut {
				t.Errorf("Expected %d, Actual %d", tt.expectedOut, actualOut)
			}
		})
	}
}
