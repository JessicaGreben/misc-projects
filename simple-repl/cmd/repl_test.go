package cmd

import "testing"

func TestParseArgs(t *testing.T) {
	cases := []struct {
		args  string
		cmd   string
		key   string
		value string
	}{
		{
			"",
			"",
			"",
			"",
		},
		{
			" WRITE a hello\n",
			"WRITE",
			"a",
			"hello",
		},
		{
			"QUIT\n",
			"QUIT",
			"",
			"",
		},
		{
			"READ a \n",
			"READ",
			"a",
			"",
		},
	}

	for _, exp := range cases {
		cmd, key, value := parseArgs(exp.args)

		// Is the correct command parsed from the input?
		if cmd != exp.cmd {
			t.Fatalf("Failed.\nActual: %v.\nExpected: %v.\n",
				cmd,
				exp.cmd,
			)
		}

		// Is the correct key parsed from the input?
		if key != exp.key {
			t.Fatalf("Failed.\nActual: %v.\nExpected: %v.\n",
				key,
				exp.key,
			)
		}

		// Is the correct value parsed from the input?
		if value != exp.value {
			t.Fatalf("Failed.\nActual: %v.\nExpected: %v.\n",
				value,
				exp.value,
			)
		}
	}
}
