package main

import (
	"bytes"
	"os"
	"testing"
)

func TestListCommand(t *testing.T) {
	// Redirect stdout to capture the output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute the list command
	rootCmd.SetArgs([]string{"list"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Restore stdout and capture the output
	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Errorf("Expected output, got empty string")
	}
}

func TestPunchCommand(t *testing.T) {
	// Create a dummy process for testing
	cmd := exec.Command("sleep", "60")
	if err := cmd.Start(); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	defer cmd.Process.Kill() // Ensure the process is killed after test

	// Redirect stdout to capture the output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute the punch command
	rootCmd.SetArgs([]string{"punch", "sleep"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Restore stdout and capture the output
	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Killing sleep with one punch") {
		t.Errorf("Expected 'Killing sleep with one punch' in output, got %v", output)
	}
}
