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
	router.HandleFunc("/ansHotel", hotelInit.ResponseHotelAns)
	router.HandleFunc("/ansHotel/edit", hotelInit.EditHotelAns)
	router.HandleFunc("/category/get", hotelInit.GetParentCategory)

	router.HandleFunc("/rfp/recieved", hotelInit.RfpRecieved) //rfp recived by hotel

	router.HandleFunc("/basic/list", companyInit.ListBasic)
	router.HandleFunc("/basic/ans", companyInit.RfpBasic)
	router.HandleFunc("/rfp", companyInit.RfpQuestion)
	router.HandleFunc("/rfp/edit", companyInit.RfpEdit)
	router.HandleFunc("/rfp/show", companyInit.RfpPreview)
	router.HandleFunc("/rfp/send", companyInit.RfpSend)

	router.HandleFunc("/hotel/list", companyInit.ListHotel)
	router.HandleFunc("/rfp/published", companyInit.RfpPublished)

	//hotelInit.SayHi("atddds")
	log.Fatal(http.ListenAndServe(":9000", router))

}
