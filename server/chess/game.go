package chess

import (
	"strings"
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
	ended         bool
	endState      string
	turn          int
	halfmoves     int
	fullmoves     int

	whiteCastling *CastlingRights
	blackCastling *CastlingRights

	you      *Player
	opponent *Player

	timeOfLastTick int
	deltaTime      int
}

// Setup creates the initial game state
func (g *GameController) Setup() {
	// create board
	g.board.SetupBoard()

	startingFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	g.GenerateFromFENString(startingFEN)
}

// GenerateFromFENString creates a particular board position from a provided valid FEN string
func (g *GameController) GenerateFromFENString(fen string) {
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
	g.fullmoves = int(fields[4][0] - '0')

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

// NextTurn performs end game state checks and if game does not end then proceeds to next turn
func (g *GameController) NextTurn(color int, opponentColor int) {
	// check for winning conditions
	if g.board.IsStalemate(opponentColor, color) {
		g.ended = true

		if g.board.IsKingInCheck(opponentColor, color, nil) {
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
