package chat

import (
	"fmt"
	"log"
	"net"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func NewServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) Run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_JOIN:
			s.join(cmd.client)
		case CMD_MSG:
			s.msg(cmd.client, cmd.msg)
		case CMD_QUIT:
			s.quit(cmd.client)
		}
	}
}

func (s *server) NewClient(conn net.Conn) *client {
	log.Printf("new client has joined: %s", conn.RemoteAddr().String())

	fmt.Fprint(conn, "Welcome to TCP-Chat!\n"+LogMessage+"\n[ENTER YOUR NAME]:")
	return &client{
		conn:     conn,
		commands: s.commands,
	}
}

func (s *server) join(c *client) {
	roomName := "room"

	_, ok := s.rooms[roomName]
	if !ok {
		R = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = R
	}

	R.members[c.conn.RemoteAddr()] = c
	c.room = R

	R.broadcast(c, fmt.Sprintf(Green+"\n%s joined the room\n"+Reset, c.nick))
}

func (s *server) msg(c *client, msg Message) {
	Message := "\n" + msg.FormatMessage(msg) + "\n" 
	c.room.broadcast(c, Message)
}

func (s *server) quit(c *client) {
	log.Printf("client has left the chat: %s", c.conn.RemoteAddr().String())
	Count--

	s.quitCurrentRoom(c)
	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf(Red+"\n%s has left the room\n"+Reset, c.nick))
	}
}
