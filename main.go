package main

import (
	"fmt"
	"os"
)

var filename = "concat.txt"

func main() {
	var n int
	var s []byte

	fmt.Print(">> Enter lines number: ")
	_, err := fmt.Scan(&n)
	if err != nil {
		panic(err)
	}

	bytesChan := make(chan []byte, n)
	errChan := make(chan error)
	doneChan := make(chan struct{})

	go saveText(errChan, bytesChan, doneChan)

	err = <-errChan
	if err != nil {
		fmt.Println("an error occurred:", err)
		return
	}

	fmt.Println("Ready for lines")
	for i := 1; i <= n; i++ {
		fmt.Printf(">> Enter a line %d: ", i)
		_, err := fmt.Scan(&s)
		if err != nil {
			panic(err)
		}
		bytesChan <- s
	}
	close(bytesChan)

	<-doneChan
	fmt.Println("Done")
}

func saveText(err chan<- error, bytes <-chan []byte, done chan<- struct{}) {
	defer close(err)
	defer close(done)

	fmt.Printf("Creating the file %q. Please, wait...\n", filename)
	f, openErr := os.Create(filename)
	if openErr != nil {
		err <- openErr
		return // явно выходим из функции, чтобы последующий код не выполнился
	}

	for s := range bytes {
		_, err := f.Write(s)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("Closing the file %q...\n", filename)
	errClose := f.Close()
	if errClose != nil {
		panic(errClose)
	}

	done <- struct{}{}
}
