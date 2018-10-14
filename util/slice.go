package util

import "golang.org/x/net/websocket"

func RemoveConn(s []*websocket.Conn, i int) []*websocket.Conn {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
