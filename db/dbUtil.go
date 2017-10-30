package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	commons "../commons"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

//=====================================================================================
// **********************hotel(stage - 1)****************************
//=====================================================================================

//get all parent category
func ParentCategory() string {
	db = GetDB()
	var parentCategoryVar commons.ParentCategory
	var parentCategoryVars []commons.ParentCategory
	retr_stmt, err := db.Query("Select questionCategoryId, questionCategory from questioncategory where parentId = 0")
	commons.CheckErr(err)
	for retr_stmt.Next() {
		err := retr_stmt.Scan(&parentCategoryVar.QuestionCategoryId, &parentCategoryVar.QuestionCategory)
		commons.CheckErr(err)
		parentCategoryVar = commons.ParentCategory{
			QuestionCategoryId: parentCategoryVar.QuestionCategoryId,
			QuestionCategory:   parentCategoryVar.QuestionCategory,
		}
		parentCategoryVars = append(parentCategoryVars, parentCategoryVar)
	}
	b, err := json.Marshal(parentCategoryVars)
	commons.CheckErr(err)
	//	fmt.Println(string(b))
	return string(b)
}

//getting question from db category wise

func QuestionJsonByCat(parentId string) string {
	db = GetDB()
	var parentCatVar commons.ParentCat
	var questionCatVar commons.QuestionCat
	var questionCatVars []commons.QuestionCat
	var groupQuestionVar commons.GroupQuestion

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
				QuestionId:        quesMVar.QuestionId,
				QuestionText:      quesMVar.QuestionText,
				QuestionSubTypeId: quesMVar.QuestionSubTypeId,
				GroupQuestionId:   quesMVar.GroupQuestionId,
				IsMandatory:       quesMVar.IsMandatory,
				TabColumn:         groupQuestionVars,
				ConcatAns:         ansMVars,
				Answer:            nil,
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

		TravelAgencyMasterId:     "",
		QuestionCategoryParentId: parentId,
		QuestionCategoryParent:   "test",
		QuesCategory:             questionCatVars,
	}
	b, err := json.Marshal(parentCatVar)
	commons.CheckErr(err)
	//	fmt.Println(string(b))
	return string(b)
}

//editing and previewing own ans setting  ------edit required
func HotelEditResponse(parentId, travelAgencyMasterId string) string {
	db = GetDB()
	var parentCatVar commons.ParentCat
	var questionCatVar commons.QuestionCat
	var questionCatVars []commons.QuestionCat
	var groupQuestionVar commons.GroupQuestion

	var quesMVar commons.QuesM
	//var quesMVars []commons.QuesM
	var ansMVar commons.AnsM
	var answersVar commons.Answers
	//var ansMVars []commons.AnsM

	var count string
	var method string
	err := db.QueryRow("Select COUNT(*) from clientanswer where travelAgencyMasterId = '" + travelAgencyMasterId + "'").Scan(&count)
	commons.CheckErr(err)
	if count == "0" {
		method = "create"
	} else {
		method = "edit"
	}

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
			var answersVars []commons.Answers
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
			if quesMVar.QuestionSubTypeId != "10" {

				retr_stmt4, err := db.Query("Select answerId, answer, questionSubTypeId, groupQuestionId from clientanswer where questionMasterId = '" + quesMVar.QuestionId + "' and travelAgencyMasterId = '" + travelAgencyMasterId + "' and groupQuestionId = '0'")
				commons.CheckErr(err)
				for retr_stmt4.Next() {
					err := retr_stmt4.Scan(&answersVar.AnswerId, &answersVar.Answer, &answersVar.QuestionSubTypeId, &answersVar.GroupQuestionMasterId)
					commons.CheckErr(err)
					answersVar = commons.Answers{
						AnswerId:              answersVar.AnswerId,
						Answer:                answersVar.Answer,
						Priority:              "",
						QuestionSubTypeId:     answersVar.QuestionSubTypeId,
						GroupQuestionMasterId: answersVar.GroupQuestionMasterId,
					}
					answersVars = append(answersVars, answersVar)
				}

			} else {

				fmt.Println("/////")
				retr_stmt5, err := db.Query("Select groupQuestionMasterId from groupquestion where groupQuestionId = '" + quesMVar.GroupQuestionId + "'")
				commons.CheckErr(err)
				for retr_stmt5.Next() {
					err := retr_stmt5.Scan(&quesMVar.QuestionId)
					commons.CheckErr(err)

					retr_stmt6, err := db.Query("Select answerId, answer, questionSubTypeId, groupQuestionId from clientanswer where questionMasterId = '" + quesMVar.QuestionId + "' and travelAgencyMasterId = '" + travelAgencyMasterId + "' and groupQuestionId = '" + quesMVar.GroupQuestionId + "'")
					commons.CheckErr(err)
					for retr_stmt6.Next() {
						err := retr_stmt6.Scan(&answersVar.AnswerId, &answersVar.Answer, &answersVar.QuestionSubTypeId, &answersVar.GroupQuestionMasterId)
						commons.CheckErr(err)
						answersVar = commons.Answers{
							AnswerId:              answersVar.AnswerId,
							Answer:                answersVar.Answer,
							Priority:              "",
							QuestionSubTypeId:     answersVar.QuestionSubTypeId,
							GroupQuestionMasterId: answersVar.GroupQuestionMasterId,
						}
						answersVars = append(answersVars, answersVar)
					}
				}
			}
			quesMVar = commons.QuesM{
				QuestionId:        quesMVar.QuestionId,
				QuestionText:      quesMVar.QuestionText,
				QuestionSubTypeId: quesMVar.QuestionSubTypeId,
				GroupQuestionId:   quesMVar.GroupQuestionId,
				IsMandatory:       quesMVar.IsMandatory,
				TabColumn:         groupQuestionVars,
				ConcatAns:         ansMVars,
				Answer:            answersVars,
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

		Method:                   method,
		TravelAgencyMasterId:     "",
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

	err := db.QueryRow("Select parentId from questioncategory where questionCategoryId = '" + questionCategory + "'").Scan(&ParentSubCatVar.QuestionCategoryParentId)
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
	fmt.Println("Storing hotel initial answers!!")
	//insrtstmt := `INSERT INTO users (age, email, first_name, last_name) VALUES ($1, $2, $3, $4)`

	for i := range hres.Ans {
		for j := range hres.Ans[i].Answer {

			if CheckHotelAnswers(hres.Ans[i].QuestionId, hres.TravelAgencyMasterId, hres.Ans[i].Answer[j].GroupQuestionMasterId) {
				// if hres.Ans[i].Answer[j].GroupQuestionMasterId != "0" {
				// 	hres.Ans[i].QuestionId = hres.Ans[i].Answer[j].GroupQuestionMasterId
				// }
				del, err := db.Exec("delete from clientanswer where questionMasterId = '" + hres.Ans[i].QuestionId + "' and travelAgencyMasterId = '" + hres.TravelAgencyMasterId + "' and groupQuestionId = '" + hres.Ans[i].Answer[j].GroupQuestionMasterId + "'")
				commons.CheckErr(err)
				fmt.Println(del, " Row deleted")
			}

			if hres.Ans[i].QuestionSubTypeId == "10" {
				insrtstmt, err := db.Prepare(`INSERT INTO clientanswer SET questionMasterId = ?, answerId = ?, answer = ?, groupQuestionId = ?, questionSubTypeId = ?, clientTypeMasterId = ?, travelAgencyMasterId= ?`)
				fmt.Println("/////")
				commons.CheckErr(err)
				res, err := insrtstmt.Exec(hres.Ans[i].Answer[j].GroupQuestionMasterId, hres.Ans[i].Answer[j].AnswerId, hres.Ans[i].Answer[j].Answer, hres.Ans[i].GroupQuestionId, hres.Ans[i].Answer[j].QuestionSubTypeId, hres.ClientTypeMasterId, hres.TravelAgencyMasterId)
				fmt.Println(".....,,,,")
				commons.CheckErr(err)
				fmt.Println(res.LastInsertId)
			} else {

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
	}

	return hres.Ans[0].Answer[0].Answer

}

//=====================================================================================
// **********************company(stage - 2)****************************
//=====================================================================================

func GetBasicList() string {
	db = GetDB()

	var QuestionVar commons.BQuestion
	var QuestionsVar []commons.BQuestion
	var BasicQuestionVar commons.BasicQuestion

	retr_stmt, err := db.Query("Select basicQuestionId, basicQuestion, bSubTypeId, division from basicquestion")
	commons.CheckErr(err)

	for retr_stmt.Next() {
		err := retr_stmt.Scan(&QuestionVar.BqId, &QuestionVar.BqText, &QuestionVar.BSubType, &QuestionVar.Divison)
		commons.CheckErr(err)
		QuestionVar = commons.BQuestion{
			BSubType: QuestionVar.BSubType,
			BqId:     QuestionVar.BqId,
			BqText:   QuestionVar.BqText,
			Divison:  QuestionVar.Divison,
			Answer:   " ",
			AnswerId: " ",
		}
		QuestionsVar = append(QuestionsVar, QuestionVar)
	}
	BasicQuestionVar = commons.BasicQuestion{
		RfpId:                " ",
		RfpName:              " ",
		TravelAgencyMasterId: " ",
		Ques:                 QuestionsVar,
	}
	b, err := json.Marshal(BasicQuestionVar)
	commons.CheckErr(err)
	return string(b)
}

func RfpBasicAns(RfpBasicQ *commons.BasicQuestion) string {
	db = GetDB()
	var result string
	insrtstmt, err := db.Prepare(`INSERT INTO rfpmaster SET rfpName = ?, travelAgencyMasterId = ?, clientTypeMasterId = ?, active = ?, completionStatus = ?, createDate = ?`)
	commons.CheckErr(err)
	createDate := time.Now()
	fmt.Println("rfp.TravelAgencyMasterId", RfpBasicQ.TravelAgencyMasterId)
	res, err := insrtstmt.Exec(RfpBasicQ.RfpName, RfpBasicQ.TravelAgencyMasterId, "4", "1", "0", createDate.String())
	commons.CheckErr(err)
	rfpid, err := res.LastInsertId()
	count := 0
	length := len(RfpBasicQ.Ques)
	for i := range RfpBasicQ.Ques {
		insrtstmt1, err := db.Prepare(`INSERT INTO basicrfpinfo SET rfpId = ?, travelAgencyMasterId = ?, basicQuestionId = ?, answer = ?, basicAnswerId = ?`)
		commons.CheckErr(err)
		res1, err := insrtstmt1.Exec(rfpid, RfpBasicQ.TravelAgencyMasterId, RfpBasicQ.Ques[i].BqId, RfpBasicQ.Ques[i].Answer, RfpBasicQ.Ques[i].AnswerId)
		commons.CheckErr(err)
		count++
		fmt.Println(res1)
	}
	if count == length {
		result = "all questions are saved"
	} else {
		result = "Failed to save some responses"
	}
	return result
}

func RfpRequest(rfp *commons.Rfp) string {

	db = GetDB()
	fmt.Println("Creating rfp!!")

	fmt.Println(rfp.CustomiseQues[0].QuestionText)
	var QuestionCategoryId string

	insrtstmt, err := db.Prepare(`INSERT INTO rfpmaster SET rfpName = ?, travelAgencyMasterId = ?, clientTypeMasterId = ?, active = ?, completionStatus = ?, createDate = ?`)
	commons.CheckErr(err)
	createDate := time.Now()
	fmt.Println("rfp.TravelAgencyMasterId", rfp.TravelAgencyMasterId)
	res, err := insrtstmt.Exec(rfp.RfpName, rfp.TravelAgencyMasterId, "4", "1", rfp.Status, createDate.String())
	commons.CheckErr(err)
	rowid, err := res.LastInsertId()
	commons.CheckErr(err)
	fmt.Println("rfp master updated, id : ", rowid)

	for i := range rfp.Ques {
		insrtstmt1, err := db.Prepare(`INSERT INTO rfpquestion SET rfpId = ?, questionMasterId = ?, groupQuestionId = ?, isMandatory = ?, travelAgencyMasterId = ?`)
		commons.CheckErr(err)
		res1, err := insrtstmt1.Exec(rowid, rfp.Ques[i].QuestionMasterId, rfp.Ques[i].GroupQuestionId, rfp.Ques[i].IsMandatory, rfp.TravelAgencyMasterId)
		commons.CheckErr(err)
		rowid1, err := res1.LastInsertId()
		commons.CheckErr(err)
		fmt.Println("rfp questions updated, id : ", rowid1)

		for j := range rfp.Ques[i].Interests {
			//fmt.Println(".......", hres.Ans[i].QuestionId, hres.Ans[i].Answer[j].AnswerId, hres.Ans[i].Answer[j].Answer, hres.Ans[i].GroupQuestionId, hres.Ans[i].QuestionSubTypeId, hres.ClientTypeMasterId, hres.TravelAgencyMasterId)

			fmt.Println(rfp.Ques[i].Interests[j])

			insrtstmt2, err := db.Prepare(`INSERT INTO rfpquestionchoices SET rfpId = ?, rfpQuestionId = ?, answerMasterId = ?, groupQuestionId = ?`)
			commons.CheckErr(err)
			res2, err := insrtstmt2.Exec(rowid, rowid1, rfp.Ques[i].Interests[j].AnswerId, rfp.Ques[i].GroupQuestionId)
			commons.CheckErr(err)
			rowid2, err := res2.LastInsertId()
			commons.CheckErr(err)
			fmt.Println("rfp questions updated", rowid2)

			insrtstmt3, err := db.Prepare(`INSERT INTO rfpchoicepriority SET rfpQuestionChoiceId = ?, rfpId = ?, priorityNumber = ?`)
			commons.CheckErr(err)
			res3, err := insrtstmt3.Exec(rowid2, rowid, rfp.Ques[i].Interests[j].Priority)
			commons.CheckErr(err)
			rowid3, err := res3.LastInsertId()
			commons.CheckErr(err)
			fmt.Println("rfp questions updated", rowid3)

		}
	}

	for i := range rfp.CustomiseQues {
		err := db.QueryRow("Select parentId from questioncategory where questionCategoryId = '" + rfp.CustomiseQues[i].QuestionCategoryId + "'").Scan(&QuestionCategoryId)
		commons.CheckErr(err)
		insrtstmt4, err := db.Prepare(`INSERT INTO questionmaster SET questionCategoryParentId = ?, questionCategoryId = ?, questionPatternMasterId = ?, questionText = ?, questionSubTypeId = ?, clientTypeMasterId = ?, travelAgencyMasterId = ?, isMandatory = ?, status = ?, approved = ?, isCostRel = ?`)
		commons.CheckErr(err)
		res, err := insrtstmt4.Exec(QuestionCategoryId, rfp.CustomiseQues[i].QuestionCategoryId, "3", rfp.CustomiseQues[i].QuestionText, "1", "0", rfp.TravelAgencyMasterId, "1", "1", "0", "0")
		commons.CheckErr(err)
		questionMasterIdINST, err := res.LastInsertId()
		commons.CheckErr(err)

		insrtstmt5, err := db.Prepare(`INSERT INTO rfpquestion SET rfpId = ?, questionMasterId = ?, groupQuestionId = ?, isMandatory = ?, travelAgencyMasterId = ?`)
		commons.CheckErr(err)
		res5, err := insrtstmt5.Exec(rowid, questionMasterIdINST, "0", "1", rfp.TravelAgencyMasterId)
		commons.CheckErr(err)
		rowid5, err := res5.LastInsertId()
		commons.CheckErr(err)
		fmt.Println("rfp questions updated with customised question", rowid5)

	}
	return ""
}

func RfpEditor(rfp *commons.Rfp) string {
	db = GetDB()
	rfpId := rfp.RfpId
	var hotels string

	updatestmt, err := db.Prepare(`UPDATE rfpmaster SET completionStatus = ? WHERE rfpId = ?`)
	commons.CheckErr(err)

	res, err := updatestmt.Exec(rfp.Status, rfpId)
	commons.CheckErr(err)
	fmt.Println(res, "> status updated")

	for i := range rfp.Ques {
		if CheckDuplicate("rfpquestion", "questionMasterId", rfpId, rfp.Ques[i].QuestionMasterId) {

			fmt.Println(rfp.Ques[i].QuestionMasterId, " Question already there in rfpid ", rfpId)
		} else {
			insrtstmt1, err := db.Prepare(`INSERT INTO rfpquestion SET rfpId = ?, questionMasterId = ?, groupQuestionId = ?, isMandatory = ?, travelAgencyMasterId = ?`)
			commons.CheckErr(err)
			res1, err := insrtstmt1.Exec(rfpId, rfp.Ques[i].QuestionMasterId, rfp.Ques[i].GroupQuestionId, rfp.Ques[i].IsMandatory, rfp.TravelAgencyMasterId)
			commons.CheckErr(err)
			rowid1, err := res1.LastInsertId()
			commons.CheckErr(err)
			fmt.Println("rfp questions updated, id : ", rowid1)

			for j := range rfp.Ques[i].Interests {
				//fmt.Println(".......", hres.Ans[i].QuestionId, hres.Ans[i].Answer[j].AnswerId, hres.Ans[i].Answer[j].Answer, hres.Ans[i].GroupQuestionId, hres.Ans[i].QuestionSubTypeId, hres.ClientTypeMasterId, hres.TravelAgencyMasterId)

				fmt.Println(rfp.Ques[i].Interests[j])

				insrtstmt2, err := db.Prepare(`INSERT INTO rfpquestionchoices SET rfpId = ?, rfpQuestionId = ?, answerMasterId = ?, groupQuestionId = ?`)
				commons.CheckErr(err)
				res2, err := insrtstmt2.Exec(rfpId, rowid1, rfp.Ques[i].Interests[j].AnswerId, rfp.Ques[i].GroupQuestionId)
				commons.CheckErr(err)
				rowid2, err := res2.LastInsertId()
				commons.CheckErr(err)
				fmt.Println("rfp questions updated", rowid2)

				insrtstmt3, err := db.Prepare(`INSERT INTO rfpchoicepriority SET rfpQuestionChoiceId = ?, rfpId = ?, priorityNumber = ?`)
				commons.CheckErr(err)
				res3, err := insrtstmt3.Exec(rowid2, rfpId, rfp.Ques[i].Interests[j].Priority)
				commons.CheckErr(err)
				rowid3, err := res3.LastInsertId()
				commons.CheckErr(err)
				fmt.Println("rfp questions updated", rowid3)

			}
		}
	}
	//customised question editing
	var QuestionCategoryId string

	for j := range rfp.CustomiseQues {

		err := db.QueryRow("Select parentId from questioncategory where questionCategoryId = '" + rfp.CustomiseQues[j].QuestionCategoryId + "'").Scan(&QuestionCategoryId)
		commons.CheckErr(err)
		insrtstmt4, err := db.Prepare(`INSERT INTO questionmaster SET questionCategoryParentId = ?, questionCategoryId = ?, questionPatternMasterId = ?, questionText = ?, questionSubTypeId = ?, clientTypeMasterId = ?, travelAgencyMasterId = ?, isMandatory = ?, status = ?, approved = ?, isCostRel = ?`)
		commons.CheckErr(err)
		res, err := insrtstmt4.Exec(QuestionCategoryId, rfp.CustomiseQues[j].QuestionCategoryId, "3", rfp.CustomiseQues[j].QuestionText, "1", "0", rfp.TravelAgencyMasterId, "1", "1", "0", "0")
		commons.CheckErr(err)
		questionMasterIdINST, err := res.LastInsertId()
		commons.CheckErr(err)

		insrtstmt5, err := db.Prepare(`INSERT INTO rfpquestion SET rfpId = ?, questionMasterId = ?, groupQuestionId = ?, isMandatory = ?, travelAgencyMasterId = ?`)
		commons.CheckErr(err)
		res5, err := insrtstmt5.Exec(rfpId, questionMasterIdINST, "0", "1", rfp.TravelAgencyMasterId)
		commons.CheckErr(err)
		rowid5, err := res5.LastInsertId()
		commons.CheckErr(err)
		fmt.Println("rfp questions updated with customised question", rowid5)

		if rfp.CustomiseQues[j].CustomiseQuestionId != "" {
			del, err := db.Exec("delete from rfpquestion where rfpId = '" + rfpId + "' and questionMasterId = '" + rfp.CustomiseQues[j].CustomiseQuestionId + "'")
			commons.CheckErr(err)

			affect, err := del.RowsAffected()
			commons.CheckErr(err)
			fmt.Println("question deleted from rfp questions : ", affect)
		}
	}

	if rfp.Status == "1" {
		fmt.Println("List Hotels as comapny want to send the rfp")
		hotels = ListHotels("1")
	}

	return hotels
}

func ListHotels(CityId string) string {
	db = GetDB()
	fmt.Println("hotels...")
	var ListHotelVar commons.ListHotel
	var ListHotelsVar []commons.ListHotel

	retr_stmt, err := db.Query("Select hotelsMasterId, hotelName, rstarRating, cityLocalityId, cityMasterId, distanceFromCity from hotelsmaster where cityMasterId = '" + CityId + "'")
	commons.CheckErr(err)

	for retr_stmt.Next() {
		//var quesVar commons.Ques
		err = retr_stmt.Scan(&ListHotelVar.HotelId, &ListHotelVar.HotelName, &ListHotelVar.Star, &ListHotelVar.Locality, &ListHotelVar.City, &ListHotelVar.DistanceFromCity)
		commons.CheckErr(err)
		ListHotelVar = commons.ListHotel{
			HotelName:        ListHotelVar.HotelName,
			HotelId:          ListHotelVar.HotelId,
			Star:             ListHotelVar.Star,
			Locality:         ListHotelVar.Locality,
			City:             ListHotelVar.City,
			DistanceFromCity: ListHotelVar.DistanceFromCity,
		}
		ListHotelsVar = append(ListHotelsVar, ListHotelVar)
	}
	b, err := json.Marshal(ListHotelsVar)
	commons.CheckErr(err)
	return string(b)
}

func RfpSend(rfp *commons.RfpSend) string {
	db = GetDB()

	for i := range rfp.Hotels {
		insrtstmt, err := db.Prepare(`INSERT INTO rfphotelmapping SET rfpId = ?, travelAgencyMasterId = ?, hotelMasterId = ?, quotedPrice = ?, accepted = ?, status = ?`)
		fmt.Println("/////")
		commons.CheckErr(err)
		res, err := insrtstmt.Exec(rfp.RfpId, rfp.TravelAgencyMasterId, rfp.Hotels[i], "0", "0", "1")
		fmt.Println(".....,,,,")
		commons.CheckErr(err)
		fmt.Println(res.LastInsertId)
	}
	return ""
}

func GetRfp(RfpId string) string {
	db = GetDB()
	//var RfpViewVar commons.RfpView
	var quesVar commons.Ques
	var questions []commons.Ques
	var RfpQuestionId string
	retr_stmt, err := db.Query("Select rfpQuestionId, questionMasterId, groupQuestionId from rfpquestion where rfpId = '" + RfpId + "'")
	commons.CheckErr(err)

	for retr_stmt.Next() {
		//var quesVar commons.Ques
		err = retr_stmt.Scan(&RfpQuestionId, &quesVar.QuestionId, &quesVar.GroupQuestionId)
		commons.CheckErr(err)

		err := db.QueryRow("Select questionText from questionmaster where questionMasterId = '" + quesVar.QuestionId + "'").Scan(&quesVar.QuestionText)
		commons.CheckErr(err)

		retr_stmt1, err := db.Query("Select answerMasterId from rfpquestionchoices where rfpQuestionId = '" + RfpQuestionId + "'")
		commons.CheckErr(err)

		var Interests []commons.Interest
		for retr_stmt1.Next() {
			var InterestVar commons.Interest

			err = retr_stmt1.Scan(&InterestVar.AnswerId)
			commons.CheckErr(err)

			err := db.QueryRow("Select answerText from answermaster where answerMasterId = '" + InterestVar.AnswerId + "'").Scan(&InterestVar.Answer)
			commons.CheckErr(err)

			InterestVar = commons.Interest{
				AnswerId: InterestVar.AnswerId,
				Answer:   InterestVar.Answer,
			}
			Interests = append(Interests, InterestVar)
		}

		quesVar = commons.Ques{
			QuestionId:      quesVar.QuestionId,
			QuestionText:    quesVar.QuestionText,
			GroupQuestionId: quesVar.GroupQuestionId,
			Interests:       Interests,
		}
		questions = append(questions, quesVar)

	}

	RfpViewVar := commons.RfpView{
		RfpID: RfpId,
		Ques:  questions,
	}
	b, err := json.Marshal(RfpViewVar)
	commons.CheckErr(err)
	return string(b)
}

func CheckDuplicate(tablename, columnname, rfpId, qId string) bool {
	db = GetDB()
	var columnvalue string
	err := db.QueryRow("Select COUNT(" + columnname + ") from " + tablename + " where rfpId = '" + rfpId + "' and questionMasterId = '" + qId + "'").Scan(&columnvalue)
	commons.CheckErr(err)
	if columnvalue != "0" {
		return true
	} else {
		return false
	}
}

func CheckHotelAnswers(questionMasterId, travelAgencyMasterId, groupQuestionId string) bool {
	fmt.Println("....")
	db = GetDB()
	var columnvalue string
	err := db.QueryRow("Select COUNT(*) from clientanswer where questionMasterId = '" + questionMasterId + "' and travelAgencyMasterId = '" + travelAgencyMasterId + "' and groupQuestionId = '" + groupQuestionId + "'").Scan(&columnvalue)
	commons.CheckErr(err)
	if columnvalue != "0" {
		return true
	} else {
		return false
	}
}

func GetDB() *sql.DB {
	//fmt.Println("I am inside db")
	var err error
	if db == nil {
		//db, err = sql.Open("mysql", "root:@/company_policy?parseTime=true&charset=utf8")
		db, err = sql.Open("mysql", "sriram:sriram123@tcp(127.0.0.1:3306)/hotnix_dev?parseTime=true&charset=utf8")

		commons.CheckErr(err)
	}

	return db
}
