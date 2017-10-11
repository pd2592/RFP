package db

import (
	"database/sql"
	"encoding/json"
	"fmt"

	commons "../commons"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func QuestionJsonByCat(parentId string) string {
	db = GetDB()
	var parentCatVar commons.ParentCat
	var questionCatVar commons.QuestionCat
	var questionCatVars []commons.QuestionCat
	var quesMVar commons.QuesM
	//var quesMVars []commons.QuesM
	var ansMVar commons.AnsM
	//var ansMVars []commons.AnsM

	stmt, err := db.Query("Select qm.questionCategoryParentId, qm.questionCategoryId, qq.questionCategory from questionmaster as qm JOIN questioncategory as qq ON qm.questionCategoryId = qq.questionCategoryId where qq.parentId = '" + parentId + "' GROUP BY qq.questionCategoryId")
	commons.CheckErr(err)
	for stmt.Next() {
		err := stmt.Scan(&parentCatVar.QuestionCategoryParentId, &questionCatVar.QuestionCategoryId, &questionCatVar.QuestionCategory)
		commons.CheckErr(err)

		var quesMVars []commons.QuesM

		stmt1, err := db.Query("Select questionMasterId, questionText, groupQuestionId, questionSubTypeId, isMandatory from questionmaster where questionCategoryId = '" + questionCatVar.QuestionCategoryId + "'")
		for stmt1.Next() {
			err := stmt1.Scan(&quesMVar.QuestionId, &quesMVar.QuestionText, &quesMVar.GroupQuestionId, &quesMVar.QuestionSubTypeId, &quesMVar.IsMandatory)
			commons.CheckErr(err)

			var ansMVars []commons.AnsM

			stmt2, err := db.Query("Select answerMasterId, answerText from answermaster where questionMasterId = '" + quesMVar.QuestionId + "'")
			for stmt2.Next() {
				err := stmt2.Scan(&ansMVar.AnswerId, &ansMVar.AnswerText)
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
