package main

import (
	"net"
	"testing"
)

func setupFs() *fileStorage {
	name := [fileNameSize]byte{}
	copy(name[:], []byte("test.txt")[:])
	fs := fileStorage{
		name: name,
		size: int64(5),
		data: make(chan []byte, chanBufferSize),
	}
	return &fs
}

func TestProcessSenderFile(t *testing.T) {
	senderTestCases := []struct {
		name        string
		size        int64
		fileStorage *fileStorage
		clientBytes []byte
		expected    int
	}{
		{
			name:        "5 Bytes",
			size:        int64(5),
			fileStorage: setupFs(),
			clientBytes: []byte("11111"),
			expected:    5,
		},
		{
			name:        "0 Bytes",
			size:        int64(0),
			fileStorage: setupFs(),
			clientBytes: []byte(""),
			expected:    0,
		},
	}

	for _, tc := range senderTestCases {
		t.Run(tc.name, func(t *testing.T) {
			server, client := net.Pipe()
			ch2 := make(chan int)
			go func() {
				actual, err := processSenderFile(server, tc.size, tc.fileStorage)
				if err != nil {
					t.Fatal(err)
				}
				ch2 <- actual
				close(ch2)
				server.Close()
			}()

			if _, err := client.Write(tc.clientBytes); err != nil {
				t.Fatal(err)
			}
			actual := <-ch2
			if actual != tc.expected {
				t.Errorf("processSenderFile: expected %d, actual %d", tc.expected, actual)
			}
			client.Close()
		})
	}
}

func TestProcessReceiverFile(t *testing.T) {
	receiverTestCases := []struct {
		name        string
		fileStorage *fileStorage
		serverBytes []byte
		expected    int
	}{
		{
			name:        "Five Bytes",
			fileStorage: setupFs(),
			serverBytes: []byte("11111"),
			expected:    5,
		},
		{
			name:        "Zero bytes",
			fileStorage: setupFs(),
			serverBytes: []byte(""),
			expected:    0,
		},
	}

	for _, tc := range receiverTestCases {
		t.Run(tc.name, func(t *testing.T) {
			ch := tc.fileStorage.data
			ch <- tc.serverBytes
			close(ch)

			server, client := net.Pipe()
			go func() {
				_, err := processReceiverFile(server, tc.fileStorage)
				if err != nil {
					t.Fatal(err)
				}
			}()

			actual, err := client.Read(tc.serverBytes)
			if err != nil {
				t.Fatal(err)
			}
			if actual != tc.expected {
				t.Errorf("processReceiverFile: expected %d, actual %d", tc.expected, actual)
			}
			client.Close()
		})
	}
}
