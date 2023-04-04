package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
	"L0/models"
	stan "github.com/nats-io/stan.go"
	"L0/pkg/handler"
	repo "L0/pkg/repository"
)

func main() {
	repository := repo.NewRepo()
	handler := handler.NewHandler(*repository)

	handler.InitRouters()

	fmt.Println("[LOG]: Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
	
	go startNats(repository)
}

func SubscriberNats(repo *repo.Repository, conn stan.Conn) {
	var err error

	_, err = conn.Subscribe("NewOrder", func(msg *stan.Msg) {

		var ord models.Order
		if err = json.Unmarshal(msg.Data, &ord); err != nil {
			log.Println(err)
			return
		}
		repo.InsertOrder(&ord)

		fmt.Printf("seq = %d [redelivered = %v] mes= %s \n", msg.Sequence, msg.Redelivered, msg.Data)

		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

		msg.Ack()

	}, stan.DurableName("i-will-remember"), stan.MaxInflight(100), stan.SetManualAckMode())

	if err != nil {
		log.Println(err)
	}
}

func startNats(repo *repo.Repository) {
	if err := runNats(repo); err != nil {
		log.Println(err)
	}
}

func runNats(repo *repo.Repository) error {
	conn, err := stan.Connect(
		"test-cluster",
		"test-client",
	)
	checkFail("Connect NATS Streaming", err)

	fmt.Println("nats err:", err.Error())

	if err != nil {
		return err
	}

	done := make(chan struct{})
	time.Sleep(time.Duration(rand.Intn(4000)) * time.Millisecond)
	SubscriberNats(repo, conn)
	<-done

	return nil
}

func checkFail(funcname string, err error) {
	if err != nil {
		log.Println(err)
	} else {
		log.Println(funcname + ": OK")
	}
}
