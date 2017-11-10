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
	router.HandleFunc("/quesReqHotel", hotelInit.RequestHotelQues) //listing initial question for hotels
	router.HandleFunc("/quesReqBySubCat", hotelInit.QuesBySubCat)
	router.HandleFunc("/ansHotel", hotelInit.ResponseHotelAns)  //storing hotel initial answers
	router.HandleFunc("/ansHotel/edit", hotelInit.EditHotelAns) //editing hotels answers
	router.HandleFunc("/category/get", hotelInit.GetParentCategory)

	router.HandleFunc("/rfp/recieved", hotelInit.RfpRecieved)      //rfp recived by hotel
	router.HandleFunc("/slab/list", hotelInit.ListSlab)            //list slab available for that hotel
	router.HandleFunc("/slab/assign", hotelInit.AssignSlab)        //assign one slab for one rfp
	router.HandleFunc("/hotel/getDetails", hotelInit.HotelDetails) //get hotel details for review

	router.HandleFunc("/basic/list", companyInit.ListBasic) //listing basic question, common for all corporates
	router.HandleFunc("/basic/ans", companyInit.RfpBasic)   //answering basic ques and also creating rfp
	router.HandleFunc("/rfp", companyInit.RfpQuestion)      //creating rfp(not in use)
	router.HandleFunc("/rfp/edit", companyInit.RfpEdit)     // edting rfp in case user returns and add more question
	router.HandleFunc("/rfp/get", companyInit.GetRfp)       // getting all question category wise along with the choices that corporate made. For editing purpose
	router.HandleFunc("/rfp/show", companyInit.RfpPreview)  //Rfp Preview
	router.HandleFunc("/rfp/send", companyInit.RfpSend)     //Sending the rfp to multiple hotel
	router.HandleFunc("/rfp/preview", companyInit.RfpView)  //rfp preview for companies and hotels

	router.HandleFunc("/hotel/list", companyInit.ListHotel)          // list hotels to which rfp should be send
	router.HandleFunc("/rfp/published", companyInit.RfpPublished)    //listing rfp which are sent by corporate
	router.HandleFunc("/rfp/quotes", companyInit.RfpQuotes)          // listing rfps for which quotes are recieved
	router.HandleFunc("/rfp/hotelRes", companyInit.RfpHotelResponse) //hotel response for particular rfp
	router.HandleFunc("/rfp/acceptQuote", companyInit.AcceptQuote)

	//hotelInit.SayHi("atddds")
	log.Fatal(http.ListenAndServe(":9000", router))

}
