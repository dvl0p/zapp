package main

import (
	"os"
	"strconv"
)

type Process struct {
	Name string
	PID int
}

func (p Process) getPath() string {
	pidStr := strconv.Itoa(p.PID)
	return "/proc/" + pidStr
}

func newProcess(pid int, name string) Process {
	return Process{PID: pid, Name: name}
}

func getProcesses() ([]Process, error) {
	
	processes := []Process{}

	procDirs, err := os.ReadDir("/proc")
	if err != nil {
		return []Process{}, err
	}
	for i := range procDirs {
		if procDirs[i].IsDir() {
			pidStr := procDirs[i].Name()
			pid, err := strconv.Atoi(pidStr)
			if err != nil {
				continue
			}
			name, err := os.ReadFile("/proc/" + pidStr + "/comm")
			if err != nil {
				processes = append(processes, newProcess(pid, ""))
			} else {
				processes = append(processes, newProcess(pid, string(name)))
			}
		}
	}

	return processes, nil
}
