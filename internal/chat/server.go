package chat

import (
	"bufio"
	lib2 "chat-server/internal/chat/lib"
	"chat-server/pkg/codec"
	"context"
	"fmt"
	"log/slog"
	"net"
	"strings"
)

var help = `Server commands
  /help - print this message
  /name {name} - change name
  /rooms - list rooms
  /join {room} - join room
  /list - list server users
  /users - list users in room
  /quit - quit server`

type MessageChan struct {
	roomChan chan string
	userChan chan string
}

type Server struct {
	rooms *lib2.Rooms
	users *lib2.Users
}

func NewServer() *Server {
	return &Server{
		rooms: lib2.NewRooms(),
		users: lib2.NewUsers(),
	}
}

func (s *Server) Run() error {
	fmt.Println("start server")
	Listener, err := net.Listen("tcp", "127.0.0.1:42069")
	if err != nil {
		slog.Error("failed to connect", err)
		return err
	}
	defer Listener.Close()
	for {
		accept, err := Listener.Accept()
		if err != nil {
			slog.Error("failed to accept", err)
			continue
		}
		go s.handle(accept)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	ctx, cancelFunc := context.WithCancel(context.Background())
	name := s.users.GetUnique()
	s.users.Insert(name)
	Write(conn, help)
	room := lib2.DEFAULT_ROOM
	mc := &MessageChan{}
	mc.roomChan, mc.userChan = s.rooms.Join(room, name)
	mc.roomChan <- name + ": 欢迎进入[" + room + "]房间"

	go func(uc *MessageChan, ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-uc.userChan:
				Write(conn, msg)
			}
		}
	}(mc, ctx)

	for {
		msg, err := Read(conn)
		if err != nil {
			break
		}
		fmt.Printf("房间:[%s]-用户:[%s]-消息:[%s]\n", room, name, msg)
		// 字符串前缀匹配
		split := strings.Split(msg, " ")
		if strings.HasPrefix(msg, "/help") {
			Write(conn, help)
		} else if strings.HasPrefix(msg, "/name") {
			if len(split) < 2 {
				Write(conn, "命令传递错误")
				continue
			}
			nameSig := split[1]
			ok := s.users.Insert(nameSig)
			if ok {
				mc.roomChan <- name + ":修改名称为[" + nameSig + "]"
				s.users.Remove(name)
				s.rooms.ChangeName(room, name, nameSig)
				name = nameSig
				Write(conn, nameSig+":名称修改成功")
			} else {
				Write(conn, nameSig+"名称已存在")
			}
		} else if strings.HasPrefix(msg, "/quit") {
			break
		} else if strings.HasPrefix(msg, "/join") {
			if len(split) < 2 {
				Write(conn, "命令传递错误")
				continue
			}
			roomName := split[1]

			mc.roomChan, mc.userChan = s.rooms.Change(room, roomName, name)
			mc.roomChan <- name + ":进入[" + roomName + "]房间"
			room = roomName
		} else if strings.HasPrefix(msg, "/rooms") {
			list := ""
			for k, v := range s.rooms.ListRoom() {
				list += fmt.Sprintf("[%s:%d] ", k, v)
			}
			Write(conn, list)
		} else if strings.HasPrefix(msg, "/users") {
			list := s.rooms.ListUser(room)
			join := strings.Join(list, ",")
			Write(conn, join)
		} else if strings.HasPrefix(msg, "/list") {
			list := s.users.List()
			join := strings.Join(list, ",")
			Write(conn, join)
		} else {
			mc.roomChan <- name + ": " + msg
		}
	}

	fmt.Printf("[%s]退出[%s]房间\n", name, room)

	cancelFunc()
	s.users.Remove(name)
	s.rooms.Leave(room, name)
}

func Write(conn net.Conn, message string) {
	encode, err := codec.Encode(message)
	_, err = conn.Write(encode)
	if err != nil {
		slog.Error("failed to write", err)
	}
}

func Read(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	msg, err := codec.Decode(reader)
	if err != nil {
		slog.Error("failed to read", err)
		return "", err
	}
	return msg, nil
}
