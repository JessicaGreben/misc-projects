package bucket

import (
	"testing"

	"github.com/jessicagreben/kademlia/pkg/types"
)

func setupBucket(contactCount int) Bucket {
	b := Bucket{}
	for i := 0; i < contactCount; i++ {
		c := types.Contact{
			NodeID: types.NodeID{byte(i)},
		}
		b = b.Push(c)
	}
	return b
}

func TestShift(t *testing.T) {
	bEmpty := setupBucket(1)
	d := types.Contact{
		NodeID: types.NodeID{123},
	}
	bNotEmpty := Bucket{d}

	var testCases = []struct {
		name        string
		bucket      Bucket
		expectedOut types.Contact
		expectedLen int
	}{
		{"empty", bEmpty, types.Contact{}, 0},
		{"full", bNotEmpty, d, 0},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actualB, actualOut := tt.bucket.Shift()
			if tt.expectedOut.NodeID != actualOut.NodeID {
				t.Fatalf("Expect %v, Actual %v", tt.expectedOut.NodeID, actualOut.NodeID)
			}
			if tt.expectedLen != len(actualB) {
				t.Fatalf("Expect %v, Actual %v", tt.expectedLen, len(tt.bucket))
			}
		})
	}
}

func TestIsFull(t *testing.T) {
	bEmpty := Bucket{}
	bFull := setupBucket(20)

	var testCases = []struct {
		name        string
		bucket      Bucket
		expectedOut bool
	}{
		{"empty", bEmpty, false},
		{"full", bFull, true},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actualOut := tt.bucket.IsFull()
			if actualOut != tt.expectedOut {
				t.Fatalf("Expect %v, Actual %v", tt.expectedOut, actualOut)
			}
		})
	}
}

func TestPush(t *testing.T) {
	bEmpty := setupBucket(0)
	bFull := setupBucket(1)

	var testCases = []struct {
		name        string
		bucket      Bucket
		c           types.Contact
		expectedLen int
	}{
		{"empty", bEmpty, types.Contact{}, 1},
		{"one", bFull, types.Contact{}, 2},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actualOut := tt.bucket.Push(tt.c)
			if tt.expectedLen != len(actualOut) {
				t.Fatalf("Expect %v, Actual %v", tt.expectedLen, len(actualOut))
			}
		})
	}
}

func TestFind(t *testing.T) {
	id := types.NodeID{123}
	c := types.Contact{
		NodeID: id,
	}
	bExists := Bucket{c}
	bNotExists := Bucket{}

	var testCases = []struct {
		name          string
		bucket        Bucket
		c             types.Contact
		expectedOut   int
		expectedFound bool
	}{
		{"exists", bExists, c, 0, true},
		{"not exists", bNotExists, types.Contact{}, 0, false},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actualOut, _, actualFound := tt.bucket.Find(tt.c.NodeID)
			if tt.expectedOut != actualOut {
				t.Fatalf("Expect %v, Actual %v", tt.expectedOut, actualOut)
			}
			if tt.expectedFound != actualFound {
				t.Fatalf("Expect %v, Actual %v", tt.expectedFound, actualFound)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	bFirst := setupBucket(3)
	bSecond := setupBucket(3)
	bThird := setupBucket(3)
	second := types.Contact{
		NodeID: types.NodeID{byte(1)},
	}
	third := types.Contact{
		NodeID: types.NodeID{byte(2)},
	}

	var testCases = []struct {
		name            string
		bucket          Bucket
		index           int
		expectedRemoved types.Contact
		expectedLen     int
	}{
		{"first index", bFirst, 0, second, 2},
		{"mid index", bSecond, 1, third, 2},
		{"last index", bThird, 2, second, 2},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actualOut := tt.bucket.Remove(tt.index)
			if tt.expectedLen != len(actualOut) {
				t.Fatalf("Expect %v, Actual %v", tt.expectedLen, len(actualOut))
			}

			if tt.index < len(actualOut) {
				if tt.expectedRemoved.NodeID != actualOut[tt.index].NodeID {
					t.Fatalf("Expect %v, Actual %v", tt.expectedRemoved.NodeID, actualOut[tt.index].NodeID)
				}
			} else {
				if tt.expectedRemoved.NodeID != actualOut[len(actualOut)-1].NodeID {
					t.Fatalf("Expect %v, Actual %v", tt.expectedRemoved.NodeID, actualOut[len(actualOut)-1].NodeID)
				}
			}
		})
	}

}
