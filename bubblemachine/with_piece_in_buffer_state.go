package bubblemachine

var withPieceInBufferState State = WithPieceInBufferState{}

type WithPieceInBufferState struct {
	machine *Machine
}

// GetStateName implements State.
func (w WithPieceInBufferState) GetStateName() StateName {
	return WithPieceInBufferStateName
}

func newWithPieceInBuffer(machine *Machine) WithPieceInBufferState {
	return WithPieceInBufferState{machine}
}

func (w WithPieceInBufferState) PutMoney(piece Piece) {
	w.machine.incrementNumberOfIgnoredTransition()
}

func (w WithPieceInBufferState) Turn() *Bubble {
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