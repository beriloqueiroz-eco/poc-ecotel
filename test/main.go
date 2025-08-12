package main

import "net/http"

// requests
func main() {
	go http.DefaultClient.Get("http://localhost:8081/hello")
	go http.DefaultClient.Get("http://localhost:8081/hello")
	go http.DefaultClient.Get("http://localhost:8082/hello")
	for i := 0; i < 4; i++ {
		http.DefaultClient.Get("http://localhost:8083/hello")
	}
}
