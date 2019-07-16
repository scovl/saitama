// killprocess project main.go
package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// args holds the commandline args
var args []string

func findAndKillProcess(path string, info os.FileInfo, err error) error {

	// We are only interested in files with a path looking like /proc/<pid>/status.
	if strings.Count(path, "/") == 3 {
		if strings.Contains(path, "/status") {

			pid, _ := strconv.Atoi(path[6:strings.LastIndex(path, "/")])

			readIndex, _ := ioutil.ReadFile(path)

			// Extract the process name from within the first line in the buffer
			processName := string(readIndex[6:bytes.IndexByte(readIndex, '\n')])

			oh :=
				`⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
		⠀⠀⠀⠀⣠⣶⡾⠏⠉⠙⠳⢦⡀⠀⠀⠀⢠⠞⠉⠙⠲⡀⠀
		⠀⠀⠀⣴⠿⠏⠀⠀⠀⠀⠀⠀⢳⡀⠀ ⡏⠀⠀⠀⠀ ⢷
		⠀⠀⢠⣟⣋⡀⢀⣀⣀⡀⠀⣀⡀⣧⠀⢸⠀⠀⠀⠀⠀ ⡇
		⠀⠀⢸⣯⡭⠁⠸⣛⣟⠆⡴⣻⡲⣿⠀⣸⠀⠀Oh!⠀⡇
		⠀⠀⣟⣿⡭⠀⠀⠀⠀⠀⢱⠀⠀⣿⠀⢹⠀⠀⠀⠀⠀ ⡇
		⠀⠀⠙⢿⣯⠄⠀⠀⠀⢀⡀⠀⠀⡿⠀⠀⡇⠀⠀⠀⠀⡼
		⠀⠀⠀⠀⠹⣶⠆⠀⠀⠀⠀⠀⡴⠃⠀⠀⠘⠤⣄⣠⠞⠀
		⠀⠀⠀⠀⠀⢸⣷⡦⢤⡤⢤⣞⣁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
		⠀⠀⢀⣤⣴⣿⣏⠁⠀⠀⠸⣏⢯⣷⣖⣦⡀⠀⠀⠀⠀⠀⠀
		⢀⣾⣽⣿⣿⣿⣿⠛⢲⣶⣾⢉⡷⣿⣿⠵⣿⠀⠀⠀⠀⠀⠀
		⣼⣿⠍⠉⣿⡭⠉⠙⢺⣇⣼⡏⠀⠀⠀⣄⢸⠀⠀⠀⠀⠀⠀
		⣿⣿⣧⣀⣿.........⣀⣰⣏⣘⣆ 
		`

			switch args[1] {
			case "--help", "-h":
				log.Fatalln("\nUsage: saitama [OPTION] <processname>\nKill process with one punch\n\nMandatory arguments\n\n-h  --help	display this help and exit\n-l  --list	list process by name\n-p  --punch <processname> punch process by name")
			case "-l", "--list":
				fmt.Printf("%s\n", processName)
			case "-p", "--punch":
				if len(args) == 2 {
					log.Fatalln("\nMissing operand\nTry 'saitama --help' for more information")
				} else if processName == args[2] {

					proc, _ := os.FindProcess(pid)

					if proc.Kill() != nil {
						fmt.Printf("\nWarning: This process owner is 'root'\nPlease use 'sudo'\n")
					} else {
						fmt.Printf("Killing %s with one punch \n", args[2])
						fmt.Printf("PID: %d %s %s .\n", pid, processName, oh)
						proc.Kill()
						

						// if error
						return io.EOF
					}
				}

			default:
				log.Fatalln("\nMissing operand\nTry 'saitama --help' for more information")
			}

		}
	}

	return nil
}

// main is the entry point of any go application
func main() {
	args = os.Args

	if len(args) == 1 {
		log.Fatalln("\nSaitama: missing operand\nTry 'saitama --help' for more information")
	}

	debugCall := filepath.Walk("/proc", findAndKillProcess)
	fmt.Printf("%s", debugCall)

}
