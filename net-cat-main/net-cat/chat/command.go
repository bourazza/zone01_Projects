package chat

type commandID int

const (
	CMD_NICK commandID = iota
	CMD_JOIN
	CMD_MSG
	CMD_QUIT
)

var Count int

type command struct {
	id     commandID
	client *client
	msg    Message
}
