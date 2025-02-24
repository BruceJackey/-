package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name   string
	Conn   net.Conn
	Addr   string
	C      chan string
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	useraddr := conn.RemoteAddr().String()
	user := &User{
		Name:   "User" + useraddr,
		Conn:   conn,
		Addr:   useraddr,
		C:      make(chan string),
		server: server,
	}
	go user.ListenMessage()
	fmt.Println("New User:", user.Name)
	return user
}

// 用户上线功能
func (this *User) Online() {
	//用户上线消息广播，加入onlineMap
	this.server.mapLock.Lock()
	this.server.OnlineUser[this.Name] = this
	this.server.mapLock.Unlock()
	//广播用户上线消息
	this.server.BroadCast(this, "online")

}

// 用户下线功能
func (this *User) Offline() {
	//用户下线消息广播，加入onlineMap
	this.server.mapLock.Lock()
	delete(this.server.OnlineUser, this.Name)
	this.server.mapLock.Unlock()
	//广播用户下线消息
	this.server.BroadCast(this, "offline")

}
func (this *User) SendMessage(msg string) {
	this.Conn.Write([]byte(msg + "\n"))
}

// 发送消息功能
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		//查询在线用户
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineUser {
			this.SendMessage(user.Name + " is online")
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		//重命名功能 rename|newname
		newname := msg[7:]
		this.server.mapLock.Lock()
		if _, ok := this.server.OnlineUser[newname]; ok {
			this.SendMessage("The name has been used, please choose another name")
		} else {
			delete(this.server.OnlineUser, this.Name)
			this.Name = newname
			this.server.OnlineUser[this.Name] = this
			this.server.mapLock.Unlock()
			this.SendMessage("You have changed your name to " + newname)
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// 私聊功能 to|username|message
		//获取对方用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			this.SendMessage("Please input the username of the user you want to chat with")
			return
		}
		//判断对方是否在线
		// 根据用户名得到对方的User对象
		this.server.mapLock.Lock()
		remoteUser, ok := this.server.OnlineUser[remoteName]
		this.server.mapLock.Unlock()
		if !ok {
			this.SendMessage(remoteName + " is not online, please try again later")
			return
		}

		// 发送私聊消息
		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.SendMessage("Please input the message you want to send")
			return
		}
		this.SendMessage("You have sent a message to " + remoteName + ": " + content)
		remoteUser.SendMessage(this.Name + " says: " + content)

	} else {
		//普通消息广播
		this.server.BroadCast(this, msg)
	}

}
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.Conn.Write([]byte(msg + "\n"))
	}
}
