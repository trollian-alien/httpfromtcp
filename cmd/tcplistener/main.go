package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"log"
	"net"
)

const port = ":42069"

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
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s\n", err.Error())
	}
	defer listener.Close()
	
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		fmt.Println("Accepted connection from", conn.RemoteAddr())

		ch := getLinesChannel(conn)
		for line := range(ch) {
		fmt.Println(line)
		}
		fmt.Println("Connection closed")
	}
}
