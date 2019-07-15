package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var args []string

func findAndKillByName(path string, info os.FileInfo, err error) error {

	if strings.Count(path, "/") == 3 {
		if strings.Contains(path, "/status") {
			//pid, _ := strconv.Atoi(path[6:strings.LastIndex(path, "/")])

			statesIndex, _ := ioutil.ReadFile(path)

			processName := string(statesIndex[6:bytes.IndexByte(statesIndex, '\n')])
			//threadsNum := string(statesIndex[])

			//fmt.Printf("%s", pid)
			//fmt.Printf("%s", f)
			fmt.Println(processName)
		}
	}

	return nil
}

func main() {

	args = os.Args

	debugFunc := filepath.Walk("/proc", findAndKillByName)

	fmt.Println(debugFunc)

	/* 	err := filepath.Walk("/proc", findAndKillByName)

	   	if err != nil {
	   		if err == io.EOF {
	   			// Not an error, just a signal when we are done
	   			err = nil
	   		} else {
	   			log.Fatal(err)
	   		}
	   	} */
}
