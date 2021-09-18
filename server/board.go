package main

import (
	"fmt"
	"strings"
	"unicode"
)

// Board handles game logic about the board and drawing board to console
type Board struct {
	grid *[Size][Size]Spot
}

// Setup creates the initial chess board
func (b *Board) Setup() {
	board := [Size][Size]Spot{}

	for rank := 0; rank < Size; rank++ {
		for file := 0; file < Size; file++ {
			board[file][rank] = Spot{file: file, rank: rank}
		}
	}

	b.grid = &board

	// generate starting position
	startingFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	b.GenerateFromFENString(startingFEN)
}

// GenerateFromFENString creates a particular board position from a provided valid FEN string
func (b *Board) GenerateFromFENString(fen string) {
	piecePlacements := strings.Split(fen, "/")

	last := strings.Split(piecePlacements[7], " ")
	piecePlacements[7] = last[0]
	fields := last[1:]

	// current turn
	if fields[0] == "b" {
		Game.turn = Black
	} else {
		Game.turn = White
	}

	// castling rights
	Game.blackCastling = &CastlingRights{false, false}
	Game.whiteCastling = &CastlingRights{false, false}

	for _, rights := range fields[1] {
		if unicode.IsLower(rights) {
			if rights == 'k' {
				Game.blackCastling.kingside = true
			} else {
				Game.blackCastling.queenside = true
			}
		} else {
			if rights == 'K' {
				Game.whiteCastling.kingside = true
			} else {
				Game.whiteCastling.queenside = true
			}
		}
	}

	// en passant targets
	if fields[2] != "-" {
		file, rank := b.locationToFileAndRank(fields[2])
		b.grid[file][rank].passantTarget = 2
	}

	// fullmoves and halfmoves
	Game.halfmoves = int(fields[3][0] - '0')
	Game.fullmoves = int(fields[4][0] - '0')

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

			b.grid[file][rank].containsPiece = true
			b.grid[file][rank].piece = &Piece{color: color, class: class}
			file++
		}
	}
}

func (b *Board) locationToFileAndRank(loc string) (int, int) {
	file := int(loc[0] - 'a')
	rank := 8 - int(loc[1]-'0')
	fmt.Printf(" file: %v rank: %v", file, rank)
	return file, rank
}

// GetValidMoves returns the moves a piece can play if the given spot contains a piece else returns {}
func (b *Board) GetValidMoves(s *Spot) []Spot {
	if !s.containsPiece {
		return []Spot{}
	}

	validMoves, _ := s.piece.FindValidMoves(b.grid, s.file, s.rank, Black, true)
	return validMoves
}

// IsSpotOffBoard returns true if the spot is not on the board
func (b *Board) IsSpotOffBoard(file int, rank int) bool {
	return file < 0 || file > Size-1 || rank < 0 || rank > Size-1
}

// MovePiece moves a piece from start to destination
func (b *Board) MovePiece(start *Spot, destination *Spot) {
	turnSuccessful := true

	piece := start.piece
	destinationPiece := destination.piece
	piece.moves++

	start.piece = nil
	start.containsPiece = false

	destination.piece = piece
	destination.containsPiece = true

	// if it is pawns first move and moved 2 places make spot behind pawn en passant target for next turn
	if destination.piece.class == Pawn && destination.piece.moves == 1 && (destination.rank == 4 || destination.rank == 3) {
		b.grid[destination.file][destination.rank+1].passantTarget = 2
	}

	// if pawn move results in en passant, take piece behind destination
	var passantPiece *Piece
	if destination.passantTarget > 0 {
		passantSpot := &b.grid[destination.file][start.rank]
		passantPiece = passantSpot.piece
		passantSpot.piece = nil
		passantSpot.containsPiece = false
	}

	// if move puts player's king in check then revert the move
	opponentColor := Black
	if Game.turn == Black {
		opponentColor = White
	}

	if b.IsKingInCheck(Game.turn, opponentColor, nil) {
		piece.moves--

		start.piece = piece
		start.containsPiece = true

		destination.piece = destinationPiece
		destination.containsPiece = destinationPiece != nil

		if destination.passantTarget > 0 {
			passantSpot := &b.grid[destination.file][start.rank]
			passantSpot.piece = passantPiece
			passantSpot.containsPiece = true
		}

		turnSuccessful = false
	}

	// if turn was successfully played then pass turn to opponent
	if turnSuccessful {
		Game.NextTurn(Game.turn, opponentColor)
	}
}

// IsKingInCheck goes through each opponent piece on the board and checks if they are attacking
// color's king
// returns either true (the king is in check) or false (the king is not in check)
func (b *Board) IsKingInCheck(color int, opponentColor int, simulatedBoard *[Size][Size]Spot) bool {
	board := b.grid
	if simulatedBoard != nil {
		board = simulatedBoard
	}

	// check if any opponent's piece puts the king in check
	for rank := 0; rank < Size; rank++ {
		for file := 0; file < Size; file++ {
			if board[file][rank].containsPiece && board[file][rank].piece.color == opponentColor {
				_, inCheck := board[file][rank].piece.FindValidMoves(board, file, rank, color, false)
				if inCheck {
					return true
				}
			}
		}
	}

	return false
}

// IsStalemate returns true if color cannot play any moves but is not in check
func (b *Board) IsStalemate(color int, opponentColor int) bool {
	king := b.GetKingSpot(color)
	kingMoves, _ := king.piece.FindValidMoves(b.grid, king.file, king.rank, opponentColor, true)

	// when king cannot move out of check
	// check if any move by color can get king out of check
	if len(kingMoves) == 0 {
		for rank := 0; rank < Size; rank++ {
			for file := 0; file < Size; file++ {
				if b.grid[file][rank].containsPiece && b.grid[file][rank].piece.color == color && b.grid[file][rank].piece.class != King {
					piece := b.grid[file][rank].piece
					moves, _ := piece.FindValidMoves(b.grid, file, rank, opponentColor, true)

					// since moves are pruned for illegal moves
					// if any move is available then it will put king out of check
					if len(moves) > 0 {
						return false
					}
				}
			}
		}
	} else {
		return false
	}

	return true
}

// GetKingSpot returns the spot that contains the king of color color
func (b *Board) GetKingSpot(color int) *Spot {
	var king *Spot
	for rank := 0; rank < Size; rank++ {
		for file := 0; file < Size; file++ {
			if b.grid[file][rank].containsPiece && b.grid[file][rank].piece.class == King && b.grid[file][rank].piece.color == color {
				king = &b.grid[file][rank]
			}
		}
	}

	return king
}
