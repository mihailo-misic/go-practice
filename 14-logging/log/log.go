package log

import (
	"os"
	"github.com/sirupsen/logrus"
)

func Start() *os.File {
	// Create your file with desired read/write permissions
	f, err := os.OpenFile("./go.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	// Configure logrus
	logrus.SetOutput(f)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	return f
}
