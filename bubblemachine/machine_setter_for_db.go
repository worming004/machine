package bubblemachine

import (
	"fmt"
	"os"
)

// This is a helper for database. Do not use it elsewhere
type machineSetterForDb struct {
	mach *Machine
}

func NewMachineSetterForDb(m *Machine) *machineSetterForDb {
	if m.inspectWriter == nil {
		m.inspectWriter = os.Stdout
	}
	return &machineSetterForDb{m}
}

func (m *machineSetterForDb) SetBubbles(bs []Bubble) *machineSetterForDb {
	m.mach.bubbles = bs
	return m
}

func (m *machineSetterForDb) SetPieces(ps []Piece) *machineSetterForDb {
	m.mach.pieces = ps
	return m
}

func (m *machineSetterForDb) SetCount(c int) *machineSetterForDb {
	m.mach.countOfIgnoredTransition = c
	return m
}

func (m *machineSetterForDb) SetStateByName(sn StateName) (*machineSetterForDb, error) {
	switch sn {
	case IddleStateName:
		m.mach.currentState = newIddleState(m.mach)
		break
	case WithPieceInBufferStateName:
		m.mach.currentState = newWithPieceInBuffer(m.mach)
		break
	default:
		return nil, fmt.Errorf("Invalid StateName: %s", sn)
	}

	return m, nil
}
