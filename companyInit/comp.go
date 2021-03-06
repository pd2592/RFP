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

func GetRfp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get some question for Hotels !!!!")
	err := r.ParseForm()
	commons.CheckErr(err)
	questionCategoryParent := r.FormValue("questionCategoryParent")
	travelAgencyMasterId := r.FormValue("travelAgencyMasterId")
	rfpId := r.FormValue("rfpId")

	//fmt.Fprintf(w, db.HotelEditResponse(questionCategoryParent, travelAgencyMasterId))
	fmt.Fprintf(w, db.CompanyEditRfp(rfpId, questionCategoryParent, travelAgencyMasterId))

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

func RfpView(w http.ResponseWriter, r *http.Request) {
	fmt.Println("View RFP ......")
	err := r.ParseForm()
	commons.CheckErr(err)
	RfpId := r.FormValue("rfpId")
	fmt.Fprintf(w, db.RfpFullView(RfpId))
}

func ListHotel(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Listing Hotels ......")
	// err := r.ParseForm()
	// commons.CheckErr(err)
	// cityId := r.FormValue("cityId")
	//fmt.Println(rfpId)
	//a := db.ListHotels(cityId)   //need to send array of cityId
	//fmt.Fprintln(w, a)

}

func RfpDrafted(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Listing Sent Rfp ......")
	err := r.ParseForm()
	commons.CheckErr(err)
	travelAgencyMasterId := r.FormValue("travelAgencyMasterId")
	//fmt.Println(rfpId)
	a := db.ListRfpDrafted(travelAgencyMasterId)
	fmt.Fprintln(w, a)
}

func RfpPublished(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Listing Sent Rfp ......")
	err := r.ParseForm()
	commons.CheckErr(err)
	travelAgencyMasterId := r.FormValue("travelAgencyMasterId")
	//fmt.Println(rfpId)
	a := db.ListRfpPublished(travelAgencyMasterId)
	fmt.Fprintln(w, a)
}

func RfpQuotes(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Listing Quotes Recieved for Rfp ......")
	err := r.ParseForm()
	commons.CheckErr(err)
	rfpId := r.FormValue("rfpId")
	//fmt.Println(rfpId)
	a := db.ListRfpQuotes(rfpId)
	fmt.Fprintln(w, a)
}

func DeclineQuote(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Decline Quote Recieved for Rfp ......")
	err := r.ParseForm()
	commons.CheckErr(err)
	rfpId := r.FormValue("rfpId")
	hotelId := r.FormValue("hotelId")

	//fmt.Println(rfpId)
	a := db.RejectQuotes(rfpId, hotelId)
	fmt.Fprintln(w, a)
}

func TrashRfp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Moving the rfp to trash ......")
	err := r.ParseForm()
	commons.CheckErr(err)
	rfpId := r.FormValue("rfpId")
	//fmt.Println(rfpId)
	a := db.TrashedRfp(rfpId)
	fmt.Fprintln(w, a)
}

func ShortListQuotes(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ShortListing Quotes Recieved ......")
	body, err := ioutil.ReadAll(r.Body)
	commons.CheckErr(err)
	RfpSent := commons.UnmarshalRFPSend(string(body))

	//fmt.Println(db.HotelResponse(RfpQues))
	fmt.Fprintln(w, db.RfpShortlist(RfpSent))
}

func UnShortListQuotes(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ShortListing Quotes Recieved ......")
	body, err := ioutil.ReadAll(r.Body)
	commons.CheckErr(err)
	RfpSent := commons.UnmarshalRFPSend(string(body))

	//fmt.Println(db.HotelResponse(RfpQues))
	fmt.Fprintln(w, db.RfpUnShortlist(RfpSent))
}

func RfpHotelResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hotels Response for RFP ......")
	err := r.ParseForm()
	commons.CheckErr(err)
	RfpId := r.FormValue("rfpId")
	HotelId := r.FormValue("hotelId")

	fmt.Fprintf(w, db.GetRfpResponse(RfpId, HotelId))
}

func AcceptQuote(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hotels Response for RFP ......")
	err := r.ParseForm()
	commons.CheckErr(err)
	RfpId := r.FormValue("rfpId")
	HotelId := r.FormValue("hotelId")

	fmt.Fprintf(w, db.AcceptQuote(RfpId, HotelId))
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
