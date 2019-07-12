/* 
Saitama a way to kill your processes with a single punch. 
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

func findAndKillProcess(path string, info os.FileInfo, err error) error {

    // We are only interested in files with a path looking like /proc/<pid>/status.
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
                fmt.Printf("PID: %d %s %s .\n", pid, name, asciArt)

                proc, _ := os.FindProcess(pid)

                // Kill the process
                proc.Kill()

            }

        }
    }

    return nil
}

// main is the entry point of any go application
func main() {
    args = os.Args

    if len(args) != 2 {
        log.Fatalln("Usage: saitama <processname>")
    }

    fmt.Printf("Killing %s with one punch \n", args[1])

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