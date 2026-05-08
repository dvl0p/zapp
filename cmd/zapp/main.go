package main


import (
	"fmt"
	"os"
	"log"
	"bytes"
	"iter"
)


func main() {
	output, err := os.ReadFile("/proc/net/tcp")
	if err != nil {
		log.Fatal(err)
	}
	filtered := getInodes(output)
	for i := range filtered {
		fmt.Print(string(filtered[i]),"\n")
	}
}


func getInodes(content []byte) [][]byte {
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
		if !bytes.Equal(fields[3],[]byte("0A")) {
			break
		}
		inodes = append(inodes, fields[9])
	}
	return inodes
}


