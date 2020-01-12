package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func checkError(err error) {
	if err != nil {
		fmt.Println("error:", err.Error())
		//0成功
		os.Exit(1)
	}
}
func MessageSend(conn net.Conn) {
	var input string
	for {
		//fmt.Println("Please input your message(Enter键发送):")
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		input = string(data)
		if strings.ToUpper(input) == "EXIT" {
			conn.Close()
			break
		} else {
			_, err := conn.Write([]byte(input))
			if err != nil {
				conn.Close()
				fmt.Println("client connect failure:", err.Error())

				break
			} else {
				fmt.Println("发送成功!")
			}

		}
		//fmt.Println(data)
	}

}
func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	checkError(err)
	defer conn.Close()
	go MessageSend(conn)
	buf := make([]byte, 1024)
	for {
		_, err := conn.Read(buf)
		checkError(err)
		//打印收到的消息
		fmt.Println("server received message content:", string(buf))
	}
	fmt.Println("The client program end!")

}
