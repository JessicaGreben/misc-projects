package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jessicagreben/misc-projects/simple-repl/pkg/storage"
)

// Run starts the REPL, reading from stdin and executing the commands.
func Run() {
	root := storage.Transaction{
		Operations: map[string]string{},
	}
	currentTx := root

	for {
		var cmd, key, value, out string
		var err error

		var reader = bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')

		cmd, key, value = parseArgs(input)

		currentTx, out, err = storage.ExecuteOp(cmd, key, value, currentTx)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		}

		fmt.Print(out)
	}
}

// parseArgs parses the input from the REPL.
func parseArgs(input string) (cmd string, key string, value string) {
	args := strings.Split(strings.Trim(input, "\n "), " ")

	switch len(args) {
	case 3:
		cmd, key, value = args[0], args[1], args[2]
	case 2:
		cmd, key = args[0], args[1]
	case 1:
		cmd = args[0]
	}

	return cmd, key, value
}
