package main

import "log"

const (
	workers = 4
)

func run() error {
	sockets, err := discoverSockets()
	if err != nil {
		return err
	}
	log.Printf("Sockets found: %+v", sockets)

	processes, err := getProcesses()
	if err != nil {
		return err
	}

	matched := matchInodes(processes, sockets, workers)
	log.Printf("Matching processes: %+v", matched)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

