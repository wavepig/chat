package main

import (
	"bufio"
	"chat-server/pkg/codec"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	message := make(chan string, 10)

	go handle(ctx, cancel, message)
	inputReader := bufio.NewReader(os.Stdin)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			input, _ := inputReader.ReadString('\n') // 读取用户输入
			inputInfo := strings.Trim(input, "\r\n")
			if strings.ToUpper(inputInfo) == "Q" { // 如果输入q就退出
				cancel()
				return
			}
			message <- inputInfo
		}
	}
}

func handle(ctx context.Context, cancel context.CancelFunc, message <-chan string) {
	conn, err := net.Dial("tcp", "127.0.0.1:42069")
	if err != nil {
		fmt.Println("连接失败:", err)
		cancel()
		return
	}

	defer conn.Close() // 关闭连接

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case data := <-message:
				encode, _ := codec.Encode(data)
				_, err = conn.Write(encode) // 发送数据
				if err != nil {
					fmt.Println("conn write error:", err.Error())
					cancel()
					return
				}
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			reader := bufio.NewReader(conn)
			msg, err := codec.Decode(reader)
			if err == io.EOF {
				continue
			}
			if err != nil {
				fmt.Println("decode msg failed, err:", err)
				cancel()
				return
			}
			if msg != "" {
				fmt.Println(msg)
			}
		}
	}
}
