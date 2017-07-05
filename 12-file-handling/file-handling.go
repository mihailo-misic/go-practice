package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	CreateDir("source")
	file, err := os.Create("source/data.txt")
	PrintFatalError(err)
	fmt.Println("File data has been created.")
	defer file.Close()

	wb := bufio.NewWriter(file)
	for i := 1; i <= 1000; i++ {
		wb.WriteString(fmt.Sprintf("%v) User connected\n", i))
	}
	wb.Flush()

	CopyFile("source/data.txt", "dest/data.txt")
}

func CopyFile(fname1, fname2 string) {
	CreateDir("dest")

	fOld, err := os.Open(fname1)
	PrintFatalError(err)
	defer fOld.Close()

	fNew, err := os.Create(fname2)
	PrintFatalError(err)
	defer fNew.Close()

	_, err = io.Copy(fNew, fOld)
	PrintFatalError(err)

	err = fNew.Sync()
	PrintFatalError(err)
}

func PrintFatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CreateDir(dirname string) {
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		os.MkdirAll(dirname, os.ModePerm)
	}
}
