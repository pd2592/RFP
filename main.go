package main

import (
	"log"
	"net/http"

	"./hotelInit" //importing hotelInit package
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/quesReqHotel", hotelInit.RequestHotelQues)
	//hotelInit.SayHi("atddds")
	log.Fatal(http.ListenAndServe(":8080", router))

}