package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"time"
)

func main() {
	sleep := flag.Duration("sleep", 0, "How long to wait after each print")
	flag.Parse()

	listener, _ := net.Listen("tcp", "localhost:3040")
	conn, _ := listener.Accept()

	for {
		message, _ := bufio.NewReader(conn).ReadBytes('\n')
		fmt.Println(string(message))
		time.Sleep(*sleep)
	}
}
