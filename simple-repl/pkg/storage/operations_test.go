package storage

import (
	"errors"
	"reflect"
	"testing"
)

func TestReadTx(t *testing.T) {

	// Case 1 setup: key exists in current transaction, but it has been deleted.
	deletedOp := map[string]string{"b": ""}

	currentTx1 := Transaction{
		Operations: deletedOp,
	}

	// Case 2 setup: key exists in the current transaction and it has not been deleted.
	currentTx2 := currentTx1
	currentTx2.Operations["a"] = "hello"

	// Case 4 setup: key does not exists in current transaction nor in parent transaction.
	currentTx3 := Transaction{
		parent: &currentTx1,
	}

	// Case 5 setup: key does not exists in current transaction, but it does exist
	// in parent, but it is deleted.
	currentTx5 := Transaction{
		parent: &currentTx1,
	}

	// Case 6 setup: key does not exists in current transaction, does exist in
	// parent and it is not deleted.
	currentTx6 := Transaction{
		parent: &currentTx2,
	}

	cases := []struct {
		inputKey string
		output   string
		inputTx  Transaction
	}{

		// case 1: key exists in current transaction, but it has been deleted.
		{
			"b",
			"",
			currentTx1,
		},

		// Case 2: key exists in the current transaction and it has not been deleted.
		{
			"a",
			"hello",
			currentTx2,
		},

		// Case 3: key does not exists in current transaction and there is no parent transaction.
		{
			"c",
			"",
			currentTx2,
		},

		// Case 4: key does not exists in current transaction nor in parent transaction.
		{
			"c",
			"",
			currentTx3,
		},

		// Case 5: key does not exists in current transaction, but it does exist in
		// parent, but it's deleted.
		{
			"b",
			"",
			currentTx5,
		},

		// Case 6: key does not exists in current transaction, but it does exist in
		// parent and it's not deleted.
		{
			"a",
			"hello",
			currentTx6,
		},
	}

	for _, exp := range cases {
		actualOutput, _ := readTx(exp.inputKey, exp.inputTx)

		// Is the correct value of the key read?
		if actualOutput != exp.output {
			t.Fatalf("Failed.\nActual: %v.\nExpected: %v.\n",
				actualOutput,
				exp.output,
			)
		}
	}
}

func TestWriteTx(t *testing.T) {

	// Case 1 setup: add a new operation to a transaction.
	currentTx1 := Transaction{
		Operations: map[string]string{},
	}

	// Case 2 setup: update an existing opertaion with a new value.
	currentTx2 := Transaction{
		Operations: map[string]string{"a": "hello"},
	}

	cases := []struct {
		key       string
		value     string
		currentTx Transaction
	}{

		// Case 1: add a new operation to a transaction.
		{
			key:       "a",
			value:     "hello",
			currentTx: currentTx1,
		},

		// Case 2: update an existing opertaion with a new value.
		{
			key:       "a",
			value:     "hello-again",
			currentTx: currentTx2,
		},
	}

	for _, exp := range cases {
		writeTx(exp.key, exp.value, exp.currentTx)

		// Does the correct value get written?
		if exp.currentTx.Operations[exp.key] != exp.value {
			t.Fatalf("Failed.\nActual: %v.\nExpected: %v.\n",
				exp.currentTx.Operations[exp.key],
				exp.value,
			)
		}
	}
}

func TestDeleteTx(t *testing.T) {

	// Case 1 Setup: there are no operations in the current transaction.
	currentTx1 := Transaction{
		Operations: map[string]string{},
	}

	// Case 2 Setup: there are is an operation in the current transaction with a
	// matching key that is to be delted.
	currentTx2 := Transaction{
		Operations: map[string]string{"a": "hello"},
	}

	cases := []struct {
		key             string
		currentTx       Transaction
		operationsCount int
	}{

		// Case 1: there are no operations in the current transaction.
		{
			key:             "a",
			currentTx:       currentTx1,
			operationsCount: 1,
		},

		// Case 2: there are is an operation in the current transaction with a
		// matching key that is to be delted.
		{
			key:             "a",
			currentTx:       currentTx2,
			operationsCount: 1,
		},

		// Case 3: there are is an operation in the current transaction with no
		// matching key to be delted.
		{
			key:             "b",
			currentTx:       currentTx2,
			operationsCount: 2,
		},
	}

	for _, exp := range cases {
		deleteTx(exp.key, exp.currentTx)

		// Does a delete operation get added to the operations
		// of the current transaction?
		if len(exp.currentTx.Operations) != exp.operationsCount {
			t.Fatalf("Failed.\nActual: %v.\nExpected: %v.\n",
				len(exp.currentTx.Operations),
				exp.operationsCount,
			)
		}

		// Is the delete operation recorded correctly?
		// i.e. the value of the deleted key is a string's zero value.
		if exp.currentTx.Operations[exp.key] != "" {
			t.Fatalf("Failed.\nActual: %v.\nExpected: %v.\n",
				exp.currentTx.Operations[exp.key],
				"",
			)
		}
	}
}

func TestStartTx(t *testing.T) {
	root := Transaction{
		Operations: map[string]string{"a": "hello"},
	}

	cases := []struct {
		currentTx Transaction
	}{
		{
			currentTx: root,
		},
	}

	for _, exp := range cases {
		actualTx := startTx(exp.currentTx)

		// Are there no operations stored in the new child transactions?
		if len(actualTx.Operations) != 0 {
			t.Fatalf("Failed.\nActual: %v.\nExpected: 0.",
				len(actualTx.Operations),
			)
		}

		// Is the correct transaction set as the parent?
		if !reflect.DeepEqual(actualTx.parent, &exp.currentTx) {
			t.Fatalf("Failed.\nActual: %v.\nExpected: %v.\n",
				actualTx.parent,
				&exp.currentTx,
			)
		}
	}
}

func TestCommitTx(t *testing.T) {

	// Case 1 setup: The current transaction does not have a parent and therefore
	// it is the root store.
	var root Transaction

	// Case 2 setup: The current transaction has a parent and is therefore not the root store.
	currentTx2 := Transaction{
		parent: &root,
	}

	cases := []struct {
		currentTx  Transaction
		expectedTx Transaction
		err        error
	}{

		// Case 1: The current transaction does not have a parent and therefor
		// is the root store.
		{
			currentTx:  root,
			expectedTx: root,
			err:        errors.New("ERROR: COMMIT called with no active transaction.\n"),
		},

		// Case 2: The current transaction has a parent and is therefore not the root store.
		{
			currentTx:  currentTx2,
			expectedTx: root,
			err:        nil,
		},
	}

	for _, exp := range cases {
		actualTx, actualErr := commitTx(exp.currentTx)

		// Is the correct current transaction returned?
		if !reflect.DeepEqual(actualTx, exp.expectedTx) {
			t.Fatalf("Failed.\nActual: %v.\nExpected: %v.\n",
				actualTx,
				exp.expectedTx,
			)
		}

		// Is the correct error returned?
		if actualErr != nil {
			if actualErr.Error() != exp.err.Error() {
				t.Fatalf("Failed.\nActual: %v.\nExpected: %v.\n",
					actualErr.Error(),
					exp.err.Error(),
				)
			}
		}
	}
}

func TestAbortTxs(t *testing.T) {

	// Case 1 setup: the current transaction has no parent.
	var root Transaction

	// Case 2 setup: the current transaction does have a parent.
	currentTx := Transaction{
		parent: &root,
	}

	cases := []struct {
		currentTx  Transaction
		expectedTx Transaction
		err        error
	}{

		// Case 1: the current transaction has no parent.
		{
			currentTx:  root,
			expectedTx: root,
			err:        errors.New("ERROR: ABORT called with no active transaction.\n"),
		},

		// Case 2: the current transaction does have a parent.
		{
			currentTx:  currentTx,
			expectedTx: root,
			err:        nil,
		},
	}

	for _, exp := range cases {
		actualTx, actualErr := abortTx(exp.currentTx)

		// Is the correct error returned?
		if actualErr != nil {
			if actualErr.Error() != exp.err.Error() {
				t.Fatalf("Failed.\nActual: %v.\nExpected: %v.\n",
					actualErr.Error(),
					exp.err.Error(),
				)
			}
		}

		// Is the current transaction now the previous parent transaction?
		if !reflect.DeepEqual(actualTx, exp.expectedTx) {
			t.Fatalf("Failed.\nActual: %v.\nExpected: %v.\n",
				actualTx,
				exp.expectedTx,
			)
		}
	}
}
