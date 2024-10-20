package main

// 作为进程的主入口
func main() {
	server := NewServer("127.0.0.1", 8080)
	server.Start()
}
