package chess

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

// SimulateMove simualates a move on a clone of the current board
// returns if the king (of the piece's color) is in check in the new board state
func (p *Piece) SimulateMove(b *Board, s *Spot, d *Spot) bool {
	simulatedBoard := *b.grid

	start := &simulatedBoard[s.file][s.rank]
	destination := &simulatedBoard[d.file][d.rank]

	piece := start.piece
	start.piece = nil
	start.containsPiece = false

	destination.piece = piece
	destination.containsPiece = true

	opponentColor := Black
	if piece.color == Black {
		opponentColor = White
	}

	if b.IsKingInCheck(piece.color, opponentColor, &simulatedBoard) {
		return true
	}

	return false
}

// FindValidMoves finds and returns all the legal moves a piece can make from it's current position
// @returns array of all valid moves
func (p *Piece) FindValidMoves(b *Board, file int, rank int, opponentColor int, pruneChecks bool) ([]Spot, bool) {
	validMoves := make([]Spot, 0)

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
	case Rook:
		checksKing = calculateMovesFromOffsets(b, &validMoves, file, rank, rookXOffs, rookYOffs, Size, true, opponentColor)
		if checksKing {
			calculateMovesFromOffsets(b, &validMoves, file, rank, rookYOffs, rookXOffs, Size, true, opponentColor)
		} else {
			checksKing = calculateMovesFromOffsets(b, &validMoves, file, rank, rookYOffs, rookXOffs, Size, true, opponentColor)
		}
	case Bishop:
		calculateMovesFromOffsets(b, &validMoves, file, rank, bishopXOffs, bishopYOffs, Size, true, opponentColor)
	case Knight:
		calculateMovesFromOffsets(b, &validMoves, file, rank, knightXOffs, knightYOffs, 1, true, opponentColor)
		calculateMovesFromOffsets(b, &validMoves, file, rank, knightYOffs, knightXOffs, 1, true, opponentColor)
	case Pawn:
		if p.moves == 0 {
			calculateMovesFromOffsets(b, &validMoves, file, rank, []int{0}, []int{-1}, 2, false, opponentColor)
		} else {
			calculateMovesFromOffsets(b, &validMoves, file, rank, []int{0}, []int{-1}, 1, false, opponentColor)
		}

		checksKing = checkIfPawnCanTake(b, &validMoves, file, rank, opponentColor)
	}

	// remove moves that cause king to be put in check
	for i := len(validMoves) - 1; i >= 0 && pruneChecks; i-- {
		if p.SimulateMove(b, &b.grid[file][rank], &validMoves[i]) {
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

func checkIfPawnCanTake(b *Board, validMoves *[]Spot, file int, rank int, opponentColor int) bool {
	// calculate positions on board
	lFile := file - 1
	rFile := file + 1
	nextRank := rank - 1

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
		} else if !b.grid[lFile][nextRank].containsPiece && b.grid[lFile][nextRank].passantTarget > 0 { // can pawn take en passant
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
		} else if !b.grid[rFile][nextRank].containsPiece && b.grid[rFile][nextRank].passantTarget > 0 { // can pawn take en passant
			*validMoves = append(*validMoves, Spot{file: rFile, rank: nextRank})
		}
	}

	return checksKing
}
