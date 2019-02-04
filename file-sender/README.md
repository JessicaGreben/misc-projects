# File Sending Exercise

## Description

Write a file sending service.

Suppose you have two laptops (A and B) and a server. The user of laptop A wants to send a file to the user of laptop B.

Write three programs:

1. The sender - this program will send a file to the relay server.
2. The receiver - this program will retrieve a file from the relay server.
3. The relay - this program can recieve and send files.

## Prerequisites

* Golang >= v1.11 [installed](https://golang.org/dl/).

## Run

* Build and run the relay server.

`go build`

`./relay :<port>`


* Send a file to the relay server using the sender program

`go build`

`./send <relay-host>:<relay-port> <file-to-send>`

* Receive a file from the relay server using the receiver program.

`go build`

`./receive <relay-host>:<relay-port> <secret-code> <output-directory>`

## Tests

Run tests with code coverage and verbose output.

`go test -v ./... -cover`
