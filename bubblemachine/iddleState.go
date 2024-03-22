package bubblemachine

type IddleState struct {
	machine *Machine
}

func NewIddleState(machine *Machine) IddleState {
	return IddleState{machine}
}

func (i IddleState) PutMoney(piece Piece) {
	i.machine.pieces = append(i.machine.pieces, piece)
	i.machine.currentState = NewWithPieceInBuffer(i.machine)
}

func (i IddleState) Turn() *Bubble {
	i.machine.incrementNumberOfIgnoredTransition()
	return nil
}
