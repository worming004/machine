package bubblemachine

import (
	"fmt"
	"io"
	"os"
)

type Machine struct {
	pieces                   []Piece
	currentState             State
	bubbles                  []Bubble
	countOfIgnoredTransition int
	inspectWriter            io.Writer
}

func NewMachine(options ...options) *Machine {
	m := &Machine{inspectWriter: os.Stdout}
	m.currentState = newIddleState(m)
	for _, option := range options {
		option(m)
	}
	return m
}

func (m *Machine) GetStateName() StateName {
	return m.currentState.GetStateName()
}

func (m *Machine) PutMoney(piece Piece) {
	m.currentState.PutMoney(piece)
}

func (m *Machine) Turn() Bubble {
	return m.currentState.Turn()
}

func (m *Machine) incrementNumberOfIgnoredTransition() {
	m.countOfIgnoredTransition++
}

func (m *Machine) CountOfIgnoredTransition() int {
	return m.countOfIgnoredTransition
}

var m State = &Machine{}

type options func(m *Machine)

func WithBubbles(bubbles []Bubble) options {
	return func(m *Machine) {
		m.bubbles = bubbles
	}
}

func WithLogWriter(writer io.Writer) options {
	return func(m *Machine) {
		m.inspectWriter = writer
	}
}

type Piece int

func (m *Machine) PrintState() {
	fmt.Fprintf(m.inspectWriter, "Current state: {pieces: %d, state: %T, bubble: %v, countOfIgnoredTransition: %d}\n", m.pieces, m.currentState, printableBubbles(m.bubbles), m.countOfIgnoredTransition)
}

type MachineRepository interface {
	Save(m *Machine) error
	Get(id int) (*Machine, error)
}
