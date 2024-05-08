package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

func main() {
	helpFlag := flag.Bool("help", false, "Display this help and exit")
	listFlag := flag.Bool("list", false, "List process by name")
	punchFlag := flag.String("punch", "", "Punch process by name")

	flag.Parse()

	if *helpFlag {
		fmt.Println("Usage: saitama [OPTION] <processname>\n\nOptions:\n-h, --help\tDisplay this help and exit\n-l, --list\tList process by name\n-p, --punch\tPunch process by name")
		return
	}

	if *listFlag && *punchFlag != "" {
		log.Fatal("Conflicting options: use either --list or --punch")
	}

	err := filepath.Walk("/proc", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.Count(path, "/") == 3 && strings.Contains(path, "/status") {
			return handleProcess(path)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking through /proc: %v", err)
	}
}

func handleProcess(path string) error {
	pid, err := strconv.Atoi(path[6:strings.LastIndex(path, "/")])
	if err != nil {
		return fmt.Errorf("error converting PID to int: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	processName := string(data[6:bytes.IndexByte(data, '\n')])

	switch {
	case flag.Lookup("list").Value.(flag.Getter).Get().(bool):
		fmt.Println(processName)
	case flag.Lookup("punch").Value.(flag.Getter).Get().(string) == processName:
		if err := killProcess(pid); err != nil {
			return fmt.Errorf("error killing process: %v", err)
		}
		fmt.Printf("Killing %s with one punch\nPID: %d %s %s .\n", processName, pid, processName, asciiArt)
	}
	return nil
}

func killProcess(pid int) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return proc.Kill()
}
