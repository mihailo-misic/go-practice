package main

import (
	"github.com/sirupsen/logrus"
	"github.com/mihailomisic/go-practice/14-logging/log"
)

func main() {
	pop()

	logrus.Error("main")
}

func pop() {
	// Init logging
	defer log.Start().Close()

	logrus.Error("pop")
}

