package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc(`/`, func(w http.ResponseWriter, r *http.Request) {
		log.Println("wait a couple of seconds ...")
		time.Sleep(10 * time.Millisecond)
		io.WriteString(w, `Hi`)
		log.Println("Done.")
	})
	log.Println(http.ListenAndServe(":8080", nil))
}
