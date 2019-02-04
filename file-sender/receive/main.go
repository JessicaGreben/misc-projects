package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatalln("Missing arguments. Usage: ./receive <relay-host>:<relay-port> <secret-code> <output-directory>")
	}

	addr := os.Args[1]
	if _, err := url.ParseRequestURI(addr); err != nil {
		log.Fatalln("Address not formed correctly. Expected <host>:<port>, Actual:", addr)
	}

	secret := os.Args[2]

	dir := os.Args[3]
	if _, err := os.Stat(dir); err != nil {
		log.Fatalln("os.Stat dir err:", err)
	}

	if err := receive(addr, secret, dir); err != nil {
		log.Fatalln("receive err:", err)
	}
}

func receive(addr, secret, dir string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := sendHeader(conn, secret); err != nil {
		return err
	}

	n, err := createFile(conn, dir)
	if err != nil {
		return err
	}

	log.Printf("Processed bytes from relay: code %s, bytes %d", secret, n)
	return nil
}

func sendHeader(conn net.Conn, secret string) error {
	var header struct {
		Request byte
		Secret  int32
	}

	// receiveCmd tells the relay server that the request is a receive request.
	const receiveCmd = 1
	header.Request = receiveCmd

	sec, err := strconv.Atoi(secret)
	if err != nil {
		return err
	}
	header.Secret = int32(sec)

	if err := binary.Write(conn, binary.LittleEndian, &header); err != nil {
		return err
	}

	return nil
}

func createFile(conn net.Conn, dir string) (int, error) {

	// Assume a file name with not exceed 255 chars, with max 4 bytes per char.
	const fileNameSize = 1024

	var header struct {
		FileName [fileNameSize]byte
		FileSize int64
	}

	// Get the file name and size from the relay server.
	if err := binary.Read(conn, binary.LittleEndian, &header); err != nil {
		return 0, err
	}

	// Trim off extra zero value bytes if the name is less than fileNameSize.
	fileName := bytes.Trim(header.FileName[:], "\x00")
	filePath := fmt.Sprintf("%s/%s", dir, fileName)
	f, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	const dataBuffer = 1024
	var bytesProcessed int
	for {
		data := make([]byte, dataBuffer)
		n, err := conn.Read(data)
		bytesProcessed += n

		switch {
		case err != nil:
			if n > 0 {
				if _, err = f.Write(data[:n]); err != nil {
					return bytesProcessed, err
				}
			}
			return bytesProcessed, err
		case bytesProcessed > int(header.FileSize):
			return bytesProcessed, fmt.Errorf("critical: read more bytes than expected. Expected: %d, Actual: %d", header.FileSize, bytesProcessed)
		default:
			if _, err = f.Write(data[:n]); err != nil {
				return bytesProcessed, err
			}
			if bytesProcessed == int(header.FileSize) {
				return bytesProcessed, nil
			}
		}
	}
}
