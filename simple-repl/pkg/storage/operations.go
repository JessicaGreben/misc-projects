package storage

import (
	"errors"
	"fmt"
)

// readTx reads the value stored at the key if it exists.
func readTx(key string, currentTx Transaction) (string, error) {

	// Check if the key exists in the operations of the current transaction.
	if opValue, keyExists := currentTx.Operations[key]; keyExists {

		// The key has been deleted if the value is a string's zero value.
		if opValue == "" {
			return "", fmt.Errorf("Key not found: %s", key)
		}

		return opValue, nil
	}

	// The key does not exist in current transaction so
	// check if it exists in it's parent transaction.
	if currentTx.parent != nil {
		return readTx(key, *currentTx.parent)
	}

	// The key does not exist in the current transaction and
	// there is no parent transaction.
	return "", fmt.Errorf("Key not found: %s", key)
}

// writeTx adds or updates a key/value in the current transactions operations.
func writeTx(key, value string, currentTx Transaction) {
	currentTx.Operations[key] = value
	return
}

// deleteTx "deletes" a value from the current transaction.
// A deletion is represented by a string's zero value.
func deleteTx(key string, currentTx Transaction) {
	currentTx.Operations[key] = ""
	return
}

// startTx creates a new child transaction where the currentTx is the parent.
func startTx(currentTx Transaction) Transaction {
	childTx := Transaction{
		Operations: map[string]string{},
		parent:     &currentTx,
	}
	return childTx
}

// commitTx commits all operations in the current transaction to its parent
func commitTx(currentTx Transaction) (Transaction, error) {

	// Return an error if the current transaction does not have a parent, i.e is root.
	if currentTx.parent == nil {
		return currentTx, errors.New("ERROR: COMMIT called with no active transaction.\n")
	}

	// Copy all operations from the current transaction to the parent transaction.
	for opKey, opValue := range currentTx.Operations {
		currentTx.parent.Operations[opKey] = opValue
	}

	// Return the parent as the current transaction.
	return *currentTx.parent, nil
}

// abortTx discards all operations in the current transaction.
func abortTx(currentTx Transaction) (Transaction, error) {

	// Return an error if the current transaction does not have a parent, i.e is root.
	if currentTx.parent == nil {
		return currentTx, errors.New("ERROR: ABORT called with no active transaction.\n")
	}

	// Set the parent as the current transaction causing all pending transactions to be discarded.
	return *currentTx.parent, nil
}
