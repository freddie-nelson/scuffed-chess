package chess

import "log"

// Enum type of piece
const (
	Queen int = iota
	King
	Rook
	Bishop
	Knight
	Pawn
)

// Enum color of piece
const (
	Black int = iota
	White
)

// Piece : generic class for a chess piece
type Piece struct {
	color int
	class int
	moves int
}

// PruneMove simualates a move on a clone of the current board
// returns wether the move should be pruned or not (king is in check)
func (p *Piece) PruneMove(originalBoard *Board, file, rank, dFile, dRank int) bool {
	b := NewBoard()

	for rank := 0; rank < Size; rank++ {
		for file := 0; file < Size; file++ {
			s := originalBoard.grid[file][rank]

			var piece *Piece = nil
			if s.piece != nil {
				piece = &Piece{s.piece.color, s.piece.class, s.piece.moves}
			}

			b.grid[file][rank] = Spot{piece, s.containsPiece, file, rank, s.passantTarget}
		}
	}

	// log.Printf("start: (%d, %d) | dest: (%d, %d) | piece: (%d, %d) \n", start.file, start.rank, destination.file, destination.rank, piece.class, piece.color)

	return !b.MovePiece(&b.grid[file][rank], &b.grid[dFile][dRank], originalBoard.grid[file][rank].piece.color)
}

// FindValidMoves finds and returns all the legal moves a piece can make from it's current position
// @returns array of all valid moves
func (p *Piece) FindValidMoves(b *Board, file int, rank int, opponentColor int, pruneChecks bool, castlingRights *CastlingRights) ([]Spot, bool) {
	validMoves := make([]Spot, 0, 30)

	// bishop offsets
	bishopXOffs := []int{1, -1}
	bishopYOffs := []int{-1, 1}

	// rook offsets
	rookXOffs := []int{0, 0}
	rookYOffs := []int{-1, 1}

	// queen offsets
	queenOffs := []int{0, 0, -1, 1}

	// knight offsets
	knightXOffs := []int{2, -2}
	knightYOffs := []int{1, -1}

	checksKing := false

	switch p.class {
	case Queen:
		checksKing = calculateMovesFromOffsets(b, &validMoves, file, rank, queenOffs, queenOffs, Size, true, opponentColor)
	case King:
		checksKing = calculateMovesFromOffsets(b, &validMoves, file, rank, queenOffs, queenOffs, 1, true, opponentColor)

		// castling
		if p.moves == 0 && castlingRights != nil && !b.IsKingInCheck(p.color, opponentColor) {
			checkCastling(b, &validMoves, file, rank, castlingRights)
		}
	case Rook:
		checks := calculateMovesFromOffsets(b, &validMoves, file, rank, rookXOffs, rookYOffs, Size, true, opponentColor)
		checks2 := calculateMovesFromOffsets(b, &validMoves, file, rank, rookYOffs, rookXOffs, Size, true, opponentColor)
		checksKing = checks || checks2
	case Bishop:
		checksKing = calculateMovesFromOffsets(b, &validMoves, file, rank, bishopXOffs, bishopYOffs, Size, true, opponentColor)
	case Knight:
		checks := calculateMovesFromOffsets(b, &validMoves, file, rank, knightXOffs, knightYOffs, 1, true, opponentColor)
		checks2 := calculateMovesFromOffsets(b, &validMoves, file, rank, knightYOffs, knightXOffs, 1, true, opponentColor)
		checksKing = checks || checks2
	case Pawn:
		pawnYOff := -1
		if opponentColor == White {
			pawnYOff = 1
		}

		if p.moves == 0 {
			calculateMovesFromOffsets(b, &validMoves, file, rank, []int{0}, []int{pawnYOff}, 2, false, opponentColor)
		} else {
			calculateMovesFromOffsets(b, &validMoves, file, rank, []int{0}, []int{pawnYOff}, 1, false, opponentColor)
		}

		checksKing = checkIfPawnCanTake(b, &validMoves, file, rank, opponentColor)
	}

	// remove moves that cause king to be put in check
	for i := len(validMoves) - 1; i >= 0 && pruneChecks; i-- {
		// log.Printf("start: %d, i: %d, move: (%d, %d) \n", len(validMoves), i, validMoves[i].file, validMoves[i].rank)

		if p.PruneMove(b, file, rank, validMoves[i].file, validMoves[i].rank) {
			log.Printf("check pruned: (%d, %d) \n", validMoves[i].file, validMoves[i].rank)
			// remove move
			validMoves = append(validMoves[:i], validMoves[i+1:]...)
		}
	}

	return validMoves, checksKing
}

func calculateMovesFromOffsets(b *Board, validMoves *[]Spot, file int, rank int, xOffs []int, yOffs []int, stopAfter int, canTake bool, opponentColor int) bool {
	checksKing := false

	// use offsets to jump across board in the way the piece would
	for _, xOff := range xOffs {
		for _, yOff := range yOffs {

			// loop through each spot until end of board or stopAfter is reached
			// if we run into a spot that is occupied then break
			// highlight occupied spot before break if it is an opponent's piece
			for i := 1; i < Size && i <= stopAfter; i++ {
				currentFile := file + xOff*i
				currentRank := rank + yOff*i
				if b.IsSpotOffBoard(currentFile, currentRank) {
					break
				}

				spot := &b.grid[currentFile][currentRank]
				if spot.containsPiece {
					if spot.piece.color == opponentColor && canTake {
						if spot.piece.class == King {
							checksKing = true
							break
						}

						if !isMoveAlreadyAdded(validMoves, currentFile, currentRank) {
							*validMoves = append(*validMoves, Spot{file: currentFile, rank: currentRank})
						}
					}

					break
				} else {
					if !isMoveAlreadyAdded(validMoves, currentFile, currentRank) {
						*validMoves = append(*validMoves, Spot{file: currentFile, rank: currentRank})
					}
				}
			}
		}
	}

	return checksKing
}

func isMoveAlreadyAdded(validMoves *[]Spot, file int, rank int) bool {
	for _, move := range *validMoves {
		if move.file == file && move.rank == rank {
			return true
		}
	}

	return false
}

func checkCastling(b *Board, validMoves *[]Spot, f, r int, castlingRights *CastlingRights) {
	// queenside
	if castlingRights.queenside && b.grid[0][r].containsPiece && b.grid[0][r].piece.class == Rook && b.grid[0][r].piece.color == b.grid[f][r].piece.color && b.grid[0][r].piece.moves == 0 {
		canQueenside := true
		for file := 1; file < f; file++ {
			if b.grid[file][r].containsPiece || b.grid[f][r].piece.PruneMove(b, f, r, file, r) {
				canQueenside = false
			}
		}

		if canQueenside {
			*validMoves = append(*validMoves, Spot{file: 2, rank: r})
		}
	} else {
		castlingRights.queenside = false
	}

	// kingside
	if castlingRights.kingside && b.grid[Size-1][r].containsPiece && b.grid[Size-1][r].piece.class == Rook && b.grid[Size-1][r].piece.color == b.grid[f][r].piece.color && b.grid[Size-1][r].piece.moves == 0 {
		canKingside := true
		for file := f + 1; file < Size-1; file++ {
			if b.grid[file][r].containsPiece || b.grid[f][r].piece.PruneMove(b, f, r, file, r) {
				canKingside = false
			}
		}

		if canKingside {
			*validMoves = append(*validMoves, Spot{file: 6, rank: r})
		}
	} else {
		castlingRights.kingside = false
	}
}

func checkIfPawnCanTake(b *Board, validMoves *[]Spot, file int, rank int, opponentColor int) bool {
	// calculate positions on board
	lFile := file - 1
	rFile := file + 1
	nextRank := rank - 1
	if opponentColor == White {
		nextRank = rank + 1
	}

	checksKing := false

	// left file
	if !b.IsSpotOffBoard(lFile, nextRank) {

		// can pawn take diagonally
		if b.grid[lFile][nextRank].containsPiece && b.grid[lFile][nextRank].piece.color == opponentColor {
			if b.grid[lFile][nextRank].piece.class == King {
				checksKing = true
			} else {
				*validMoves = append(*validMoves, Spot{file: lFile, rank: nextRank})
			}
		} else if !b.grid[lFile][nextRank].containsPiece && b.grid[lFile][nextRank].passantTarget > 0 && b.grid[lFile][rank].containsPiece && b.grid[lFile][rank].piece.color == opponentColor { // can pawn take en passant
			*validMoves = append(*validMoves, Spot{file: lFile, rank: nextRank})
		}
	}

	// right file
	if !b.IsSpotOffBoard(rFile, nextRank) {

		// can pawn take diagonally
		if b.grid[rFile][nextRank].containsPiece && b.grid[rFile][nextRank].piece.color == opponentColor {
			if b.grid[rFile][nextRank].piece.class == King {
				checksKing = true
			} else {
				*validMoves = append(*validMoves, Spot{file: rFile, rank: nextRank})
			}
		} else if !b.grid[rFile][nextRank].containsPiece && b.grid[rFile][nextRank].passantTarget > 0 && b.grid[rFile][rank].containsPiece && b.grid[rFile][rank].piece.color == opponentColor { // can pawn take en passant
			*validMoves = append(*validMoves, Spot{file: rFile, rank: nextRank})
		}
	}

	return checksKing
}
