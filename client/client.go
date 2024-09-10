package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	WITH_DEFAULT_HTTP_CLIENT_SETTING bool = false
	CLINT_TIMEOUT                    int  = 10
	ALL_MESSAGES_COUNT               int  = 200
	NUMBER_OF_CONNECTIONS            int  = 2
)

func main() {
	log.Println("HTTP GET")

	senario1()
	log.Println("------------------------------------------")
	// senario2()
}

func senario1() {
	var wg sync.WaitGroup
	for range 1 {
		wg.Add(1)
		go createNewConnection(ALL_MESSAGES_COUNT, &wg)
	}
	wg.Wait()
}

func senario2() {
	var wg sync.WaitGroup
	message_per_connection := ALL_MESSAGES_COUNT / NUMBER_OF_CONNECTIONS
	for range NUMBER_OF_CONNECTIONS {
		wg.Add(1)
		go createNewConnection(message_per_connection, &wg)
	}
	wg.Wait()
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
	r, err := client.Get(`http://127.0.0.1:8080/`)
	if err != nil {
		log.Printf("Get error[%v][%v]\n", counter, err)
		return
	}
	defer r.Body.Close()
	_, err = ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	// log.Printf("body[%v][%v]\n", counter, string(body))
}
