package bubblemachine_test

import (
	"testing"

	"github.com/worming004/machine/bubblemachine"
)

func TestSimpleUsage(t *testing.T) {
	t.Parallel()

	uniqueBubbleName := "pokemon"

	m := bubblemachine.NewMachine(bubblemachine.WithBubbles(buildBubbles(uniqueBubbleName)))

	if m.GetStateName() != bubblemachine.IddleStateName {
		t.Errorf("Expected state to be %s, but got %s", bubblemachine.IddleStateName, m.GetStateName())
	}

	m.PutMoney(bubblemachine.Piece(1))

	if m.GetStateName() != bubblemachine.WithPieceInBufferStateName {
		t.Errorf("Expected state to be %s, but got %s", bubblemachine.WithPieceInBufferStateName, m.GetStateName())
	}

	b := m.Turn()

	if b.String() != uniqueBubbleName {
		t.Errorf("Expected bubble to be %s, but got %s", uniqueBubbleName, b.String())
	}
	if m.GetStateName() != bubblemachine.IddleStateName {
		t.Errorf("Expected state to be %s, but got %s", bubblemachine.IddleStateName, m.GetStateName())
	}

	if m.CountOfIgnoredTransition() != 0 {
		t.Errorf("Expected count of ignored transitions to be 0, but got %d", m.CountOfIgnoredTransition())
	}

	m.Turn()

	if m.CountOfIgnoredTransition() != 1 {
		t.Errorf("Expected count of ignored transitions to be 1, but got %d", m.CountOfIgnoredTransition())
	}
}

func buildBubbles(vs ...string) []*bubblemachine.Bubble {
	bubble := make([]*bubblemachine.Bubble, len(vs))
	for i, v := range vs {
		b := bubblemachine.Bubble(v)
		bubble[i] = &b
	}
	return bubble
}
