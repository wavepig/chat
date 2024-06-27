package lib

import (
	"chat-server/pkg/utils"
	"sync"
)

type HashSet struct {
	set map[string]bool
}

func NewHashSet() *HashSet {
	return &HashSet{make(map[string]bool)}
}

func (h *HashSet) Add(i string) bool {
	_, found := h.set[i]
	h.set[i] = true
	return !found
}

func (h *HashSet) Get(i string) bool {
	_, found := h.set[i]
	return found
}

func (h *HashSet) Remove(i string) {
	delete(h.set, i)
}

func (h *HashSet) Len() int {
	return len(h.set)
}

func (h *HashSet) List() []string {
	list := make([]string, 0, len(h.set))
	for k := range h.set {
		list = append(list, k)
	}
	return list
}

type Users struct {
	lock    sync.Mutex
	hashSet *HashSet
}

func NewUsers() *Users {
	return &Users{
		lock:    sync.Mutex{},
		hashSet: NewHashSet(),
	}
}

func (s *Users) Insert(name string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.hashSet.Add(name)
}

func (s *Users) Remove(name string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.hashSet.Remove(name)
}

func (s *Users) GetUnique() string {
	randStr := utils.RandString(8)
	s.lock.Lock()
	defer s.lock.Unlock()
	for s.hashSet.Get(randStr) {
		randStr = utils.RandString(8)
	}
	return randStr
}

func (s *Users) Len() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.hashSet.Len()
}

func (s *Users) List() []string {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.hashSet.List()
}
