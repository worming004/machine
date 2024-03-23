package main

import (
	"context"
	"log"
	"os"

	"github.com/worming004/machine/bubblemachine"
	"github.com/worming004/machine/sqliterepository"
)

func main() {
	shouldClean := os.Getenv("CLEAN_DB")
	if shouldClean == "true" || shouldClean == "1" {
		os.Remove("./test.db")
	}
	saveScenario()
	getScenario()
}

func saveScenario() {
	machine := bubblemachine.InitMachine(bubblemachine.WithBubbles(defaultBubbles()))
	machine.Turn()
	machine.PutMoney(4)
	createRepoAndSave(machine)
	log.Print("Machine saved successfully with Id ", machine.GetId())
}

func getScenario() {
	repo, err := sqliterepository.NewMachineRepository(sqliterepository.NewRepositoryRequest{DataSourceName: "test.db", Init: true})
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	m, err := repo.Get(context.Background(), 1)
	if err != nil {
		log.Fatal(err)
	}
	m.PrintState()
}

func defaultBubbles() []bubblemachine.Bubble {
	bubble := make([]bubblemachine.Bubble, 3)

	for i, v := range []string{"Troll", "ToyCar", "Pokemon"} {
		b := bubblemachine.Bubble(v)
		bubble[i] = b
	}
	return bubble
}

func createRepoAndSave(machine *bubblemachine.Machine) {
	repo, err := sqliterepository.NewMachineRepository(sqliterepository.NewRepositoryRequest{DataSourceName: "test.db", Init: true})
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()
	err = repo.Save(context.Background(), machine)

	if err != nil {
		log.Fatal(err)
	}
}
