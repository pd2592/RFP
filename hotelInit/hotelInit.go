package hotelInit

import (
	"fmt"
	"io/ioutil"
	"net/http"

	commons "../commons"
	db "../db" //importing db package
)

// func SayHi(ats string) {
// 	fmt.Println(ats)
// }
func GetParentCategory(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, db.ParentCategory())
}

func RequestHotelQues(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get some question for Hotels !!!!")
	err := r.ParseForm()
	commons.CheckErr(err)
	questionCategoryParent := r.FormValue("questionCategoryParent")

	fmt.Fprintf(w, db.QuestionJsonByCat(questionCategoryParent))
}

func QuesBySubCat(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Questions by subcategory !!!!")
	err := r.ParseForm()
	commons.CheckErr(err)
	questionCategory := r.FormValue("questionSubCategory")
	fmt.Fprintf(w, db.QuestionJsonBySubCat(questionCategory))

}

func ResponseHotelAns(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	commons.CheckErr(err)
	Response := commons.UnmarshalResponse(string(body))

	fmt.Println(db.HotelResponse(Response))
}

func EditHotelAns(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get some question for Hotels !!!!")
	err := r.ParseForm()
	commons.CheckErr(err)
	questionCategoryParent := r.FormValue("questionCategoryParent")
	travelAgencyMasterId := r.FormValue("travelAgencyMasterId")

	fmt.Fprintf(w, db.HotelEditResponse(questionCategoryParent, travelAgencyMasterId))
}

// body, err := ioutil.ReadAll(r.Body)
// commons.CheckErr(err)
// Question := commons.UnmarshalQuestion(string(body))
// fmt.Println(Question.QuesCategory[0])
