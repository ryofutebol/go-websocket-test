package main

import (
	"fmt"
	"net/http"
)

func hd(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World Go")
}

func main() {
	http.HandleFunc("/", hd)
	http.ListenAndServe(":3000", nil)
}
