package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"net"
	"net/url"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("Missing arguments. Usage: ./send <relay-host>:<relay-port> <file-to-send>")
	}

	addr := os.Args[1]
	if _, err := url.ParseRequestURI(addr); err != nil {
		log.Fatalln("Address not formed correctly. Expected <host>:<port>, Actual:", addr)
	}

	fileName := os.Args[2]
	if _, err := os.Stat(fileName); err != nil {
		log.Fatalln("os.Stat fileName err:", err)
	}

	if err := send(addr, fileName); err != nil {
		log.Fatalln("send err:", err)
	}
}

func send(addr, fileName string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := sendHeader(conn, fileName); err != nil {
		return err
	}

	if _, err := sendFile(conn, fileName); err != nil {
		return err
	}
	return nil
}

func sendHeader(conn net.Conn, fileName string) error {

	// Assume a file name with not exceed 255 chars, with max 4 bytes per char.
	const fileNameSize = 1024
	var header struct {
		Request  byte
		Secret   int32
		FileName [fileNameSize]byte
		FileSize int64
	}

	// sendCmd tells the relay server that the request is a send request.
	const sendCmd = 2
	header.Request = sendCmd

	header.Secret = generateSecret()
	copy(header.FileName[:], fileName)

	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	}
	header.FileSize = fi.Size()

	if err := binary.Write(conn, binary.LittleEndian, &header); err != nil {
		return err
	}

	// Print the secret to stdout after the header is successfully sent to
	// the relay server.
	fmt.Println(header.Secret)
	return nil
}

func generateSecret() int32 {
	mrand.Seed(time.Now().UnixNano())
	secret := int32(mrand.Int())
	return secret
}

func sendFile(conn net.Conn, fileName string) (int, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	const dataBuffer = 1024
	var bytesProcessed int
	for {
		data := make([]byte, dataBuffer)
		n, err := f.Read(data)
		bytesProcessed += n

		switch {
		case err == io.EOF:
			if n > 0 {
				if _, err := conn.Write(data[:n]); err != nil {
					return bytesProcessed, err
				}
			}
			return bytesProcessed, nil
		case err != nil:
			return bytesProcessed, err
		default:
			if _, err := conn.Write(data[:n]); err != nil {
				return bytesProcessed, err
			}
		}
	}
}
