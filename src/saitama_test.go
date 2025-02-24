package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func hasCommonProcess(output string) bool {
	commonProcesses := []string{"systemd", "init"}
	for _, proc := range commonProcesses {
		if strings.Contains(output, proc) {
			return true
		}
	}
	return false
}

func TestListCommand(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("Tests need to be run as root to access /proc")
	}

	output := captureOutput(t, func() {
		if err := manager.ExecuteWithArgs([]string{"list"}); err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 {
		t.Error("Expected at least one process, got none")
	}
}

func setupTestProcess(t *testing.T) *exec.Cmd {
	cmd := exec.Command("sleep", "60")
	if err := cmd.Start(); err != nil {
		t.Fatalf("Expected no error starting process, got %v", err)
	}
	time.Sleep(500 * time.Millisecond)
	return cmd
}

func captureOutput(t *testing.T, fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	
	fn()
	
	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("Expected no error reading output, got %v", err)
	}
	return buf.String()
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
		{name: "Normal punch", force: false, expected: "Killing sleep with one punch"},
		{name: "Force punch", force: true, expected: "Killing sleep with one punch"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("Running test case: %s\n", tc.name)
			done := make(chan error, 1)
			cmd := setupTestProcess(t)
			defer cmd.Process.Kill()

			go func() {
				done <- cmd.Wait()
			}()

			args := []string{"punch", "sleep"}
			if tc.force {
				args = append(args, "--force")
			}

			output := captureOutput(t, func() {
				if err := manager.ExecuteWithArgs(args); err != nil {
					t.Fatalf("Expected no error executing command, got %v", err)
				}
			})

			if !strings.Contains(output, tc.expected) {
				t.Errorf("Expected '%s' in output, got %v", tc.expected, output)
			}

			select {
			case err := <-done:
				if err == nil {
					t.Error("Process should have been killed")
				}
			case <-time.After(5 * time.Second):
				t.Error("Timeout waiting for process to be killed")
			}
		})
	}
}

func TestHandleProcessInvalidPath(t *testing.T) {
	pm := NewProcessManager()
	err := pm.handleProcess("/proc/invalid/status", true, "", nil)
	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}

func TestHandleProcessInvalidFormat(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test-status")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte("invalid format")); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	pm := NewProcessManager()
	err = pm.handleProcess(tmpfile.Name(), true, "", nil)
	if err == nil {
		t.Error("Expected error for invalid format, got nil")
	}
}
