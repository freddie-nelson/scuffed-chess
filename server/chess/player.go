package chess

import socketio "github.com/googollee/go-socket.io"

// User stores name and time of user
type Player struct {
	name           string
	time           int
	timeOfLastMove int
	opponent       bool
	s              socketio.Conn
}

func (p *Player) CompareID(id string) bool {
	return p.s.ID() == id
}

func (p *Player) GetSocket() socketio.Conn {
	return p.s
}

func NewPlayer(name string, opponent bool, s socketio.Conn) *Player {
	return &Player{name, 600000, 0, opponent, s}
}
