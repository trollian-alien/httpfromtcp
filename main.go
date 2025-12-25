package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	s := ""
	temp := ""
	go func() {
		defer f.Close()
		for {
			s += temp
			temp = ""
			b := make([]byte, 8)
			n, err := f.Read(b)
			str := string(b[:n])
			parts := strings.Split(str, "\n")

			s += parts[0]
			if len(parts) == 2 {
				temp += parts[1]
				ch <- s
				s = ""
			} 
			
			if err == io.EOF {
				close(ch)
				break
			} else if err != nil {
				os.Exit(1)
			}
		}
	}()
	return ch
}

func main() {
	f, err := os.Open("./messages.txt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	ch := getLinesChannel(f)

	for line := range(ch) {
		fmt.Println("read: " +line)
	}
}
