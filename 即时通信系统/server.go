package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int
	//在线用户列表
	OnlineUser map[string]*User
	mapLock    sync.RWMutex
	//消息广播的channel
	Message chan string
	//用户登录的channel
}

// 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:         ip,
		Port:       port,
		OnlineUser: make(map[string]*User),
		Message:    make(chan string),
	}

	return server
}

// 监听消息的接口
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message
		//...消息处理的业务
		fmt.Println("收到消息:", msg)
		this.mapLock.Lock()
		for _, user := range this.OnlineUser {
			//...消息发送的业务
			user.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// 广播消息
func (this *Server) BroadCast(user *User, msg string) {
	//...消息广播的业务
	sendMsg := fmt.Sprintf("%s:%s", user.Name, msg)
	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	//...当前链接的业务
	//fmt.Println("链接建立成功")
	user := NewUser(conn, this)
	//用户上线消息广播，加入onlineMap
	user.Online()
	//监听用户是否活跃
	isLive := make(chan bool)
	//接收消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn.Read err:", err)
				return
			}
			msg := string(buf[:n-1])
			//...消息处理的业务
			user.DoMessage(msg)
			//用户任意消息，表示活跃
			isLive <- true
		}

	}()

	//阻塞
	for {
		select {
		case <-isLive:
			//...活跃处理的业务
			//fmt.Println("活跃处理")
		case <-time.After(time.Second * 100):
			//...超时处理的业务
			user.SendMessage("offline")
			//销毁占用资源
			close(user.C)
			//关闭链接
			conn.Close()
			//用户下线消息广播，从onlineMap中删除
			user.Offline()
			return
		}
	}

}

// 启动服务器的接口
func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	//close listen socket
	defer listener.Close()
	//启动消息监听
	go this.ListenMessage()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		//do handler
		go this.Handler(conn)
	}
}
