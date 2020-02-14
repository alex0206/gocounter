package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

const maxHandlers = 5

func main() {
	var wg sync.WaitGroup
	scanner := bufio.NewScanner(os.Stdin)
	results := make(chan int)
	handlers := make(chan struct{}, maxHandlers)

	for scanner.Scan() {
		handlers <- struct{}{}
		wg.Add(1)
		go handle(scanner.Text(), &wg, results, handlers)
	}

	go func() {
		wg.Wait()
		close(results)
		close(handlers)
	}()

	fmt.Printf("Total: %d\n", countTotal(results))
}

func countTotal(results chan int) (total int) {
	for result := range results {
		total += result
	}
	return
}

func handle(url string, wg *sync.WaitGroup, results chan<- int, handlers <-chan struct{}) {
	defer wg.Done()

	wordCnt := countGoWords(url)

	fmt.Printf("Count for %s: %d\n", url, wordCnt)

	<-handlers
	results <- wordCnt
}

func countGoWords(url string) (goCnt int) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return bytes.Count(body, []byte("Go"))
}
