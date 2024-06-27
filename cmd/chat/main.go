package main

import (
	"chat-server/internal/chat"
	"fmt"
)

func main() {
	server := chat.NewServer()
	err := server.Run()
	if err != nil {
		fmt.Println("start server error: ", err)
	}
}
