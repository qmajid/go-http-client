package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	WITH_DEFAULT_HTTP_CLIENT_SETTING bool = false
	//-------------
	NUMBER_OF_CONNECTIONS int = 1
	ALL_MESSAGES_COUNT    int = 1000
	CLINT_TIMEOUT         int = 2000 //ms
)

var (
	ERROR_COUNT int
	mu          sync.Mutex
)

func main() {
	log.Println("HTTP GET")
	run()
}

func run() {
	var wg sync.WaitGroup
	message_per_connection := ALL_MESSAGES_COUNT / NUMBER_OF_CONNECTIONS
	for range NUMBER_OF_CONNECTIONS {
		wg.Add(1)
		go createNewConnection(message_per_connection, &wg)
	}
	wg.Wait()
	log.Printf("scenario_2 error count: %v\n", ERROR_COUNT)
}

func createNewConnection(messageCount int, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	var client *http.Client

	if WITH_DEFAULT_HTTP_CLIENT_SETTING {
		client = &http.Client{
			Timeout: time.Duration(CLINT_TIMEOUT) * time.Millisecond,
		}
	} else {
		t := http.DefaultTransport.(*http.Transport).Clone()
		t.MaxIdleConns = 100
		t.MaxConnsPerHost = 100
		t.MaxIdleConnsPerHost = 100

		client = &http.Client{
			Timeout:   time.Duration(CLINT_TIMEOUT) * time.Millisecond,
			Transport: t,
		}
	}

	var wg sync.WaitGroup
	for i := range messageCount {
		wg.Add(1)
		go send(client, &wg, i)
	}
	wg.Wait()
}

func send(client *http.Client, wg *sync.WaitGroup, counter int) {
	defer wg.Done()
	start_time := time.Now()
	r, err := client.Get(`http://127.0.0.1:8080/`)
	fmt.Printf("elapsed time: %v\n", time.Since(start_time))
	if err != nil {
		mu.Lock()
		// log.Printf("Get error[%v][%v]\n", counter, err)
		ERROR_COUNT++
		mu.Unlock()
		return
	}
	defer r.Body.Close()
	_, err = ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	// log.Printf("body[%v][%v]\n", counter, string(body))
}
