package bubblemachine

var withPieceInBufferStateImpl State = withPieceInBufferState{}

type withPieceInBufferState struct {
	machine *Machine
}

// GetStateName implements State.
func (w withPieceInBufferState) GetStateName() StateName {
	return WithPieceInBufferStateName
}

func newWithPieceInBuffer(machine *Machine) withPieceInBufferState {
	return withPieceInBufferState{machine}
}

func (w withPieceInBufferState) PutMoney(piece Piece) {
	w.machine.incrementNumberOfIgnoredTransition()
}

func (w withPieceInBufferState) Turn() *Bubble {
	w.machine.currentState = newIddleState(w.machine)
	b, bs := pop(w.machine.bubbles)
	w.machine.bubbles = bs
	return b
}

func pop(slice []*Bubble) (*Bubble, []*Bubble) {
	if len(slice) == 0 {
		return nil, slice // Or handle error as you wish
	}
	return slice[len(slice)-1], slice[:len(slice)-1]
}
