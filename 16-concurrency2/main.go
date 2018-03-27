// http://whipperstacker.com/2015/10/05/3-trivial-concurrency-exercises-for-the-confused-newbie-gopher/
package main

import (
	"fmt"
	"strconv"
	"time"
	"math/rand"
)

/*
Has 8 computers
First come first serve
25 people total
Usage time 15m - 2h
 */

type Tourist struct {
	Name string
}

var Tourists []Tourist

func init() {
	genTourists()
}

func main() {
	done := make(chan int, len(Tourists))
	bc := 0
	
	for _, t := range Tourists {
		ct := t
		if bc < 8 {
			go ct.UsePC(done)
		} else {
			<-done
			go ct.UsePC(done)
		}
		bc++
	}
	
	// Wait for the last 8 to finish.
	for i := 0; i < 8; i++ {
		<-done
	}
}

func (t *Tourist) UsePC(c chan int) {
	fmt.Printf("Tourist %s is online.\n", t.Name)
	d := time.Duration(randRange(5, 10)) * time.Second
	time.Sleep(d)
	fmt.Printf("Tourist %s is done - spent %s online.\n", t.Name, d)
	c <- 1
}

// Helpers
func genTourists() {
	for i := 1; i <= 25; i++ {
		Tourists = append(Tourists, Tourist{strconv.Itoa(i)})
	}
}

// Seeding the math/rand
var seed = rand.NewSource(time.Now().UnixNano())
var r = rand.New(seed)

func randRange(min, max int) int {
	return r.Intn(max-min) + min
}
