package bubblemachine

var iddleStateImpl State = iddleState{}

type iddleState struct {
	machine *Machine
}

// GetStateName implements State.
func (i iddleState) GetStateName() StateName {
	return IddleStateName
}

func newIddleState(machine *Machine) iddleState {
	return iddleState{machine}
}

func (i iddleState) PutMoney(piece Piece) {
	i.machine.pieces = append(i.machine.pieces, piece)
	i.machine.currentState = newWithPieceInBuffer(i.machine)
}

func (i iddleState) Turn() *Bubble {
	i.machine.incrementNumberOfIgnoredTransition()
	return nil
}
