package bubblemachine

type State interface {
	PutMoney(piece Piece)
	Turn() Bubble
	GetStateName() StateName
}

type StateName string

var (
	IddleStateName             StateName = "IddleState"
	WithPieceInBufferStateName StateName = "WithPieceInBufferState"
)
