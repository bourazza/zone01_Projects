package chat

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

type client struct {
	conn     net.Conn
	nick     string
	room     *room
	messages *Message
	commands chan<- command
}

func (c *client) ReadInput() {
	var name bool
	First := true

	for {
		
		if History != nil && First && name {
			c.msg(strings.Join(History, "\n") + "\n")
		}
		if name {
			c.messages = &Message{
				Time:   time.Now(),
				Sender: c.nick,
			}

			fmt.Fprintf(c.conn, "%s", c.messages.FormatMessage(*c.messages))
		}
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		
		if err != nil {
			c.commands <- command{
				id:     CMD_QUIT,
				client: c,
			}
			return
		}

		msg = strings.Trim(msg, "\r\n")

		if !name {
			if msg == "" {
				fmt.Fprint(c.conn, "[ENTER YOUR NAME]:")
				continue
			}
			if !ValidateLength(msg) {
				fmt.Fprint(c.conn, "[ENTER YOUR NAME]:")
				continue
			}
			if !ValidName(msg) {
				fmt.Fprint(c.conn, "[ENTER YOUR NAME]:")
				continue
			}
			if R != nil && !R.SameName(msg) {
				fmt.Fprint(c.conn, "[ENTER YOUR NAME]:")
				continue
			}

			c.nick = msg
			name = true
			c.messages = &Message{
				Time:   time.Now(),
				Sender: c.nick,
			}
			c.commands <- command{
				id:     CMD_JOIN,
				client: c,
			}
			continue
		}
		First = false
		if msg == "" {
			continue
		}
		if !ValidateLengthMessage(msg) {
			continue
		}
		if !ValidMessage(msg) {
			continue
		}
		c.messages = &Message{
			Time:    time.Now(),
			Sender:  c.nick,
			Content: msg,
		}
		c.history(*c.messages)
		
		c.commands <- command{
			id:     CMD_MSG,
			client: c,
			msg:    *c.messages,
		}

	}
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte(msg))
}

func (c *client) history(msg Message) {
	History = append(History, msg.FormatMessage(msg))
}
