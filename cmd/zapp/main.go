package main


import (
	"fmt"
	"os"
)


func main() {
	output, err := os.ReadFile("/proc/net/tcp")
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
	fmt.Print(string(output))
}
