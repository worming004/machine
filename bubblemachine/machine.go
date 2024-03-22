package bubblemachine

import "fmt"

type Machine struct {
	pieces                   []Piece
	currentState             State
	bubbles                  []*Bubble
	countOfIgnoredTransition int
}

type StateName string

var (
	IddleStateName             StateName = "IddleState"
	WithPieceInBufferStateName StateName = "WithPieceInBufferState"
)

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

type State interface {
	PutMoney(piece Piece)
	Turn() *Bubble
	GetStateName() StateName
}

type Piece int

type Bubble string

func (b *Bubble) String() string {
	return string(*b)
}

func (m *Machine) PrintState() {
	fmt.Printf("Current state: {pieces: %d, state: %T, bubble: %v, countOfIgnoredTransition: %d}\n", m.pieces, m.currentState, printableBubbles(m.bubbles), m.countOfIgnoredTransition)
}

func printableBubbles(b []*Bubble) string {
	var s string
	for _, v := range b {
		s += string(*v) + ", "
	}
	return "[" + s[:len(s)-2] + "]"
}
