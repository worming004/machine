package bubblemachine

import (
	"fmt"
	"io"
	"os"
)

type machine struct {
	pieces                   []Piece
	currentState             State
	bubbles                  []*Bubble
	countOfIgnoredTransition int
	inspectWriter            io.Writer
}

func NewMachine(options ...options) *machine {
	m := &machine{inspectWriter: os.Stdout}
	m.currentState = newIddleState(m)
	for _, option := range options {
		option(m)
	}
	return m
}

func (m *machine) GetStateName() StateName {
	return m.currentState.GetStateName()
}

func (m *machine) PutMoney(piece Piece) {
	m.currentState.PutMoney(piece)
}

func (m *machine) Turn() *Bubble {
	return m.currentState.Turn()
}

func (m *machine) incrementNumberOfIgnoredTransition() {
	m.countOfIgnoredTransition++
}

func (m *machine) CountOfIgnoredTransition() int {
	return m.countOfIgnoredTransition
}

var m State = &machine{}

type options func(m *machine)

func WithBubbles(bubbles []*Bubble) options {
	return func(m *machine) {
		m.bubbles = bubbles
	}
}

func WithLogWriter(writer io.Writer) options {
	return func(m *machine) {
		m.inspectWriter = writer
	}
}

type Piece int

func (m *machine) PrintState() {
	fmt.Fprintf(m.inspectWriter, "Current state: {pieces: %d, state: %T, bubble: %v, countOfIgnoredTransition: %d}\n", m.pieces, m.currentState, printableBubbles(m.bubbles), m.countOfIgnoredTransition)
}
