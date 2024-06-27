package lib

import "sync"

const (
	DEFAULT_ROOM = "DEFAULT_ROOM"
)

type Room struct {
	message  chan string
	users    *HashSet
	userChan map[string]chan string
}

func NewRoom() *Room {
	room := &Room{
		message:  make(chan string, 128),
		users:    NewHashSet(),
		userChan: make(map[string]chan string),
	}
	go room.publish()
	return room
}

func (s *Room) publish() {
	for {
		select {
		case msg := <-s.message:
			for _, user := range s.users.List() {
				s.userChan[user] <- msg
			}
		}
	}
}

type Rooms struct {
	lock sync.Mutex
	room map[string]*Room
}

func NewRooms() *Rooms {
	room := &Rooms{
		lock: sync.Mutex{},
		room: make(map[string]*Room, 0),
	}
	return room
}

// 进入房间
func (s *Rooms) Join(roomName, userName string) (chan string, chan string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	r, ok := s.room[roomName]
	if !ok {
		room := NewRoom()
		room.users.Add(userName)
		msgChan := make(chan string, 8)
		room.userChan[userName] = msgChan
		s.room[roomName] = room

		return room.message, msgChan
	}
	r.users.Add(userName)
	msgChan := make(chan string, 8)
	r.userChan[userName] = msgChan

	return r.message, msgChan
}

// 退出房间
func (s *Rooms) Leave(roomName, userName string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.room[roomName].users.Remove(userName)

	if s.room[roomName].users.Len() < 1 {
		delete(s.room, roomName)
	}
}

// 改变房间
func (s *Rooms) Change(prevRoom, nextRoom, userName string) (chan string, chan string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	r := s.room[prevRoom]
	userChan := r.userChan[userName]

	next, ok := s.room[nextRoom]
	if !ok {
		room := NewRoom()
		room.users.Add(userName)
		room.userChan[userName] = userChan

		s.room[nextRoom] = room

		r.users.Remove(userName)
		if r.users.Len() < 1 {
			delete(s.room, prevRoom)
		}

		return room.message, userChan
	}
	next.users.Add(userName)

	next.userChan[userName] = userChan

	r.users.Remove(userName)
	if r.users.Len() < 1 {
		delete(s.room, prevRoom)
	}

	return next.message, userChan
}

func (s *Rooms) ChangeName(roomName, prevName, nextName string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	room := s.room[roomName]
	room.users.Remove(prevName)
	room.users.Add(nextName)
	room.userChan[nextName] = room.userChan[prevName]
	delete(room.userChan, prevName)
}

func (s *Rooms) ListRoom() map[string]int {
	s.lock.Lock()
	defer s.lock.Unlock()

	roomList := make(map[string]int, 0)
	for k, v := range s.room {
		roomList[k] = v.users.Len()
	}
	return roomList
}

func (s *Rooms) ListUser(roomName string) []string {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.room[roomName].users.List()
}
