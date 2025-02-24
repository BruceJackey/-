# 一个基于GO语言的即时通信系统的开发
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
