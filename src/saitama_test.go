package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestListCommand(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("Tests need to be run as root to access /proc")
	}

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

	// Verifica se processos comuns estão listados
	commonProcesses := []string{"systemd", "init"}
	found := false
	for _, proc := range commonProcesses {
		if strings.Contains(output, proc) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected to find at least one common process, found none")
	}
}

func TestPunchCommand(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("Tests need to be run as root to access /proc")
	}

	testCases := []struct {
		name     string
		force    bool
		expected string
	}{
		{
			name:     "Normal punch",
			force:    false,
			expected: "Killing sleep with one punch",
		},
		{
			name:     "Force punch",
			force:    true,
			expected: "Killing sleep with one punch",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a dummy process for testing
			cmd := exec.Command("sleep", "60")
			if err := cmd.Start(); err != nil {
				t.Fatalf("Expected no error starting process, got %v", err)
			}

			// Ensure the process is killed after test
			defer func() {
				if cmd.Process != nil {
					cmd.Process.Kill()
				}
			}()

			// Wait a bit to ensure process is running
			time.Sleep(100 * time.Millisecond)

			// Redirect stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Prepare command args
			args := []string{"punch", "sleep"}
			if tc.force {
				args = append(args, "--force")
			}
			rootCmd.SetArgs(args)

			// Execute the punch command
			if err := rootCmd.Execute(); err != nil {
				t.Fatalf("Expected no error executing command, got %v", err)
			}

			// Restore stdout and capture output
			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r); err != nil {
				t.Fatalf("Expected no error reading output, got %v", err)
			}

			output := buf.String()
			if !strings.Contains(output, tc.expected) {
				t.Errorf("Expected '%s' in output, got %v", tc.expected, output)
			}

			// Verify process was actually killed
			if err := cmd.Process.Signal(0); err == nil {
				t.Error("Process should have been killed")
			}
		})
	}
}

func TestHandleProcessInvalidPath(t *testing.T) {
	err := handleProcess("/proc/invalid/status", true)
	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}

func TestHandleProcessInvalidFormat(t *testing.T) {
	// Create temporary file with invalid format
	tmpfile, err := os.CreateTemp("", "test-status")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte("invalid format")); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	err = handleProcess(tmpfile.Name(), true)
	if err == nil {
		t.Error("Expected error for invalid format, got nil")
	}
}
