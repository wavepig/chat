package utils

import (
	"math/rand"
	"sync"
	"time"
)

const _charsetRand = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@#$"

var _seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
var _lock = sync.Mutex{}

func RandStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	l := len(charset)
	_lock.Lock()
	defer _lock.Unlock()
	for i := range b {
		b[i] = charset[_seededRand.Intn(l)]
	}
	return string(b)
}

func RandString(length int) string {
	return RandStringWithCharset(length, _charsetRand)
}

func RandInt(min int, max int) int {
	if min <= 0 || max <= 0 {
		return 0
	}

	if min >= max {
		return max
	}

	return _seededRand.Intn(max-min) + min
}

func RandMax(max int) int {
	if max <= 1 {
		return 0
	}

	return _seededRand.Intn(max)
}
