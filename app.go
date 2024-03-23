package main

import (
	"context"
	"log"

	"github.com/worming004/machine/bubblemachine"
	"github.com/worming004/machine/sqliterepository"
)

func main() {
	machine := bubblemachine.InitMachine(bubblemachine.WithBubbles(defaultBubbles()))
	repo, err := sqliterepository.NewMachineRepository(sqliterepository.NewRepositoryRequest{DataSourceName: "test.db", Init: true})
	if err != nil {
		panic(err)
	}
	defer repo.Close()
	err = repo.Save(context.Background(), machine)

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Machine saved successfully with Id ", machine.GetId())

}

func defaultBubbles() []bubblemachine.Bubble {
	bubble := make([]bubblemachine.Bubble, 3)

	for i, v := range []string{"Troll", "ToyCar", "Pokemon"} {
		b := bubblemachine.Bubble(v)
		bubble[i] = b
	}
	return bubble
}
