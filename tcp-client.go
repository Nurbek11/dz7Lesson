package main

import (
	"fmt"
	"net"
)

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:6060")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	for {
		var source string
		fmt.Print("Enter the number: ")
		_, err := fmt.Scanln(&source)
		if err != nil {
			fmt.Println("incorrect format", err)
			continue
		}
		//send a message to the server
		if n, err := conn.Write([]byte(source));
			n == 0 || err != nil {
			fmt.Println(err)
			return
		}
		//get an answer
		fmt.Print("Quadratic:")
		buff := make([]byte, 1024)
		n, err := conn.Read(buff)
		if err != nil {
			break
		}
		fmt.Print(string(buff[0:n]))
		fmt.Println()
	}
}
