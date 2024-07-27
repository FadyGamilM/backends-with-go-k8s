package main

import (
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
		GetCounter(server)
		// sleep from random time within 5 seconds
		// d := randomGenerator.Intn(5000)
		time.Sleep(time.Duration(10000) * time.Millisecond)
	}
}

func GetCounter(server string) error {
	counter := &Counter{}

	// we first read the response body, then docde it into our struct
	res, err := http.Get(server)
	if err != nil {
		log.Println("error sending request : ", err)
		return err
	}
	defer res.Body.Close()

	err = json.NewDecoder((res.Body)).Decode(counter)
	if err != nil {
		log.Println("error marshling request : ", err)
		return err
	}

	log.Println("reading counter : ", counter.Counter)
	return nil
}

type Counter struct {
	Counter int `json:"counter"`
}
