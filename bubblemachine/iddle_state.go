package bubblemachine

var iddleState State = IddleState{}

type IddleState struct {
	machine *Machine
}

// GetStateName implements State.
func (i IddleState) GetStateName() StateName {
	return IddleStateName
}

func newIddleState(machine *Machine) IddleState {
	return IddleState{machine}
}

func (i IddleState) PutMoney(piece Piece) {
	i.machine.pieces = append(i.machine.pieces, piece)
	i.machine.currentState = newWithPieceInBuffer(i.machine)
}

func (i IddleState) Turn() *Bubble {
	i.machine.incrementNumberOfIgnoredTransition()
	return nil
}
