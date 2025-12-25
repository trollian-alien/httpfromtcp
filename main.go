package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	for {
		b := make([]byte, 8)
		n, err := file.Read(b)
		if err == io.EOF {
			break
		} else if err != nil {
			os.Exit(1)
		}
		fmt.Printf("read: %s\n", b[:n])
	}
}
