package main

import (
	"fmt"

	"github.com/worming004/machine/bubblemachine"
)

func main() {
	machine := bubblemachine.NewMachine(bubblemachine.WithBubbles(defaultBubbles()))
	machine.PrintState()
	machine.Turn() // Should increment the number of ignored transitions
	machine.PrintState()
	machine.PutMoney(bubblemachine.Piece(1))
	machine.PrintState()

	machine.PutMoney(bubblemachine.Piece(2)) // Should increment the number of ignored transitions
	machine.PrintState()
	bubble := machine.Turn()
	machine.PrintState()

	fmt.Printf("Bubble: %s\n", bubble)
}

func defaultBubbles() []bubblemachine.Bubble {
	bubble := make([]bubblemachine.Bubble, 3)

	for i, v := range []string{"Troll", "ToyCar", "Pokemon"} {
		b := bubblemachine.Bubble(v)
		bubble[i] = b
	}
	return bubble
}
