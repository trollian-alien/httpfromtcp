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
		for {
			b := make([]byte, 8)
			n, err := f.Read(b)
			str := string(b[:n])
			parts := strings.Split(str, "\r\n")

			//adding to the initial partial line
			if len(parts) > 0 {
				s += parts[0]
			}

			// middle parts are full lines
			l := len(parts)
			for i := 1; i < l-1; i++ {
				if s != "" {
					ch <- s
				}
				s = parts[i]
			}

			// handle the last part
			if len(parts) > 1 {
				last := parts[l-1]
				if strings.HasSuffix(str, "\r\n") {
					// line finished in this chunk
					ch <- s
					s = ""
				} else {
					// still partial, carry over
					s += last
				}
			}

			if err == io.EOF {
				if s != "" {
					ch <- s
				}
				break
			} else if err != nil {
				log.Fatalf("error listening TCP traffic: %s\n", err.Error())
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
		for line := range ch {
			fmt.Println(line)
		}
		fmt.Println("Connection closed")
	}
}
