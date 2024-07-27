package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var (
	serverUrl       = "http://localhost:8081/counter"
	randomGenerator = rand.New(rand.NewSource(time.Now().UnixMilli()))
)

func main() {
	server := os.Getenv("SERVER_URL")
	if server == "" {
		server = serverUrl
	}

	for {
		// send the post request
		SetCounter(server)
		// sleep from random time within 5 seconds
		// d := randomGenerator.Intn(5000)
		time.Sleep(time.Duration(10000) * time.Millisecond)
	}
}

func SetCounter(server string) error {
	counter := &Counter{
		Counter: randomGenerator.Intn(10000),
	}

	log.Println("setting counter as : ", counter.Counter)

	data, err := json.Marshal(counter)
	if err != nil {
		log.Println("error marshling request : ", err)
		return err
	}

	_, err = http.Post(server, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println("error sending request : ", err)
		return err
	}

	return nil
}

type Counter struct {
	Counter int `json:"counter"`
}
