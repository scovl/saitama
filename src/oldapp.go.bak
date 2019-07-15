/*
Saitama a way to kill your processes with one punch.
developed by lobocode - lobocode@fedoraproject.org
*/

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

// findAndKillProcess in /proc/<pid>/status name and number
func findAndKillProcess(path string, info os.FileInfo, err error) error {

	//if err != nil {
	//    //fmt.Printf("Process owner is root!\n Please use sudo.\n")
	//    fmt.Printf("%s\n",err)
	//    return nil
	// }

	if strings.Count(path, "/") == 3 {
		if strings.Contains(path, "/status") {

			// Let's extract the middle part of the path with the <pid> and
			// convert the <pid> into an integer. Log an error if it fails.
			pid, _ := strconv.Atoi(path[6:strings.LastIndex(path, "/")])

			f, _ := ioutil.ReadFile(path)

			// Extract the process name from within the first line in the buffer
			name := string(f[6:bytes.IndexByte(f, '\n')])

			asciArt :=
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

			if name == args[1] {

				proc, _ := os.FindProcess(pid)

				if proc.Kill() != nil {
					// If process owner is root
					fmt.Printf("Please use sudo\n")
				} else {
					// Kill the process
					fmt.Printf("Killing %s with one punch \n", args[1])
					fmt.Printf("PID: %d %s %s .\n", pid, name, asciArt)
					proc.Kill()
				}
			}

			if "-l" == args[1] {
				fmt.Printf("%s\n", name)
			}
		}
	}

	return nil
}

// main is the entry point of any go application
func main() {
	args = os.Args

	if len(args) != 2 {
		log.Fatalln("Usage: saitama -p <processname> to kill with one punch \n or saitama -l to list process names ")
	}

	err := filepath.Walk("/proc", findAndKillProcess)

	if err != nil {
		if err == io.EOF {
			// Not an error, just a signal when we are done
			err = nil
		} else {
			log.Fatal(err)
		}
	}
}
