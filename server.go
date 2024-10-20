package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

// 创建一个server类
type Server struct {
	Ip   string
	Port int

	//在线用户的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex //读写锁 确保同步机制

	// 广播消息的channel
	Message chan string
}

// 处理函数
func (s *Server) handleConnection(conn net.Conn) {
	user := NewUser(conn)
	//将上线用户加入到在线用户列表
	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()
	//广播上线消息
	s.Broadcast(user, "comes online")
	//接收用户消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf) //从客户端输入中读取数据 成功的话返回消息字节数，失败的话返回错误
			if n == 0 {
				s.Broadcast(user, "leaves")
				return
			}
			if err != nil || err != io.EOF {
				fmt.Println("Conn reading error:", err)
				return
			}

			//将字节形式的消息转化为字符串 去除结尾的\n
			msg := string(buf[:n-1])
			//广播消息
			s.Broadcast(user, msg)
		}
	}()

	//处理阻塞
	select {}
}

// 创建一个方法，用来启动服务
func (s *Server) Start() {
	//socket监听
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	// 关闭监听 程序结束后自动关闭
	defer listener.Close()

	// 开启监听协程
	go s.ListenBroadcast()

	for {
		//获取信息
		conn, err := listener.Accept() //标识当前已经有用户上线了
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}

		//处理信息
		go s.handleConnection(conn)
	}
}

// 广播消息的方法
func (s *Server) Broadcast(user *User, message string) {
	sendMsg := "[" + user.Addr + "] " + user.Name + ": " + message
	s.Message <- sendMsg
}

// 监听mseeage广播channel的协程，一旦有消息就发送给全部的user
func (s *Server) ListenBroadcast() {
	for {
		msg := <-s.Message

		// 遍历在线用户列表，发送消息
		s.mapLock.RLock()
		for _, client := range s.OnlineMap {
			client.C <- msg
		}
		s.mapLock.RUnlock()
	}
}

// 创建一个server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}
