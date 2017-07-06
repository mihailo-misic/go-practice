package main

import (
	"bufio" // For buffered reading & writing
	"fmt"
	"io" // For Copying
	"log"
	"os" // Creating, Opening, Renaming, Removing, and setting permissions
	"time"
)

// This program covers the following:
// Create
// Read
// Write (Append || Replace)
// Move/Rename
// Copy
// Status
// Watch the file
// Delete

func main() {
	srcFile := "source/data.txt"
	dstFile := "dest/data.txt"

	CreateDir("source")
	file, err := os.Create(srcFile)
	PrintFatalError(err)
	fmt.Println("File data has been created.")
	defer file.Close()

	wb := bufio.NewWriter(file)
	for i := 1; i <= 10; i++ {
		wb.WriteString(fmt.Sprintf("%v) User connected\n", i))
	}
	wb.Flush()

	CopyFile(srcFile, dstFile)

	file2, err := os.OpenFile(dstFile, os.O_APPEND|os.O_RDWR, 0644)
	PrintFatalError(err)
	defer file2.Close()

	ReadFile(file2)

	go WatchFile(dstFile)

	for i := 0; i < 10; i++ {
		WriteToFile(file2)
		time.Sleep(1 * time.Second)
	}

	ReadFile(file2)

	removeFile(srcFile)
	renameFile(dstFile, "./sacred.txt")
	removeFile(dstFile)
}

func CreateDir(dirName string) {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		os.MkdirAll(dirName, os.ModePerm)
	}
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

func ReadFile(file *os.File) {
	scanner := bufio.NewScanner(file)
	for count := 1; scanner.Scan(); count++ {
		fmt.Printf("Reading line (%v): %v\n", count, scanner.Text())
	}
	file.Seek(0, 0) // Reset the pointer to the beginning
}

func WatchFile(fname string) {
	fileStat, err := os.Stat(fname)
	PrintFatalError(err)
	for {
		time.Sleep(1 / 2 * time.Second)
		fileStat2, err := os.Stat(fname)
		PrintFatalError(err)
		if fileStat.ModTime() != fileStat2.ModTime() {
			fmt.Println("File was modified at", fileStat2.ModTime())
			fileStat, err = os.Stat(fname)
			PrintFatalError(err)
		}
	}
}

func WriteToFile(file *os.File) {
	writeBuffer := bufio.NewWriter(file)
	for i := 1; i <= 5; i++ {
		writeBuffer.WriteString(fmt.Sprintf("Added line %v\n", i))
	}
	writeBuffer.Flush()
	file.Seek(0, 0) // Reset the pointer to the beginning
}

func removeFile(dirName string) {
	err := os.RemoveAll(dirName)
	PrintFatalError(err)
}

func renameFile(oldPath, newPath string) {
	err := os.Rename(oldPath, newPath)
	PrintFatalError(err)
	fmt.Println("File from", oldPath, "moved to", newPath, "here are its stats:")
	FileStat(newPath)
}

func FileStat(fname string) {
	fileStats, err := os.Stat(fname)
	PrintFatalError(err)
	fmt.Println("File name:", fileStats.Name())
	fmt.Println("Is a dir:", fileStats.IsDir())
	fmt.Println("Permissions:", fileStats.Mode())
	fmt.Println("File size:", fileStats.Size())
	fmt.Println("Last time the file modified:", fileStats.ModTime())
}

func PrintFatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
