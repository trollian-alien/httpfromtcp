package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	b := make([]byte, 0)
	fmt.Print("read: ")
	defer file.Close()
	for {
		c := make([]byte, 8)
		n, err := file.Read(c)
		b = append(b, c[:n]...)
		if err == io.EOF {
			break
		} else if err != nil {
			os.Exit(1)
		}
	}
	s := string(b)
	s = strings.Replace(s, "\n", "\nread: ", -1)
	s = strings.TrimSuffix(s, "\nread: ")
    s = strings.TrimSuffix(s, "\n")

	fmt.Println(s)
}
