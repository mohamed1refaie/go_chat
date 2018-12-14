package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {

	listener, err := net.Listen("tcp", "localhost:8080")
	defer listener.Close()

	if err != nil {
		log.Fatal(err)
	}

	for {
		con, errs := listener.Accept()
		if errs != nil {
			log.Fatal(errs)
		}

		fmt.Println("Accepted Connection..")
		reader := bufio.NewReader(os.Stdin)
		terminalreader := bufio.NewReader(con)

		for {
			fmt.Print("You: ")
			text, _ := reader.ReadString('\n')
			io.WriteString(con, "Stranger: "+text)
			io.WriteString(con, "You: ")
			textx, _ := terminalreader.ReadString('\n')
			fmt.Print("Stanger: " + textx)
		}
		con.Close()
	}
}
