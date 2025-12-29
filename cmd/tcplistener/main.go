package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/trollian-alien/httpfromtcp/internal/request"
)

const port = ":42069"

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

		rq, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("Error listening: %v", err)
			os.Exit(1)
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %v\n", rq.RequestLine.Method)
		fmt.Printf("- Target: %v\n", rq.RequestLine.RequestTarget)
		fmt.Printf("- Version: %v\n", rq.RequestLine.HttpVersion)
		fmt.Println("Connection closed")
	}
}
