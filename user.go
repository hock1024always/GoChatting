package main

import (
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string //每一个用户都有一个消息通道 通过信道接受服务端的消息
	Conn net.Conn    //和客户端建立的连接
}

// 对channel进行监听
func (u *User) ListenMessage() {
	for {
		msg := <-u.C

		u.Conn.Write([]byte(msg + "\n")) //将数据写进连接中

	}
}

// 创建一个用户的API
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		Conn: conn,
	}

	//启动监听消息的协程
	go user.ListenMessage()

	return user
}
