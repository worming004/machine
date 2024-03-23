package bubblemachine

import (
	"fmt"
	"os"
)

// This is a helper for database. Do not use it elsewhere
type MachineSetterForDb struct {
	mach *Machine
}

func NewMachineSetterForDb(m *Machine) *MachineSetterForDb {
	if m.inspectWriter == nil {

		m.inspectWriter = os.Stdout
	}
	return &MachineSetterForDb{m}
}

func (m *MachineSetterForDb) SetBubbles(bs []Bubble) *MachineSetterForDb {
	m.mach.bubbles = bs
	return m
}

func (m *MachineSetterForDb) SetPieces(ps []Piece) *MachineSetterForDb {
	m.mach.pieces = ps
	return m
}

func (m *MachineSetterForDb) SetCount(c int) *MachineSetterForDb {
	m.mach.countOfIgnoredTransition = c
	return m
}

func (m *MachineSetterForDb) SetStateByName(sn StateName) (*MachineSetterForDb, error) {
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
