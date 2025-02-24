# 一个基于GO语言的即时通信系统的开发
## 系统介绍
### 功能实现
1.用户上线及广播功能    
2.用户名更改功能    
3.用户公聊功能    
4.用户私聊功能    
5.超时强踢功能  
### user端
```
type User struct {
	Name   string
	Conn   net.Conn
	Addr   string
	C      chan string
	server *Server
}

//User的建立并通过server进行与服务端的关联
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
```
### server端
```
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
```
### client端
```

type Client struct {
	ServerIP   string
	ServerPort int
	Name       string
	conn   net.Conn
	choice int
}

// ./client -ip 127.0.0.1 -port 8888
//通过解析命令行获取IP及port
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
```

## 简单服务启动
### Windows
```
go run main.go user.go server.go
go run client.go -ip 127.0.0.1 -sort 8888
```
### Linux
```
go build -o server main.go user.go server.go
./server
go build -o client.go -ip 127.0.0.1 -sort 8888
./client
```
