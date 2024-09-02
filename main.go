package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Aki158/School-API/pkg/db"
	"github.com/Aki158/School-API/pkg/handlers"
)

func main() {
	// MySQLに接続する
	mydb := &db.Database{}
	for i := 0; i < 10; i++ {
		mydb.Connect()
		err := mydb.UseDb.Ping()
		if err == nil {
			log.Println("Successfully connected to the database")
			break
		}
		log.Println("Waiting for the database to be ready...")
		// 5秒待機して再試行する
		time.Sleep(5 * time.Second)
	}
	defer mydb.UseDb.Close()

	http.Handle("/students", http.HandlerFunc(handlers.StudentsHandler(mydb)))

	log.Println("Server start up")

	// httpsサーバーを起動する
	http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil)
}