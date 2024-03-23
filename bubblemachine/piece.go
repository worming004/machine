package bubblemachine

type Piece int

func NewPiece(value int) Piece {
	return Piece(value)
}
func (p Piece) Value() int {
	return int(p)
}
