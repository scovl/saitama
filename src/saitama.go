// killprocess project main.go
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// findAndKillProcess searches the /proc directory to find the process with the specified name
// on the command line and sends an interrupt signal to kill it.
func findAndKillProcess(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	// It only processes directories with names similar to "/proc/<pid>/status".
	if strings.Count(path, "/") == 3 && strings.Contains(path, "/status") {
		pidStr := path[strings.Index(path, "proc")+5 : strings.Index(path, "/status")]
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			log.Printf("Could not convert pid %s for integer: %v", pidStr, err)
			return nil
		}

		// Extracts the process name from the first line of the "status" file.
		status, err := ioutil.ReadFile(path)
		if err != nil {
			log.Printf("Could not read the file %s: %v", path, err)
			return nil
		}
		processName := ""
		scanner := bufio.NewScanner(strings.NewReader(string(status)))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "Name:") {
				processName = strings.TrimSpace(strings.TrimPrefix(line, "Name:"))
				break
			}
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Error reading file %s: %v", path, err)
			return nil
		}

		// Performs the action as per the command line.
		switch args[1] {
		case "--help", "-h":
			fmt.Println(`Usage: saitama [OPTION] <processname>
Kill process with one punch

Mandatory arguments:

-h, --help		display this help and exit
-l, --list		list process by name
-p, --punch <processname>	punch process by name`)
			return io.EOF

		case "-l", "--list":
			fmt.Printf("%s\n", processName)

		case "-p", "--punch":
			if len(args) < 3 {
				log.Fatalln("\nMissing operand\nTry 'saitama --help' for more information")
			}
			if processName == args[2] {
				proc, err := os.FindProcess(pid)
				if err != nil {
					log.Printf("Could not find process %d: %v", pid, err)
					return nil
				}
				if err := proc.Kill(); err != nil {
					fmt.Printf("\nWarning: This process owner is 'root'\nPlease use 'sudo'\n")
				} else {
					fmt.Printf("Killing %s with one punch \n", args[2])
					fmt.Printf("PID: %d %s\n", pid, oh)
				}
			}

		default:
			log.Fatalln("\nMissing operand\nTry 'saitama --help' for more information")
