package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

var rootCmd = &cobra.Command{
	Use:   "saitama",
	Short: "Saitama is a tool to manage processes",
	Long: `Saitama is a command line tool to list and punch processes by name.
It uses the /proc filesystem to gather process information and allows you to punch processes with one command.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List processes by name",
	Run: func(cmd *cobra.Command, args []string) {
		err := filepath.Walk("/proc", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.Count(path, "/") == 3 && strings.Contains(path, "/status") {
				return handleProcess(path, true)
			}
			return nil
		})
		if err != nil {
			log.Fatalf("Error walking through /proc: %v", err)
		}
	},
}

var punchCmd = &cobra.Command{
	Use:   "punch [processname]",
	Short: "Punch process by name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		processName := args[0]
		force, _ := cmd.Flags().GetBool("force")
		
		err := filepath.Walk("/proc", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.Count(path, "/") == 3 && strings.Contains(path, "/status") {
				return handleProcess(path, false, processName, force)
			}
			return nil
		})
		if err != nil {
			log.Fatalf("Error walking through /proc: %v", err)
		}
	},
}

func main() {
	punchCmd.Flags().BoolP("force", "f", false, "Use SIGKILL instead of SIGTERM")
	
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(punchCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handleProcess(path string, list bool, args ...interface{}) error {
	pid, err := strconv.Atoi(path[6:strings.LastIndex(path, "/")])
	if err != nil {
		return fmt.Errorf("error converting PID to int: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	if len(data) < 7 {
		return fmt.Errorf("invalid process status file: %s", path)
	}

	nameIndex := bytes.Index(data, []byte("Name:\t"))
	if nameIndex == -1 {
		return fmt.Errorf("could not find process name in: %s", path)
	}
	
	endIndex := bytes.Index(data[nameIndex:], []byte("\n"))
	if endIndex == -1 {
		return fmt.Errorf("invalid status file format: %s", path)
	}
	
	processName := string(data[nameIndex+6 : nameIndex+endIndex])
	processName = strings.TrimSpace(processName)

	if list {
		fmt.Println(processName)
	} else if len(args) > 0 {
		targetName, ok := args[0].(string)
		if !ok {
			return fmt.Errorf("invalid process name argument type")
		}
		
		if targetName == processName {
			force := false
			if len(args) > 1 {
				force, ok = args[1].(bool)
				if !ok {
					return fmt.Errorf("invalid force argument type")
				}
			}
			
			if err := killProcess(pid, force); err != nil {
				return fmt.Errorf("error killing process: %v", err)
			}
			fmt.Printf("Killing %s with one punch\nPID: %d %s %s .\n", processName, pid, processName, asciiArt)
		}
	}
	return nil
}

func killProcess(pid int, force bool) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	
	if force {
		return proc.Kill() // SIGKILL
	}
	return proc.Signal(syscall.SIGTERM) // Mais gentil
}
