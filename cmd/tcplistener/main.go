package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const port = ":42069"

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {
		defer f.Close()
		defer close(ch)
		s := "" // for storing partial lines that must be carried to next iteration of loop
		fmt.Println("Debugging:")
		for {
			b := make([]byte, 8)
			n, err := f.Read(b)

			if err == io.EOF {
				if s != "" {
					ch <- s
					fmt.Println(s)
				}
				break
			} else if err != nil {
				log.Fatalf("error listening TCP traffic: %s\n", err.Error())
			}

			str := string(b[:n])
			parts := strings.Split(str, "\n")

			for i := 0; i < len(parts) -1 ; i++ {
				ch <- s + parts[i]
				s = ""
			}
			s += parts[len(parts)-1]
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
		for line := range ch {
			fmt.Println(line)
		}
		fmt.Println("Connection closed")
	}
}
