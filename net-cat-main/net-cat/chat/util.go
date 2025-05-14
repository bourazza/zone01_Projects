package chat

import (
	"fmt"
	"time"
)

const (
	Red   = "\033[31m"
	Green = "\033[32m"
	Reset = "\033[0m"
)

type Message struct {
	Time    time.Time
	Sender  string
	Content string
}

var LogMessage string = `         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    '.       | '' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     '-'       '--'`

var History []string

func (m *Message) FormatMessage(msg Message) string {
	t := m.Time
	formatted := t.Format("2006-01-02 15:04:05")

	return fmt.Sprintf("[%s][%s]:%s", formatted, msg.Sender, msg.Content)
}

func ValidateLength(username string) bool {
	if len(username) < 3 || len(username) > 15 {
		return false
	}
	return true
}

func ValidName(username string) bool {
	for _, i := range username {
		if !(i >= 'A' && i <= 'Z') && !(i >= 'a' && i <= 'z') {
			return false
		}
	}
	return true
}

// validators for meessages
func ValidateLengthMessage(message string) bool {
	if len(message) < 1 || len(message) > 50 {
		return false
	}
	return true
}

func ValidMessage(message string) bool {
	for _, i := range message {
		if i < 32 || i > 126 {
			return false
		}
	}
	return true
}
