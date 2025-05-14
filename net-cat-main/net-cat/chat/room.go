package chat

import (
	"net"
	"time"
)

type room struct {
	name    string
	members map[net.Addr]*client
}

var R *room

func (r *room) broadcast(sender *client, msg string) {
	for addr, m := range r.members {
		if sender.conn.RemoteAddr() != addr {
			m.messages.Time = time.Now()
			m.msg(msg + m.messages.FormatMessage(*m.messages))
		}
	}
}

func (r *room) SameName(sender string) bool {
	for _, m := range r.members {
		if sender == m.nick {
			return false
		}
	}
	return true
}
