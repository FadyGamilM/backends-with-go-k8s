package main

import (
	"bytes"
	"counterpooler/models"
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
		d := randomGenerator.Intn(5000)
		time.Sleep(time.Duration(d) * time.Millisecond)
	}
}

func SetCounter(server string) error {
	counter := &models.Counter{
		Counter: randomGenerator.Intn(10000),
	}

	data, err := json.Marshal(counter)
	if err != nil {
		log.Println("error marshling request : ", err)
		return err
	}

	res, err := http.Post(server, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println("error sending request : ", err)
		return err
	}

	log.Println(res.Body)
	return nil
}
