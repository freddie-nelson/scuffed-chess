package chess

import socketio "github.com/googollee/go-socket.io"

// User stores name and time of user
type Player struct {
	name     string
	time     int
	opponent bool
	s        socketio.Conn
}

func NewPlayer(name string, opponent bool, s socketio.Conn) *Player {
	return &Player{name, 0, opponent, s}
}
