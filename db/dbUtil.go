package db

import (
	"database/sql"
	"encoding/json"
	"fmt"

	commons "../commons"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

//getting question from db category wise

func QuestionJsonByCat(parentId string) string {
	db = GetDB()
	var parentCatVar commons.ParentCat
	var questionCatVar commons.QuestionCat
	var questionCatVars []commons.QuestionCat
	var quesMVar commons.QuesM
	//var quesMVars []commons.QuesM
	var ansMVar commons.AnsM
	//var ansMVars []commons.AnsM

	retr_stmt, err := db.Query("Select qm.questionCategoryParentId, qm.questionCategoryId, qq.questionCategory from questionmaster as qm JOIN questioncategory as qq ON qm.questionCategoryId = qq.questionCategoryId where qq.parentId = '" + parentId + "' GROUP BY qq.questionCategoryId")
	commons.CheckErr(err)
	for retr_stmt.Next() {
		err := retr_stmt.Scan(&parentCatVar.QuestionCategoryParentId, &questionCatVar.QuestionCategoryId, &questionCatVar.QuestionCategory)
		commons.CheckErr(err)

		var quesMVars []commons.QuesM

		retr_stmt1, err := db.Query("Select questionMasterId, questionText, groupQuestionId, questionSubTypeId, isMandatory from questionmaster where questionCategoryId = '" + questionCatVar.QuestionCategoryId + "'")
		for retr_stmt1.Next() {
			err := retr_stmt1.Scan(&quesMVar.QuestionId, &quesMVar.QuestionText, &quesMVar.GroupQuestionId, &quesMVar.QuestionSubTypeId, &quesMVar.IsMandatory)
			commons.CheckErr(err)

			var ansMVars []commons.AnsM

			retr_stmt2, err := db.Query("Select answerMasterId, answerText from answermaster where questionMasterId = '" + quesMVar.QuestionId + "'")
			for retr_stmt2.Next() {
				err := retr_stmt2.Scan(&ansMVar.AnswerId, &ansMVar.AnswerText)
				commons.CheckErr(err)
				ansMVar = commons.AnsM{
					AnswerId:   ansMVar.AnswerId,
					AnswerText: ansMVar.AnswerText,
				}
				ansMVars = append(ansMVars, ansMVar)
			}
			quesMVar = commons.QuesM{
				QuestionId:        quesMVar.QuestionId,
				QuestionText:      quesMVar.QuestionText,
				QuestionSubTypeId: quesMVar.QuestionSubTypeId,
				GroupQuestionId:   quesMVar.GroupQuestionId,
				IsMandatory:       quesMVar.IsMandatory,
				ConcatAns:         ansMVars,
			}
			quesMVars = append(quesMVars, quesMVar)

		}
		questionCatVar = commons.QuestionCat{
			QuestionCategoryId: questionCatVar.QuestionCategoryId,
			QuestionCategory:   questionCatVar.QuestionCategory,
			Ques:               quesMVars,
		}
		questionCatVars = append(questionCatVars, questionCatVar)

	}
	parentCatVar = commons.ParentCat{
		QuestionCategoryParentId: parentId,
		QuestionCategoryParent:   "test",
		QuesCategory:             questionCatVars,
	}
	b, err := json.Marshal(parentCatVar)
	commons.CheckErr(err)
	//	fmt.Println(string(b))
	return string(b)
}

//getting question from db sub category wise

func QuestionJsonBySubCat(questionCategory string) string {

	db = GetDB()
	var ParentSubCatVar commons.ParentSubCat
	var quesMVar commons.QuesM
	var ansMVar commons.AnsM
	var groupQuestionVar commons.GroupQuestion
	var quesMVars []commons.QuesM

	err := db.QueryRow("Select questionCategoryParentId from questionmaster where questionCategoryId = '" + questionCategory + "'").Scan(&ParentSubCatVar.QuestionCategoryParentId)
	commons.CheckErr(err)
	retr_stmt, err := db.Query("Select questionMasterId, questionText, questionSubTypeId, groupQuestionId, connectedQuestionId, answerMasterId, isMandatory from questionmaster where questionCategoryId = '" + questionCategory + "'")
	for retr_stmt.Next() {

		err := retr_stmt.Scan(&quesMVar.QuestionId, &quesMVar.QuestionText, &quesMVar.QuestionSubTypeId, &quesMVar.GroupQuestionId, &quesMVar.ConnectedQuestionId, &quesMVar.AnswerMasterId, &quesMVar.IsMandatory)
		commons.CheckErr(err)

		var ansMVars []commons.AnsM
		var groupQuestionVars []commons.GroupQuestion

		retr_stmt2, err := db.Query("Select answerMasterId, answerText from answermaster where questionMasterId = '" + quesMVar.QuestionId + "'")
		for retr_stmt2.Next() {
			err := retr_stmt2.Scan(&ansMVar.AnswerId, &ansMVar.AnswerText)
			commons.CheckErr(err)
			ansMVar = commons.AnsM{
				AnswerId:   ansMVar.AnswerId,
				AnswerText: ansMVar.AnswerText,
			}
			ansMVars = append(ansMVars, ansMVar)
		}
		if quesMVar.QuestionSubTypeId == "10" && quesMVar.GroupQuestionId != "0" {
			//get group questions here

			retr_stmt3, err := db.Query("Select groupQuestionMasterId, groupQuestionId, questionText, questionSubTypeId from groupquestion where groupQuestionId = '" + quesMVar.GroupQuestionId + "'")
			commons.CheckErr(err)

			for retr_stmt3.Next() {
				err := retr_stmt3.Scan(&groupQuestionVar.GroupQuestionMasterId, &groupQuestionVar.GroupQuestionId, &groupQuestionVar.QuestionText, &groupQuestionVar.QuestionSubTypeId)
				commons.CheckErr(err)
				groupQuestionVar = commons.GroupQuestion{
					GroupQuestionMasterId: groupQuestionVar.GroupQuestionMasterId,
					GroupQuestionId:       groupQuestionVar.GroupQuestionId,
					QuestionText:          groupQuestionVar.QuestionText,
					QuestionSubTypeId:     groupQuestionVar.QuestionSubTypeId,
				}
				groupQuestionVars = append(groupQuestionVars, groupQuestionVar)
			}
		}
		quesMVar = commons.QuesM{
			QuestionId:          quesMVar.QuestionId,
			QuestionText:        quesMVar.QuestionText,
			QuestionSubTypeId:   quesMVar.QuestionSubTypeId,
			GroupQuestionId:     quesMVar.GroupQuestionId,
			ConnectedQuestionId: quesMVar.ConnectedQuestionId,
			AnswerMasterId:      quesMVar.AnswerMasterId,
			IsMandatory:         quesMVar.IsMandatory,
			TabColumn:           groupQuestionVars,
			ConcatAns:           ansMVars,
		}
		quesMVars = append(quesMVars, quesMVar)
	}

	ParentSubCatVar = commons.ParentSubCat{
		QuestionCategoryParentId: ParentSubCatVar.QuestionCategoryParentId,
		QuestionCategoryParent:   "",
		QuestionCategoryId:       questionCategory,
		QuestionCategory:         "",
		Ques:                     quesMVars,
	}

	b, err := json.Marshal(ParentSubCatVar)
	commons.CheckErr(err)

	// mc := memcache.New("127.0.0.1:11211")

	// mc.Set(&memcache.Item{Key: "key_one", Value: []byte("michael")})
	// mc.Set(&memcache.Item{Key: "key_two", Value: []byte("programming")})

	// // Get a single value
	// val, err := mc.Get("key_one")
	// fmt.Println("////", val)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	//	fmt.Println(string(b))
	return string(b)
}

//getting hotel response at initial stage

func HotelResponse(hres *commons.HotelRes) string {

	db = GetDB()
	fmt.Println("I am on the way to store hotel initial answers!!")
	//insrtstmt := `INSERT INTO users (age, email, first_name, last_name) VALUES ($1, $2, $3, $4)`

	for i := range hres.Ans {
		for j := range hres.Ans[i].Answer {

			fmt.Println(".......", hres.Ans[i].QuestionId, hres.Ans[i].Answer[j].AnswerId, hres.Ans[i].Answer[j].Answer, hres.Ans[i].GroupQuestionId, hres.Ans[i].QuestionSubTypeId, hres.ClientTypeMasterId, hres.TravelAgencyMasterId)
			insrtstmt, err := db.Prepare(`INSERT INTO clientanswer SET questionMasterId = ?, answerId = ?, answer = ?, groupQuestionId = ?, questionSubTypeId = ?, clientTypeMasterId = ?, travelAgencyMasterId= ?`)
			fmt.Println("/////")
			commons.CheckErr(err)
			fmt.Println("\\\\\\")
			res, err := insrtstmt.Exec(hres.Ans[i].QuestionId, hres.Ans[i].Answer[j].AnswerId, hres.Ans[i].Answer[j].Answer, hres.Ans[i].GroupQuestionId, hres.Ans[i].QuestionSubTypeId, hres.ClientTypeMasterId, hres.TravelAgencyMasterId)
			fmt.Println(".....,,,,")
			commons.CheckErr(err)
			fmt.Println(res.LastInsertId)
		}
	}

	return hres.Ans[0].Answer[0].Answer

}

func GetDB() *sql.DB {
	fmt.Println("I am inside db")
	var err error
	if db == nil {
		db, err = sql.Open("mysql", "root:@/company_policy?parseTime=true&charset=utf8")
		//db, err = sql.Open("mysql", "sriram:sriram123@tcp(127.0.0.1:3306)/hotnix_dev?charset=utf8")

		commons.CheckErr(err)
	}

	return db
}
