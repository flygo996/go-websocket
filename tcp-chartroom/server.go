package main

import (
	"fmt"
	"net"
)

// 客户端结构体
type Client struct {
	//用户通信
	C chan string
	//客户端名称
	name string
	//客户端地址
	addr string
}

//map存储在线用户
var is_online map[string]Client

//消息通讯
var messages = make(chan string)

//广播全局消息到客户端
func Message() {
	is_online = make(map[string]Client)

	// 循环读取 message 通道中的数据
	for {
		// 通道 message 中有数据读到 msg 中。 没有，则阻塞
		msg := <-messages

		// 一旦执行到这里，说明message中有数据了，解除阻塞。 遍历 map
		for _, client := range is_online {
			client.C <- msg // 把从Message通道中读到的数据，写到 client 的 C 通道中。
		}
	}
}

//生产消息函数
func MakeMsg(client Client, msg string) (buf string) {
	buf = "[" + client.addr + "]" + client.name + ": " + msg
	return
}

//发送消息给在线客户端
func WriteMsgToClient(conn net.Conn, client Client) {
	for msg := range client.C {
		conn.Write([]byte(msg))
	}
}

func Handler(conn net.Conn) {

	//把当前客户端添加到在线map中
	addr := conn.RemoteAddr().String()
	client := Client{make(chan string), addr, addr}
	//将当前客户端加入在线字典列表中
	is_online[addr] = client

	//创建一个协程，专门给当前客户端发消息
	go WriteMsgToClient(conn, client)

	//将用户上线的消息放到全局消息中
	messages <- MakeMsg(client, "login")

	// 创建一个新协程，循环读取用户发送的消息，广播给在线用户
	go func() {
		for {
			buf := make([]byte, 2048)
			//读取客户端数据
			n, _ := conn.Read(buf)
			if n == 0 {
				fmt.Printf("用户%s退出登录\n", client.name)
				//将当前用户从在线字典中删除
				delete(is_online, addr)
				//通知其他客户端该用户退出登录
				messages <- MakeMsg(client, "logout")
				return
			}
			msg := string(buf[:n])
			//将客户端发的消息加入到全局消息通道中
			messages <- MakeMsg(client, msg)
		}
	}()
	//让协程不停止运行
	for {

	}
}

func main() {
	//奖励tcp监听
	listen, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer listen.Close()

	//创建协程处理消息
	go Message()

	for {
		//阻塞等待客户端连接
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}
		//创建协程处理客户端事件
		go Handler(conn)
		defer conn.Close()
	}
}