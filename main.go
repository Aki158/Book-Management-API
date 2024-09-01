package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	hello := []byte("Hello My World!!!")
	_, err := w.Write(hello)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	fmt.Println("Server Start Up........")

	// httpsサーバーを起動する
	http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil)
}