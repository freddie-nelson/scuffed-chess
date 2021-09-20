package chess

import (
	"fmt"
	"strings"
	"time"
	"unicode"
)

// CastlingRights stores what side a player can castle
type CastlingRights struct {
	queenside bool
	kingside  bool
}

// GameController controls top level game logic
type GameController struct {
	color         int
	opponentColor int
	board         *Board
	started       bool
	startTime     int
	ended         bool
	endState      string
	turn          int
	halfmoves     int
	fullmoves     int

	whiteCastling *CastlingRights
	blackCastling *CastlingRights

	You      *Player
	Opponent *Player

	code string
}

// NewGame creates the initial game state
func NewGame(code string) *GameController {
	var g GameController

	g.whiteCastling = &CastlingRights{true, true}
	g.blackCastling = &CastlingRights{true, true}

	// create board
	g.board = NewBoard()

	startingFEN := "p7/8/1k6/8/3K4/8/8/2Q5 b - - 9 5"
	g.fromFENString(startingFEN)

	return &g
}

func (g *GameController) StartGame() {
	if g.started {
		return
	}

	g.started = true
	g.startTime = int(time.Now().UnixNano() / 100000)
}

// MakeMove checks if move is valid and then plays that move if it is
// @returns wether the move was made or not
func (g *GameController) MakeMove(file, rank, dFile, dRank int) bool {
	if g.board.IsSpotOffBoard(file, rank) || g.board.IsSpotOffBoard(dFile, dRank) {
		return false
	}

	start := &g.board.grid[file][rank]
	dest := &g.board.grid[dFile][dRank]

	if !start.containsPiece || start.piece.color != g.turn {
		return false
	}

	valid := false
	validMoves := g.board.GetValidMoves(start, g.GetOpponentColor(g.turn))

	// check if dest is in validMoves
	for _, m := range validMoves {
		if m.file == dFile && m.rank == dRank {
			valid = true
		}
	}

	movedPiece := false
	if valid {
		movedPiece = g.board.MovePiece(start, dest, g.turn)
	}

	if movedPiece {
		// update time control
		player := g.CurrentlyPlaying()
		now := int(time.Now().UnixNano() / 100000)

		// correct last move time for first moves
		if player.timeOfLastMove == 0 {
			if g.turn == White {
				player.timeOfLastMove = g.startTime
			} else {
				player.timeOfLastMove = g.You.timeOfLastMove
			}
		}

		player.time -= now - player.timeOfLastMove
		player.timeOfLastMove = now

		g.NextTurn(g.turn, g.GetOpponentColor(g.turn))

		return true
	} else {
		return false
	}
}

func (g *GameController) GetValidMoves(file, rank, opponentColor int) []Spot {
	if g.board.IsSpotOffBoard(file, rank) {
		return []Spot{}
	}

	return g.board.GetValidMoves(&g.board.grid[file][rank], opponentColor)
}

func (g *GameController) GetOpponentColor(color int) int {
	if color == White {
		return Black
	} else {
		return White
	}
}

// CurrentlyPlaying returns the player that is currently playing or nil
func (g *GameController) CurrentlyPlaying() *Player {
	if !g.started {
		return nil
	}

	if g.turn == White {
		return g.You
	} else {
		return g.Opponent
	}
}

func (g *GameController) IsCurrentlyPlaying(p *Player) bool {
	return p != nil && p == g.CurrentlyPlaying()
}

func (g *GameController) BroadcastData() {
	fen := g.toFENString()

	g.You.s.Emit("game:fen", fen)
	g.Opponent.s.Emit("game:fen", fen)

	playersJSON := fmt.Sprintf("{ \"you\": { \"username\": \"%s\", \"time\": %d }, \"opponent\": { \"username\": \"%s\", \"time\": %d } }", g.You.name, g.You.time, g.Opponent.name, g.Opponent.time)
	g.You.s.Emit("game:players", playersJSON)
	g.Opponent.s.Emit("game:players", playersJSON)

	g.You.s.Emit("game:end-state", g.endState)
	g.Opponent.s.Emit("game:end-state", g.endState)
}

func (g *GameController) toFENString() string {
	fen := ""

	enPassantTarget := "-"

	empty := '0'
	for rank := 0; rank < Size; rank++ {
		for file := 0; file < Size; file++ {
			s := g.board.grid[file][rank]

			if s.passantTarget > 0 {
				enPassantTarget = g.fileAndRankToLocation(file, rank)
			}

			if s.containsPiece {
				if empty > '0' {
					fen += string(empty)
				}
				empty = '0'

				piece := 'k'
				switch s.piece.class {
				case King:
					piece = 'k'
				case Queen:
					piece = 'q'
				case Rook:
					piece = 'r'
				case Knight:
					piece = 'n'
				case Bishop:
					piece = 'b'
				case Pawn:
					piece = 'p'
				}

				if s.piece.color == White {
					fen += strings.ToUpper(string(piece))
				} else {
					fen += string(piece)
				}

			} else {
				empty++
			}
		}

		if empty > '0' {
			fen += string(empty)
		}
		empty = '0'
		fen += "/"
	}
	fen = fen[:len(fen)-1]

	// turn
	if g.turn == White {
		fen += " w "
	} else {
		fen += " b "
	}

	// castling rights
	if g.whiteCastling.kingside {
		fen += "K"
	}
	if g.whiteCastling.queenside {
		fen += "Q"
	}
	if g.blackCastling.kingside {
		fen += "k"
	}
	if g.blackCastling.queenside {
		fen += "q"
	}

	fen += " " + enPassantTarget

	fen += fmt.Sprintf(" %d", g.halfmoves)
	fen += fmt.Sprintf(" %d", g.fullmoves)

	return fen
}

// fromFENString creates a particular board position from a provided valid FEN string
func (g *GameController) fromFENString(fen string) {
	piecePlacements := strings.Split(fen, "/")

	last := strings.Split(piecePlacements[7], " ")
	piecePlacements[7] = last[0]
	fields := last[1:]

	// current turn
	if fields[0] == "b" {
		g.turn = Black
	} else {
		g.turn = White
	}

	// castling rights
	g.blackCastling = &CastlingRights{false, false}
	g.whiteCastling = &CastlingRights{false, false}

	for _, rights := range fields[1] {
		if unicode.IsLower(rights) {
			if rights == 'k' {
				g.blackCastling.kingside = true
			} else {
				g.blackCastling.queenside = true
			}
		} else {
			if rights == 'K' {
				g.whiteCastling.kingside = true
			} else {
				g.whiteCastling.queenside = true
			}
		}
	}

	// en passant targets
	if fields[2] != "-" {
		file, rank := g.locationToFileAndRank(fields[2])
		g.board.grid[file][rank].passantTarget = 2
	}

	// fullmoves and halfmoves
	g.halfmoves = int(fields[3][0] - '0')
	if len(fields[3]) == 2 {
		g.halfmoves += int(fields[3][1] - '0')
	}

	g.fullmoves = int(fields[4][0] - '0')
	if len(fields[4]) == 2 {
		g.halfmoves += int(fields[4][1] - '0')
	}

	// place pieces
	for rank, fenRank := range piecePlacements {
		file := 0

		for _, char := range fenRank {
			var color int
			var class int

			if unicode.IsNumber(char) {
				file += int(char - '0')
				continue
			} else if unicode.IsLower(char) {
				color = Black
			} else {
				color = White
			}

			switch unicode.ToUpper(char) {
			case 'Q':
				class = Queen
			case 'K':
				class = King
			case 'R':
				class = Rook
			case 'B':
				class = Bishop
			case 'N':
				class = Knight
			case 'P':
				class = Pawn
			}

			g.board.grid[file][rank].containsPiece = true
			g.board.grid[file][rank].piece = &Piece{color: color, class: class}
			file++
		}
	}
}

func (g *GameController) locationToFileAndRank(loc string) (int, int) {
	file := int(loc[0] - 'a')
	rank := 8 - int(loc[1]-'0')
	return file, rank
}

func (g *GameController) fileAndRankToLocation(file, rank int) string {
	return string('a'+file) + string('0'+(8-rank))
}

// NextTurn performs end game state checks and if game does not end then proceeds to next turn
func (g *GameController) NextTurn(color int, opponentColor int) {
	// check for winning conditions
	if g.board.IsStalemate(opponentColor, color) {
		g.ended = true

		if g.board.IsKingInCheck(opponentColor, color) {
			g.endState = "checkmate"
		} else {
			g.endState = "stalemate"
		}
	}

	g.halfmoves++
	if g.turn == Black {
		g.fullmoves++
		g.turn = White
	} else {
		g.turn = Black
	}

	if g.fullmoves == 50 {
		g.ended = true
	}
}
