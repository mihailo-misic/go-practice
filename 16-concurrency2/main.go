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
	Done bool
}

var Tourists []Tourist

var seed = rand.NewSource(time.Now().UnixNano())
var r = rand.New(seed)

var bc = 0

func init() {
	genTourists()
}

func main() {
	busy := make(chan int, len(Tourists))
	busy <- 0

	for _, t := range Tourists {
		if bc < 8 {
			go t.UsePC(busy)
			<-busy
			<-busy
		}
	}
}

func (t *Tourist) UsePC(c chan int) {
	fmt.Printf("Tourist %s is online.\n", t.Name)
	bc++
	c <- bc
	d := time.Duration(randRange(5, 10)) * time.Second
	time.Sleep(d)
	fmt.Printf("Tourist %s is done - spent %s online.\n", t.Name, d)
	bc--
	c <- bc
}

// Helpers
func genTourists() {
	for i := 1; i <= 25; i++ {
		Tourists = append(Tourists, Tourist{strconv.Itoa(i), false})
	}
}

func randRange(min, max int) int {
	return r.Intn(max-min) + min
}
