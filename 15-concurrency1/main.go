// http://whipperstacker.com/2015/10/05/3-trivial-concurrency-exercises-for-the-confused-newbie-gopher/
// Exercise #1
package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Person struct {
	Name string
}

var Bob = Person{Name: "Bob"}
var Alice = Person{Name: "Alice"}

func main() {
	ch := make(chan string)

	fmt.Println("Let's go for a walk!")

	// Getting ready.
	go Bob.GetReady(ch)
	go Alice.GetReady(ch)

	for i := 0; i < 2; i++ {
		<-ch
	}

	// Arming alarm and putting shoes.
	go Alarm(ch)
	time.Sleep(1 * time.Second)
	go Bob.PutShoes(ch)
	go Alice.PutShoes(ch)

	for i := 0; i < 2; i++ {
		<-ch
	}

	fmt.Println("Exiting and locking the door.")

	<-ch
}

func Alarm(ch chan string) {
	fmt.Println("Arming alarm.")
	time.Sleep(5 * time.Second)
	fmt.Println("Alarm is counting down...")
	time.Sleep(20 * time.Second)
	fmt.Println("Alarm is armed!")
	ch <- "done"
}

func (p *Person) GetReady(ch chan string) {
	fmt.Println(p.Name + " started getting ready...")
	d := time.Duration(RandRange(10, 25)) * time.Second
	time.Sleep(d)
	fmt.Println(p.Name + " spent " + d.String() + " getting ready.")
	ch <- "done"
}

func (p *Person) PutShoes(ch chan string) {
	fmt.Println(p.Name + " started putting on shoes...")
	d := time.Duration(RandRange(10, 20)) * time.Second
	time.Sleep(d)
	fmt.Println(p.Name + " spent " + d.String() + " putting on shoes.")
	ch <- "done"
}

// Helpers
func RandRange(min, max int) int {
	return rand.Intn(max-min) + min
}
