package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

const asciiArt = `
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⣠⣶⡾⠏⠉⠙⠳⢦⡀⠀⠀⠀⢠⠞⠉⠙⠲⡀⠀
⠀⠀⠀⣴⠿⠏⠀⠀⠀⠀⠀⠀⢳⡀⠀⡏⠀⠀⠀⠀⢷
⠀⠀⢠⣟⣋⡀⢀⣀⣀⡀⠀⣀⡀⣧⠀⢸⠀⠀⠀⠀⠀⡇
⠀⠀⢸⣯⡭⠁⠸⣛⣟⠆⡴⣻⡲⣿⠀⣸⠀⠀Oh!⠀⡇
⠀⠀⣟⣿⡭⠀⠀⠀⠀⠀⢱⠀⠀⣿⠀⢹⠀⠀⠀⠀⠀⡇
⠀⠀⠙⢿⣯⠄⠀⠀⠀⢀⡀⠀⠀⡿⠀⠀⡇⠀⠀⠀⠀⡼
⠀⠀⠀⠀⠹⣶⠆⠀⠀⠀⠀⠀⡴⠃⠀⠀⠘⠤⣄⣠⠞⠀
⠀⠀⠀⠀⠀⢸⣷⡦⢤⡤⢤⣞⣁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⢀⣤⣴⣿⣏⠁⠀⠀⠸⣏⢯⣷⣖⣦⡀⠀⠀⠀⠀⠀⠀
⢀⣾⣽⣿⣿⣿⣿⠛⢲⣶⣾⢉⡷⣿⣿⠵⣿⠀⠀⠀⠀⠀⠀
⣼⣿⠍⠉⣿⡭⠉⠙⢺⣇⣼⡏⠀⠀⠀⣄⢸⠀⠀⠀⠀⠀⠀
⣿⣿⣧⣀⣿.........⣀⣰⣏⣘⣆
`

type Process struct {
	PID  int
	Name string
}


type ProcessManager struct{}

func NewProcessManager() *ProcessManager {
	return &ProcessManager{}
}

var rootCmd = &cobra.Command{
	Use:   "saitama",
	Short: "Saitama is a tool to manage processes",
	Long: `Saitama is a command line tool to list and punch processes by name.
It uses the /proc filesystem to gather process information and allows you to punch processes with one command.`,
	Run: func(cmd *cobra.Command, args []string) {
		if runtime.GOOS != "linux" {
			fmt.Printf("This tool only works on Linux systems (current: %s)\n", runtime.GOOS)
			os.Exit(1)
		}
		cmd.Help()
	},
}

var listCmd = &cobra.Command{
	Use:   "list [pattern]",
	Short: "List all processes or filter by pattern",
	Long:  "List all running processes. Optionally filter by process name pattern.",
	Run: func(cmd *cobra.Command, args []string) {
		pm := NewProcessManager()
		
		var pattern string
		if len(args) > 0 {
			pattern = args[0]
		}
		
		processes, err := pm.ListProcesses(pattern)
		if err != nil {
			log.Fatalf("Error listing processes: %v", err)
		}
		
		if len(processes) == 0 {
			if pattern != "" {
				fmt.Printf("No processes found matching pattern: %s\n", pattern)
			} else {
				fmt.Println("No processes found")
			}
			return
		}
		
		fmt.Printf("%-8s %s\n", "PID", "NAME")
		fmt.Println(strings.Repeat("-", 40))
		for _, proc := range processes {
			fmt.Printf("%-8d %s\n", proc.PID, proc.Name)
		}
		fmt.Printf("\nTotal: %d processes\n", len(processes))
	},
}

var punchCmd = &cobra.Command{
	Use:   "punch <processname>",
	Short: "Punch (kill) processes by name",
	Long:  "Kill all processes that match the given name. Use --force for SIGKILL or --dry-run to see what would be killed.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		processName := args[0]
		force, _ := cmd.Flags().GetBool("force")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		
		pm := NewProcessManager()
		processes, err := pm.ListProcesses(processName)
		if err != nil {
			log.Fatalf("Error finding processes: %v", err)
		}
		
		if len(processes) == 0 {
			fmt.Printf("No processes found with name: %s\n", processName)
			return
		}
		
		var matchingProcesses []Process
		for _, proc := range processes {
			if proc.Name == processName {
				matchingProcesses = append(matchingProcesses, proc)
			}
		}
		
		if len(matchingProcesses) == 0 {
			fmt.Printf("No processes found with exact name: %s\n", processName)
			return
		}
		
		if dryRun {
			fmt.Printf("Would kill %d process(es) with name '%s':\n", len(matchingProcesses), processName)
			for _, proc := range matchingProcesses {
				signal := "SIGTERM"
				if force {
					signal = "SIGKILL"
				}
				fmt.Printf("  PID: %d, Signal: %s\n", proc.PID, signal)
			}
			return
		}
		
		if len(matchingProcesses) > 1 {
			fmt.Printf("Found %d processes with name '%s'. Continue? (y/N): ", len(matchingProcesses), processName)
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
				fmt.Println("Operation cancelled.")
				return
			}
		}
		
		killedCount := 0
		for _, proc := range matchingProcesses {
			if err := pm.KillProcess(proc.PID, force); err != nil {
				fmt.Printf("Failed to kill process %d (%s): %v\n", proc.PID, proc.Name, err)
			} else {
				killedCount++
				signal := "SIGTERM"
				if force {
					signal = "SIGKILL"
				}
				fmt.Printf("Successfully killed process %d (%s) with %s\n", proc.PID, proc.Name, signal)
			}
		}
		
		if killedCount > 0 {
			fmt.Printf("\n%s\n", asciiArt)
			fmt.Printf("One punch! Killed %d process(es) with name '%s'\n", killedCount, processName)
		}
	},
}

func main() {
	punchCmd.Flags().BoolP("force", "f", false, "Use SIGKILL instead of SIGTERM")
	punchCmd.Flags().BoolP("dry-run", "d", false, "Show what would be killed without actually killing")
	
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(punchCmd)
	
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func (pm *ProcessManager) ListProcesses(pattern string) ([]Process, error) {
	var processes []Process
	
	procDirs, err := filepath.Glob("/proc/[0-9]*")
	if err != nil {
		return nil, fmt.Errorf("failed to read /proc directory: %w", err)
	}
	
	for _, dir := range procDirs {
		pidStr := filepath.Base(dir)
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			continue // Skip non-numeric directories
		}
		
		processName, err := pm.getProcessName(pid)
		if err != nil {
			continue // Skip processes we can't read
		}
		
		if pattern != "" {
			if !strings.Contains(strings.ToLower(processName), strings.ToLower(pattern)) {
				continue
			}
		}
		
		processes = append(processes, Process{
			PID:  pid,
			Name: processName,
		})
	}
	
	return processes, nil
}

func (pm *ProcessManager) getProcessName(pid int) (string, error) {
	commPath := fmt.Sprintf("/proc/%d/comm", pid)
	data, err := os.ReadFile(commPath)
	if err != nil {
		return pm.getProcessNameFromCmdline(pid)
	}
	
	name := strings.TrimSpace(string(data))
	if name == "" {
		return pm.getProcessNameFromCmdline(pid)
	}
	
	return name, nil
}


func (pm *ProcessManager) getProcessNameFromCmdline(pid int) (string, error) {
	cmdlinePath := fmt.Sprintf("/proc/%d/cmdline", pid)
	data, err := os.ReadFile(cmdlinePath)
	if err != nil {
		return "", fmt.Errorf("failed to read cmdline for PID %d: %w", pid, err)
	}
	
	if len(data) == 0 {
		return fmt.Sprintf("[pid_%d]", pid), nil 
	}
	
	cmdParts := strings.Split(string(data), "\x00")
	if len(cmdParts) > 0 && cmdParts[0] != "" {
		return filepath.Base(cmdParts[0]), nil
	}
	
	return fmt.Sprintf("[pid_%d]", pid), nil
}

func (pm *ProcessManager) KillProcess(pid int, force bool) error {
	if !pm.processExists(pid) {
		return fmt.Errorf("process %d does not exist", pid)
	}
	
	proc, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process %d: %w", pid, err)
	}
	
	if force {
		return proc.Kill() // SIGKILL
	}
	return proc.Signal(syscall.SIGTERM) // SIGTERM
}

func (pm *ProcessManager) processExists(pid int) bool {
	_, err := os.Stat(fmt.Sprintf("/proc/%d", pid))
	return err == nil
}