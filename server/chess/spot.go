package chess

// Spot identifies a location on the board
type Spot struct {
	piece         *Piece
	containsPiece bool
	file          int
	rank          int
	passantTarget int
}

func (s *Spot) GetFile() int {
	return s.file
}

func (s *Spot) GetRank() int {
	return s.rank
}
