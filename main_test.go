package main

import (
	"bytes"
	"os"
	"testing"
)

func TestRootCommand(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() {
		os.Stdout = oldStdout
	}()

	rootCmd.Run(rootCmd, []string{})

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)

	expectedParts := []string{
		"Welcome to CSS Color Variable Creator!",
		"Use --help to see available commands",
	}

	for _, part := range expectedParts {
		if !bytes.Contains(buf.Bytes(), []byte(part)) {
			t.Errorf("Expected output to contain %q", part)
		}
	}
}

func TestMain(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"css-color-variable-creator", "--help"}
	main()
}
