package main

import (
	"UTS/Controllers"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/rooms", Controllers.GetAllRooms).Methods("GET")
	router.HandleFunc("/rooms", Controllers.InsertRoom).Methods("POST")
	router.HandleFunc("/rooms/{room_id}", Controllers.GetRoomDetails).Methods("DELETE")
	router.HandleFunc("/rooms/{room_id}", Controllers.GetRoomDetails).Methods("GET")

	http.Handle("/", router)
	fmt.Println("Connected to port 8888")
	log.Println("Connected to port 8888")
	log.Fatal(http.ListenAndServe(":8888", router))
}
