package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc(`/`, func(w http.ResponseWriter, r *http.Request) {
		log.Println("wait a couple of seconds ...")
		// time.Sleep(time.Duration(rand.IntN(100)) * time.Millisecond)
		io.WriteString(w, `Hi`)
		log.Println("Done.")
	})
	log.Println(http.ListenAndServe(":8080", nil))
}
