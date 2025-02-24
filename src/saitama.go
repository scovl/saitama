package main

import (
	"bytes"
	"fmt"
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

type ProcessManager struct {
	rootCmd  *cobra.Command
	listCmd  *cobra.Command
	punchCmd *cobra.Command
}

func NewProcessManager() *ProcessManager {
	pm := &ProcessManager{}
	
	pm.rootCmd = &cobra.Command{
		Use:   "saitama",
		Short: "Saitama is a tool to manage processes",
		Long:  `Saitama is a command line tool to list and punch processes by name.`,
	}

	pm.listCmd = &cobra.Command{
		Use:   "list",
		Short: "List processes by name",
		RunE:  pm.runList,
	}

	pm.punchCmd = &cobra.Command{
		Use:   "punch [processname]",
		Short: "Punch process by name",
		Args:  cobra.ExactArgs(1),
		RunE:  pm.runPunch,
	}

	pm.punchCmd.Flags().BoolP("force", "f", false, "Use SIGKILL instead of SIGTERM")
	pm.rootCmd.AddCommand(pm.listCmd, pm.punchCmd)
	
	return pm
}

func (pm *ProcessManager) runList(cmd *cobra.Command, args []string) error {
	return pm.walkProc(true, "")
}

func (pm *ProcessManager) runPunch(cmd *cobra.Command, args []string) error {
	force, _ := cmd.Flags().GetBool("force")
	return pm.walkProc(false, args[0], force)
}

func (pm *ProcessManager) walkProc(list bool, processName string, args ...interface{}) error {
	return filepath.Walk("/proc", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsNotExist(err) || strings.Contains(err.Error(), "no such") {
				return nil
			}
			return err
		}
		
		if strings.Count(path, "/") == 3 && strings.Contains(path, "/status") {
			return pm.handleProcess(path, list, processName, args...)
		}
		return nil
	})
}

func (pm *ProcessManager) handleProcess(path string, list bool, processName string, args ...interface{}) error {
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash < 6 {
		return fmt.Errorf("invalid path format: %s", path)
	}

	pid, err := strconv.Atoi(path[6:lastSlash])
	if err != nil {
		return fmt.Errorf("error converting PID to int: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) || strings.Contains(err.Error(), "no such process") {
			return nil // Ignora processos que não existem mais
		}
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
	
	pName := string(data[nameIndex+6 : nameIndex+endIndex])
	pName = strings.TrimSpace(pName)

	if list {
		fmt.Println(pName)
	} else if len(args) > 0 {
		if processName == pName {
			if len(args) > 0 {
				force, ok := args[0].(bool)
				if !ok {
					return fmt.Errorf("invalid force argument type")
				}
				
				if err := killProcess(pid, force); err != nil {
					return fmt.Errorf("error killing process: %v", err)
				}
				fmt.Printf("Killing %s with one punch\nPID: %d %s %s .\n", pName, pid, pName, asciiArt)
			}
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

func (pm *ProcessManager) Execute() error {
	return pm.rootCmd.Execute()
}

// Para testes
func (pm *ProcessManager) ExecuteWithArgs(args []string) error {
	pm.rootCmd.SetArgs(args)
	return pm.rootCmd.Execute()
}

var manager = NewProcessManager()

func main() {
	if err := manager.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
