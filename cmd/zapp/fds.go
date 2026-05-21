package main

import (
	"errors"
	"strconv"
	"os"
	"sync"
	"strings"
)

type ProcessSocket struct {
	Process Process
	Socket Socket
}

func matchInodes(processes []Process, sockets []Socket, 
		workers int) []ProcessSocket {

	worker := func(inodeSet map[Inode]Socket, jobs <-chan Process,
			results chan<- ProcessSocket) {

		for process := range jobs {
			fdPathRoot := getFDPath(process)
			fds, err := os.ReadDir(fdPathRoot)
			if err != nil {
				continue
			}

			for _, fd := range fds {
				inode, err := getInode(fdPathRoot + fd.Name())
				if err != nil {
					continue
				}

				if socket, ok := inodeSet[inode]; ok {
					results <- ProcessSocket{
						Process: process,
						Socket: socket,
					}
				}
			}
		}
	}

	inodeSet := make(map[Inode]Socket, len(sockets))
	for _, s := range sockets {
		inodeSet[s.Inode] = s
	}

	workers = max(workers, 1) 
	jobs := make(chan Process)
	results := make(chan ProcessSocket)
	var wg sync.WaitGroup

	go func() {
		defer close(jobs)
		for _, process := range processes {
			jobs <- process
		}
	}()

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(inodeSet, jobs, results)
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	processSockets := []ProcessSocket{}
	for r := range results {
		processSockets = append(processSockets, r)
	}

	return processSockets
}

func getFDPath(process Process) string {
	return process.getPath() + "/fd/"
}

func getInode(fdPath string) (Inode, error) {
	fdStr, err := os.Readlink(fdPath)
	if err != nil {
		return 0, err
	}
	part, found := strings.CutPrefix(fdStr, "socket:")
	if !found {
		return 0, errors.New("wrong file descriptor type")
	}
	inodeStr := strings.Trim(part, "[]")
	inodeInt, err := strconv.Atoi(inodeStr)
	if err != nil {
		return 0, err
	}
	return Inode(inodeInt), nil
}
