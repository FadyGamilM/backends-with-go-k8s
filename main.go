package main

import (
	"errors"
	"log"
	"runtime"
)

func main() {

	// go func() {
	// 	log.Println("hello")
	// }()

	// select {}

	// log.Println("done")

	// var x X
	// x = Y{}
	// r := Y{}

	// log.Println(x == r)

	// y := Y{}
	// var x X = y
	// err := x.Do()
	// log.Println(err)

	c := make(chan int, 3)
	go func(c chan int) {
		for i := 0; i < 4; i++ {
			num := <-c
			log.Println(num * num)
		}
	}(c)

	log.Println(runtime.NumGoroutine())
	c <- 1
	c <- 2
	c <- 3
	c <- 4
	log.Println(runtime.NumGoroutine())

	x := interface{}(&Y{})
	y := interface{}(&Y{})
	log.Println(x == y)
	a := []string{"1"}
	a = nil
	log.Println(a, len(a), cap(a))

	b := []int{1, 2, 3}
	b = b[:0]
	log.Println(b, len(b), cap(b))

}

type X interface {
	Do() error
}

type Y struct {
}

func (y Y) Do() error {
	return errors.New("type")
}
