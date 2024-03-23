package bubblemachine

import (
	"context"
	"fmt"
	"io"
	"os"
)

type Machine struct {
	id                       int64
	pieces                   []Piece
	currentState             State
	bubbles                  []Bubble
	countOfIgnoredTransition int
	inspectWriter            io.Writer
}

func (m *Machine) GetId() int64 {
	return m.id
}

func (m *Machine) SetId(id int64) {
	m.id = id
}

func InitMachine(options ...options) *Machine {
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

func (m *Machine) GetBubbles() []Bubble {
	// TODO should return a copy
	return m.bubbles
}

func (m *Machine) GetPieces() []Piece {
	// TODO should return a copy
	return m.pieces
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

func (m *Machine) PrintState() {
	fmt.Fprintf(m.inspectWriter, "Current state: {pieces: %d, state: %T, bubble: %v, countOfIgnoredTransition: %d}\n", m.pieces, m.currentState, printableBubbles(m.bubbles), m.countOfIgnoredTransition)
}

type MachineRepository interface {
	Save(ctx context.Context, m *Machine) error
	Get(ctx context.Context, id int) (*Machine, error)
}
