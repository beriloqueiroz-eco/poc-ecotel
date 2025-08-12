package main

import "net/http"

// requests
func main() {
	go http.DefaultClient.Get("http://localhost:8080/hello")
	go http.DefaultClient.Get("http://localhost:8080/hello")
	for i := 0; i < 10; i++ {
		http.DefaultClient.Get("http://localhost:8080/hello")
	}
}
