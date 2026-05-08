package main

import (
	"bytes"
	"fmt"
	"iter"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	output, err := os.ReadFile("/proc/net/tcp")
	if err != nil {
		log.Fatal(err)
	}
	filtered := getTCPInodes(output)

	procs, err := scanProc(filtered)
	if err != nil {
		log.Fatal(err)
	}
	for i := range procs {
		fmt.Printf("Process has socket: pid=%d\n", procs[i].Pid)
	}
}

func getTCPInodes(content []byte) [][]byte {
	inodes := [][]byte{}
	next, stop := iter.Pull(bytes.Lines(content))
	defer stop()
	_, ok := next()
	if !ok {
		return inodes
	}
	for {
		line, ok := next()
		if !ok {
			break
		}
		fields := bytes.Fields(line)
		if !bytes.Equal(fields[3], []byte("0A")) {
			continue
		}
		inodes = append(inodes, fields[9])
	}
	return inodes
}

func scanProc(inodes [][]byte) ([]*os.Process, error) {
	processes := []*os.Process{}
	files, err := os.ReadDir("/proc")
	if err != nil {
		// proc might not be mounted
		return []*os.Process{}, 
			fmt.Errorf("could not list \"/proc\": %w", err)
	}
	for _, file := range files {
		if file.IsDir() {
			pidStr := file.Name()
			isPid := true
			for i := 0; i < len(pidStr); i++ {
				if pidStr[i] < '0' || pidStr[i] > '9' {
					isPid = false
					break
				}
			}
			if !isPid {
				// directory is not a PID
				continue
			}
			currentDir := "/proc/" + pidStr + "/fd"
			fds, err := os.ReadDir(currentDir)
			if err != nil {
				// TODO: error type checking
				// most likely error cause: read perms -- maybe implement: 
				// log.Warning("no read permission for dir %s", currentDir)
				continue
			}
			if len(fds) == 0 {
				// no file descriptors for this process
				continue
			}
			for _, fd := range fds {
				linked, err := os.Readlink(currentDir + "/" + fd.Name())
				if err != nil {
					// Skip this fd due to error: this is unlikely to occur
					// err != nil would indicate non symlink in path
					// /proc/<pid>/fd/*
					continue
				}
				for j := range len(inodes) {
					if strings.Contains(
							linked, "["+string(inodes[j])+"]") {
						pidInt, err := strconv.Atoi(pidStr)
						if err != nil {
							// if this happens, something changed in the
							// kernel proc implementation
							continue
						}
						process, err := os.FindProcess(pidInt)
						if err != nil {
							// if this happens, something changed in the
							// kernel proc implementation
							continue
						}
						processes = append(processes, process)
						break
					}
				}
			}
		}
	}
	return processes, nil
}
