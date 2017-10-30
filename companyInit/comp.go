package companyInit

import (
	"fmt"
	"io/ioutil"
	"net/http"

	commons "../commons"
	db "../db" //importing db package
)

func RfpQuestion(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	commons.CheckErr(err)
	RfpQues := commons.UnmarshalRFP(string(body))

	//fmt.Println(db.HotelResponse(RfpQues))
	fmt.Println(db.RfpRequest(RfpQues))
}

func RfpEdit(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	commons.CheckErr(err)
	RfpEdits := commons.UnmarshalRFP(string(body))

	//fmt.Println(db.HotelResponse(RfpQues))
	fmt.Fprintln(w, db.RfpEditor(RfpEdits))
}

func RfpPreview(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Previewing RFP ......")
	err := r.ParseForm()
	commons.CheckErr(err)
	RfpId := r.FormValue("rfpId")
	fmt.Fprintf(w, db.GetRfp(RfpId))
}

func RfpSend(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Sending RFP ......")

	body, err := ioutil.ReadAll(r.Body)
	commons.CheckErr(err)
	RfpSent := commons.UnmarshalRFPSend(string(body))

	//fmt.Println(db.HotelResponse(RfpQues))
	fmt.Fprintln(w, db.RfpSend(RfpSent))
}

func ListHotel(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Listing Hotels ......")
	err := r.ParseForm()
	commons.CheckErr(err)
	cityId := r.FormValue("cityId")
	//fmt.Println(rfpId)
	a := db.ListHotels(cityId)
	fmt.Fprintln(w, a)

}

func ListBasic(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Listing Basic question ......")
	fmt.Fprintln(w, db.GetBasicList())
}

func RfpBasic(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Recieving RFP basics ......")

	body, err := ioutil.ReadAll(r.Body)
	commons.CheckErr(err)
	RfpBasicQ := commons.UnmarshalRFPBasic(string(body))

	//fmt.Println(db.HotelResponse(RfpQues))
	fmt.Fprintln(w, db.RfpBasicAns(RfpBasicQ))
}
