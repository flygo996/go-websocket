package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

//原文链接：https://blog.csdn.net/rj2017211811/article/details/89629529

var onlineConns = make(map[string]net.Conn)

var messageQueue = make(chan string, 1000)
var quitChan chan bool

func CheckError(err error) {
	if err != nil {
		fmt.Println("error:", err.Error())
		//0成功
		os.Exit(1)
	}
}
func processInfo(conn net.Conn) {
	//规定最大输入汉字为1024
	buf := make([]byte, 1024)
	defer conn.Close()
	for {
		numOfBytes, err := conn.Read(buf)
		if err != nil {
			break
		}

		/*if numOfBytes!=0{
			remoteAddr:=conn.RemoteAddr()
			fmt.Println(remoteAddr)
			//因为buf是缓存的数据，所以可能会接收到垃圾数据
			//只接受当前输入的字符
			fmt.Printf("Has received message: %s\n",string(buf[0:numOfBytes]))
		}*/
		if numOfBytes != 0 {
			message := string(buf[0:numOfBytes])
			//打印服务器接收到的消息
			fmt.Printf("Received message content: %s\n", message)
			//把message写到messageQueue里
			messageQueue <- message
		}

	}
}
func ConsumeMessage() {
	for {
		select {

		case message := <-messageQueue:
			//对消息进行解析
			doProcessMessage(message)
		case <-quitChan:
			{
				break
			}

		}
	}
}
func doProcessMessage(message string) {
	contents := strings.Split(message, "#")
	if len(contents) > 1 {
		addr := contents[0]
		sendMessage := contents[1]
		fmt.Println(addr, sendMessage)
		addr = strings.Trim(addr, " ")
		if conn, ok := onlineConns[addr]; ok {
			_, err := conn.Write([]byte(sendMessage))
			if err != nil {
				fmt.Println("online send failure")
			} else {
				//fmt.Println("size:",size)

			}
		} else {
			//fmt.Println("11")
		}

	}

}
func main() {

	listen_scoket, err := net.Listen("tcp", "127.0.0.1:8080")
	CheckError(err)
	defer listen_scoket.Close()
	fmt.Println("Server is waiting....")
	go ConsumeMessage()
	for {
		conn, err := listen_scoket.Accept()
		CheckError(err)
		//把conn存储到onlineConns映射表里
		addr1 := fmt.Sprintf("%s", conn.RemoteAddr())
		onlineConns[addr1] = conn
		for addr := range onlineConns {
			fmt.Println(addr)
		}
		//通过协程处理连接
		go processInfo(conn)

	}
}
