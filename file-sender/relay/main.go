package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"sync"
	"time"
)

const (

	// The chan (fileStoreage.data) is a slice of bytes, where each slice is dataBuffer (1024) bytes.
	// So that we don't exceed maxDataStorage (4MB), create a buffer with chanBufferSize.
	chanBufferSize = maxDataStorage / dataBuffer
	maxDataStorage = 4 * 1024 * 1024 // Max storage of 4MB.
	dataBuffer     = 1024            // Read dataBuffer bytes at a time of file data.
	fileNameSize   = 1024            // Assume a file name with not exceed 255 chars, with max 4 bytes per char.
)

const (
	receiveCmd = 1 // receiveCmd indicates that relay should execute the receiver code.
	sendCmd    = 2 // sendCmd indicates that relay should execute the sender code.
)

type secretCode int32

type fileStorage struct {
	name [fileNameSize]byte
	size int64
	data chan []byte
}

type dataStore struct {
	fs map[secretCode]*fileStorage
	mu sync.RWMutex
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Missing arguments. Usage: ./relay :<port>")
	}

	port := os.Args[1]
	var rxPat = regexp.MustCompile(`^:\d{4,5}$`)
	if !rxPat.MatchString(port) {
		log.Fatalf("Port not formed correctly. Expected :<port>, Actual: %s\n", port)
	}

	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	defer ln.Close()
	log.Printf("Listening on port %s", port)

	ds := dataStore{
		fs: make(map[secretCode]*fileStorage),
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("ln.Accept err:", err)
			continue
		}

		go handleConn(conn, &ds)
	}
}

func handleConn(conn net.Conn, dataStore *dataStore) {

	// Read the first byte of the request which will indicate
	// if the request is from the sender or receiver.
	var request byte
	if err := binary.Read(conn, binary.LittleEndian, &request); err != nil {
		return
	}

	switch request {
	case sendCmd:
		processSender(conn, dataStore)
	case receiveCmd:
		processReceiver(conn, dataStore)
	default:
		log.Println("No request provided.")
	}
}

func processSender(conn net.Conn, dataStore *dataStore) error {
	start := time.Now()
	var header struct {
		Secret   secretCode
		FileName [fileNameSize]byte
		FileSize int64
	}

	// Read the header from the sender which contains information
	// about the file that the sender is about to send.
	if err := binary.Read(conn, binary.LittleEndian, &header); err != nil {
		return err
	}

	// Create a file store object to store the information about the file
	// that the sender is about to send.
	fileStore := fileStorage{
		name: header.FileName,
		size: header.FileSize,
		data: make(chan []byte, chanBufferSize),
	}

	dataStore.mu.Lock()
	{
		// Add the file store object from above to a map that will be available to the
		// reciever. The key is the secret code and the value is the file store object.
		if _, ok := dataStore.fs[header.Secret]; ok {
			return errors.New("there is already data saved for the secret code")
		}
		dataStore.fs[header.Secret] = &fileStore
	}
	dataStore.mu.Unlock()

	// Read all the file bytes off the wire and send them down a channel that
	// will block when 4MB of data is stored.
	n, err := processSenderFile(conn, header.FileSize, &fileStore)
	if err != nil {
		return err
	}

	log.Printf("Processed bytes from sender: code %d, bytes %d, time %.5fs", header.Secret, n, time.Since(start).Minutes())
	return nil
}

func processSenderFile(conn net.Conn, fileSize int64, fileStore *fileStorage) (int, error) {
	var bytesProcessed int
	for {
		data := make([]byte, dataBuffer)
		n, err := conn.Read(data)
		bytesProcessed += n

		switch {
		case err != nil:
			if n > 0 {
				fileStore.data <- data[:n]
			}
			return bytesProcessed, err
		case bytesProcessed > int(fileSize):
			return bytesProcessed, fmt.Errorf("critical: read more bytes than expected. Expected: %d, Actual: %d", fileSize, bytesProcessed)
		default:
			fileStore.data <- data[:n]
			if bytesProcessed == int(fileSize) {
				close(fileStore.data)
				return bytesProcessed, nil
			}
		}
	}
}

func processReceiver(conn net.Conn, dataStore *dataStore) error {
	start := time.Now()

	// Read the secret code that the receiever sent.
	var header struct {
		Secret secretCode
	}
	if err := binary.Read(conn, binary.LittleEndian, &header); err != nil {
		return err
	}

	// Confirm the secret code exists in the data store structure and if it does
	// retrieve fileStorage.
	dataStore.mu.RLock()
	fileStore, ok := dataStore.fs[header.Secret]
	if !ok {
		return errors.New("invalid secret code")
	}
	dataStore.mu.RUnlock()

	// Send information about the file to the receiver before the file data
	// is sent.
	fileHeader := struct {
		FileName [fileNameSize]byte
		FileSize int64
	}{
		FileName: fileStore.name,
		FileSize: fileStore.size,
	}
	if err := binary.Write(conn, binary.LittleEndian, &fileHeader); err != nil {
		return err
	}

	// Write all the bytes from the file.data channel to the receiver.
	n, err := processReceiverFile(conn, fileStore)
	if err != nil {
		return err
	}

	// Once the receiver client has all file data then delete the fileStore
	// from the dataStore.
	delete(dataStore.fs, header.Secret)

	log.Printf("Processed bytes to receiver: code %d, bytes %d, time %.5fs", header.Secret, n, time.Since(start).Minutes())
	return nil
}

func processReceiverFile(conn net.Conn, fileStore *fileStorage) (int, error) {
	var bytesProcessed int
	for data := range fileStore.data {
		n, err := conn.Write(data)
		bytesProcessed += n
		if err != nil {
			return bytesProcessed, err
		}
	}

	return bytesProcessed, nil
}
