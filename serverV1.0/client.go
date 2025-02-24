package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIP   string
	ServerPort int
	Name       string
	// Password   string
	conn   net.Conn
	choice int
}

// func NewClient(serverIP string, serverPort int, username string, password string) *Client {
func NewClient(serverIP string, serverPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIP:   serverIP,
		ServerPort: serverPort,
		//Username:   username,
		//Password:   password,
		Name:   "bruce",
		choice: 999,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
	if err != nil {
		fmt.Println("connect server failed:", err)
		return nil
	}
	client.conn = conn
	return client

}

// 处理server回应消息，显示到输出
func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
	//有数据就copy到stdout，永久阻塞监听
	// for{
	// 	buf := make([]byte, 1024)
	// 	client.conn.Read(buf)
	// 	fmt.Println(string(buf))
	// }
}

func (client *Client) menu() bool {
	var choice int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&choice)

	if choice >= 0 && choice <= 3 {
		client.choice = choice
		return true
	} else {
		fmt.Println(">>>>请输入合法范围内的数字<<<<")
		return false
	}
}

// 更新用户名
func (client *Client) updateUsername() bool {
	fmt.Println(">>>>请输入用户名<<<<")
	fmt.Scanln(&client.Name)
	//发送更新用户名消息
	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("send rename msg failed:", err)
		return false
	}
	fmt.Println("update username success")
	return true
}

// 公聊模式
func (client *Client) publicChat() bool {
	var Msg string
	//fmt.Println(">>>>请输入公聊消息内容<<<<")
	fmt.Scanln(&Msg)
	//发送公聊消息
	for Msg != "exit" {
		sendMsg := Msg + "\n"
		_, err := client.conn.Write([]byte(sendMsg))
		if err != nil {
			fmt.Println("send public chat msg failed:", err)
			return false
		}
		fmt.Println("send public chat msg success")
		Msg = ""
		//fmt.Println(">>>>请输入公聊消息内容<<<<")
		fmt.Scanln(&Msg)

	}
	return true
}

// 私聊模式
func (client *Client) privateChat() bool {
	var Msg string
	var TargetName string
	fmt.Println(">>>>请输入私聊用户<<<<")
	fmt.Scanln(&TargetName)

	//发送私聊消息
	for TargetName != "exit" {
		fmt.Println(">>>>请输入私聊消息内容<<<<")
		fmt.Scanln(&Msg)
		for Msg != "exit" {
			sendMsg := "to|" + TargetName + "|" + Msg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("send private chat msg failed:", err)
				return false
			}
			fmt.Println("send private chat msg success")
			fmt.Println(">>>>请输入私聊消息内容<<<<")
			fmt.Scanln(&Msg)
		}
		fmt.Println(">>>>请输入私聊用户<<<<")
		fmt.Scanln(&TargetName)
	}
	return true
}

func (client *Client) Run() {
	for client.choice != 0 {
		for client.menu() != true {
			continue
		}
		//fmt.Println("选择模式")
		switch client.choice {
		case 1:
			// fmt.Println("公聊模式")
			client.publicChat()
			break
		case 2:
			//fmt.Println("私聊模式")
			client.privateChat()
			break
		case 3:
			client.updateUsername()
			//fmt.Println("更新用户名")
			break
		}

	}
}

var serverIp string
var serverPort int

// ./client -ip 127.0.0.1 -port 8888 -username bruce -password 0325
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "server ip")
	flag.IntVar(&serverPort, "port", 8888, "server port")
}

func main() {
	//命令行解析
	flag.Parse()
	client := NewClient(serverIp, serverPort)

	if client == nil {
		fmt.Println("connect server failed")
		return
	}
	//单独开启一个goroutine去处理server的回执消息
	go client.DealResponse()
	//发送登录消息
	fmt.Println("connect server success")
	//启动客户端业务
	client.Run()
}
