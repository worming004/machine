package bubblemachine

import "fmt"

type Machine struct {
	pieces                   []Piece
	currentState             State
	bubbles                  []*Bubble
	countOfIgnoredTransition int
}

func (m *Machine) GetStateName() StateName {
	return m.currentState.GetStateName()
}

// PutMoney implements State.
func (m *Machine) PutMoney(piece Piece) {
	m.currentState.PutMoney(piece)
}

// Turn implements State.
func (m *Machine) Turn() *Bubble {
	return m.currentState.Turn()
}

func (m *Machine) incrementNumberOfIgnoredTransition() {
	m.countOfIgnoredTransition++
}

func (m *Machine) CountOfIgnoredTransition() int {
	return m.countOfIgnoredTransition
}

var m State = &Machine{}

func NewMachine(options ...options) *Machine {
	m := &Machine{}
	m.currentState = newIddleState(m)
	for _, option := range options {
		option(m)
	}
	return m
}

type options func(m *Machine)

func WithBubbles(bubbles []*Bubble) options {
	return func(m *Machine) {
		m.bubbles = bubbles
	}
}

type Piece int

func (m *Machine) PrintState() {
	fmt.Printf("Current state: {pieces: %d, state: %T, bubble: %v, countOfIgnoredTransition: %d}\n", m.pieces, m.currentState, printableBubbles(m.bubbles), m.countOfIgnoredTransition)
}
