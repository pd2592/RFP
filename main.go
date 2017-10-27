package main

import (
	"log"
	"net/http"

	"./companyInit"
	"./hotelInit" //importing hotelInit package
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/quesReqHotel", hotelInit.RequestHotelQues)
	router.HandleFunc("/quesReqBySubCat", hotelInit.QuesBySubCat)
	router.HandleFunc("/ansResHotel", hotelInit.ResponseHotelAns)

	router.HandleFunc("/rfp", companyInit.RfpQuestion)
	router.HandleFunc("/rfp/edit", companyInit.RfpEdit)
	router.HandleFunc("/rfp/show", companyInit.RfpPreview)
	router.HandleFunc("/rfp/send", companyInit.RfpSend)

	router.HandleFunc("/hotel/list", companyInit.ListHotels)

	router.HandleFunc("/category/get", hotelInit.GetParentCategory)

	//hotelInit.SayHi("atddds")
	log.Fatal(http.ListenAndServe(":8080", router))

}
