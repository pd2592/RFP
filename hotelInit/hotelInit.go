package hotelInit

import (
	"fmt"
	"net/http"

	commons "../commons"
	db "../db" //importing db package
)

// func SayHi(ats string) {
// 	fmt.Println(ats)
// }

func RequestHotelQues(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get some question based on questionCategoryParent !!!!")
	err := r.ParseForm()
	commons.CheckErr(err)
	questionCategoryParent := r.FormValue("questionCategoryParent")
	db.QuestionJsonByCat(questionCategoryParent)
	fmt.Fprintf(w, questionCategoryParent)
}

// body, err := ioutil.ReadAll(r.Body)
// commons.CheckErr(err)
// Question := commons.UnmarshalQuestion(string(body))
// fmt.Println(Question.QuesCategory[0])
