package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter,req *http.Request){
	fmt.Println("Hello, Backend is Running");
}

func main() {
	http.HandleFunc("/hello",hello)

	fmt.Println("Backend Running")
	http.ListenAndServe(":3001",nil)
}