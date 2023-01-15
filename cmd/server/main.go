package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	metricsHandler := func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			io.WriteString(w, "Method Not Allowed\n")
			return
		}

		fmt.Println(req.URL)
	}

	http.HandleFunc("/", metricsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
