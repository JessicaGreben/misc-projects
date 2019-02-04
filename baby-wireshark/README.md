# Baby Wireshark

## Description:

We have recorded a packet capture of an HTTP request and response for an image, performed over an imperfect network. 

The challenge is to parse the capture file, located at `net.cap`, find and parse the packets constituting the image download, and reconstruct the image!

## Details

`tcpdump` was used to capture the network traffic and create the `net.cap` file.

The file is saved as “pcap-savefile” format. Read more about that format here:
https://www.tcpdump.org/manpages/pcap-savefile.5.txt

## Getting Started

### Prerequisites

- Golang needs to be [installed](https://golang.org/doc/install).

- Assumes you can open a jpeg with an `open` command located at `/usr/bin/open`

```
$ which open
/usr/bin/open
```

### Run

The following command should execute the program and open the image that 
was reassembled from the packet capture.

`go run main.go`
