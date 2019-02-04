package storage

import (
	"fmt"
	"os"
	"strings"
)

// Transaction contains a map of key/value operations and a pointer to it's parent transaction.
type Transaction struct {
	Operations map[string]string
	parent     *Transaction
}

// ExecuteOp executes the operation passed into the REPL.
func ExecuteOp(cmd, key, value string, currentTx Transaction) (Transaction, string, error) {
	var output string
	var err error

	switch strings.ToUpper(cmd) {
	case "READ":
		{
			output, err = readTx(key, currentTx)
			output += "\n"
		}
	case "WRITE":
		{
			writeTx(key, value, currentTx)
		}
	case "DELETE":
		{
			deleteTx(key, currentTx)
		}
	case "START":
		{
			currentTx = startTx(currentTx)
		}
	case "COMMIT":
		{
			currentTx, err = commitTx(currentTx)
		}
	case "ABORT":
		{
			currentTx, err = abortTx(currentTx)
		}
	case "QUIT":
		{
			fmt.Fprintf(os.Stderr, "Exiting...\n")
			os.Exit(0)
		}
	}
	return currentTx, output, err
}
