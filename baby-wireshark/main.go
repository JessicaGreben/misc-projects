package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
)

const linkTypeEthernet = 1
const ipv4Ethertype = 8
const protocolNumberTCP = 6

// globalHeader is the pcap file global header.
// ref: https://www.tcpdump.org/manpages/pcap-savefile.5.txt
type globalHeader struct {
	MagicNumber     [4]byte
	MajorVersion    int16
	MinorVersion    [2]byte
	ZoneOffset      [4]byte
	TimeAccuracy    [4]byte
	SnapshotLength  uint32
	LinkLayerHeader [4]byte
}

// packetHeader is the pcap per-packet header.
type packetHeader struct {
	TimestampSec       [4]byte
	TimestampMs        [4]byte
	PacketSize         uint32
	OriginalPacketSize uint32
}

// frameHeader is the per packet link layer frame header. Layer 2 Ethernet frame.
// ref: https://en.wikipedia.org/wiki/Ethernet_frame#Structure
type frameHeader struct {
	DestAddr   [6]byte
	SourceAddr [6]byte
	EtherType  uint16
}

// datagramHeader is the per packet IP layer datagram header.
// ref: https://tools.ietf.org/html/rfc791#page-11
type datagramHeader struct {
	Version        uint8
	_              [1]byte
	TotalLength    uint16
	ID             [2]byte
	_              [2]byte
	TTL            [1]byte
	Protocol       uint8
	HeaderChecksum [2]byte
	Source         [4]byte
	Dest           [4]byte
}

// segmentHeader is the per packet transport layer segment header.
// ref: https://tools.ietf.org/html/rfc793#section-3.1
type segmentHeader struct {
	SourcePort uint16
	DestPort   uint16
	Sequence   uint32
	AckNumber  [4]byte
	DataOffset [1]byte
}

func main() {
	data, err := ioutil.ReadFile("./net.cap")
	if err != nil {
		fmt.Printf("ioutil.ReadFile err: %v\n", err)
		return
	}

	// Place raw bytes into the buffer for processing.
	buffer := bytes.NewBuffer(data)

	// Read the pcap-savefile global header.
	if err := readGlobalHeader(buffer); err != nil {
		fmt.Printf("readGlobalHeader err: %v\n", err)
		return
	}

	// This mapping is the tcp sequence number to http payload bytes
	// so that we can reassemble the bytes in the correct order.
	httpData := make(map[int][]byte)

	// Read all of the packets.
	for {
		httpData, err = readPacket(buffer, httpData)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("readPacket err: %v\n", err)
			return
		}
	}

	// Create a file from the httpData and open that file.
	createFile(httpData)
}

func readPacket(buffer *bytes.Buffer, httpOrder map[int][]byte) (map[int][]byte, error) {
	pacHeader, err := readPacHeader(buffer)
	if err != nil {
		return httpOrder, err
	}

	readFrame(buffer)

	datagram, err := readDatagram(buffer)
	if err != nil {
		return httpOrder, err
	}

	datagramHeader := 4 * int(datagram.Version&0x0f)

	segment, err := readSegment(buffer)
	if err != nil {
		return httpOrder, err
	}
	tcpHeader := 4 * int(segment.DataOffset[0]>>4)

	// Read HTTP header and data.
	// The HTTP payload length is the datagram total length minus the datagram and
	//  tcp header lengths.
	httpLength := int(datagram.TotalLength) - datagramHeader - tcpHeader
	httpData := make([]byte, httpLength)
	buffer.Read(httpData)

	// We know we only want http data from packets
	// that are destined to this addr.
	ip := [4]byte{192, 168, 0, 101}
	if httpLength > 0 && datagram.Dest == ip {

		// The http header is separated by \r\n\r\n. We can tell this
		// packet contains the header if this is present.
		headerPresent := bytes.Split(httpData, []byte{'\r', '\n', '\r', '\n'})
		if len(headerPresent) > 1 {
			httpOrder[int(segment.Sequence)] = headerPresent[1]
		} else {
			httpOrder[int(segment.Sequence)] = httpData
		}
	}

	// Weird hack to remove the random extra 2 bytes that show up randomly
	// on packets of this size.
	if pacHeader.PacketSize == 68 || pacHeader.PacketSize == 76 {
		buffer.Read(make([]byte, 2))
	}

	return httpOrder, nil
}

func readGlobalHeader(buffer *bytes.Buffer) error {
	globalHeader := globalHeader{}
	if err := binary.Read(buffer, binary.BigEndian, &globalHeader); err != nil {
		return err
	}

	// Make sure link layer header is Ethernet.
	// ref: https://www.tcpdump.org/linktypes.html
	if globalHeader.LinkLayerHeader[0] != linkTypeEthernet {
		return fmt.Errorf("wrong link layer header. Expected 1, but got %d", globalHeader.LinkLayerHeader)
	}

	return nil
}

func readPacHeader(buffer *bytes.Buffer) (packetHeader, error) {
	pacHeader := packetHeader{}
	if err := binary.Read(buffer, binary.LittleEndian, &pacHeader); err != nil {
		if err == io.EOF {
			return pacHeader, err
		}
		return pacHeader, err
	}

	// This won't always be the case, but for this file we know that the packet
	// size is the same as the original size. I.e. the packet is never trucated.
	if pacHeader.PacketSize != pacHeader.OriginalPacketSize {
		return pacHeader, fmt.Errorf("wrong pcap pac header. Expected %d, but got %d", pacHeader.PacketSize, pacHeader.OriginalPacketSize)
	}
	return pacHeader, nil
}

func readFrame(buffer *bytes.Buffer) error {
	frame := frameHeader{}
	if err := binary.Read(buffer, binary.BigEndian, &frame); err != nil {
		return err
	}

	// We know that the ether type is 8, indicating protocol IPv4 for the payload.
	if frame.EtherType != ipv4Ethertype {
		return fmt.Errorf("wrong etherType. Expected 8, but got %d", frame.EtherType)
	}
	return nil
}

func readDatagram(buffer *bytes.Buffer) (datagramHeader, error) {
	datagram := datagramHeader{}
	if err := binary.Read(buffer, binary.BigEndian, &datagram); err != nil {
		return datagram, err
	}

	// Make sure the protocol is 6 for TCP.
	if datagram.Protocol != protocolNumberTCP {
		return datagram, fmt.Errorf("wrong IP Protocol. Expected 6, but got %d", datagram.Protocol)
	}

	return datagram, nil
}

func readSegment(buffer *bytes.Buffer) (segmentHeader, error) {
	segment := segmentHeader{}
	if err := binary.Read(buffer, binary.BigEndian, &segment); err != nil {
		return segment, err
	}

	tcpHeader := 4 * int(segment.DataOffset[0]>>4)

	// We only read in 13 bytes so far, but the segment header is bigger than that,
	// so we need to read the rest of it.
	readTheRestTCP := tcpHeader - 13
	buffer.Read(make([]byte, readTheRestTCP))

	return segment, nil
}

func createFile(httpOrder map[int][]byte) error {

	// Sort the values of the TCP sequence numbers.
	var sequenceNums []int
	for k := range httpOrder {
		sequenceNums = append(sequenceNums, k)
	}
	sort.Ints(sequenceNums)

	f, err := os.Create("./packet.jpeg")
	if err != nil {
		return err
	}
	defer f.Close()

	// Write the HTTP data in order to a jpeg file.
	for _, k := range sequenceNums {
		data := httpOrder[k]
		f.Write(data)
	}
	f.Sync()

	// Open the jpeg that was created. This assumes mac OSX.
	if err := exec.Command("/usr/bin/open", "./packet.jpeg").Run(); err != nil {
		return err
	}

	return nil
}
