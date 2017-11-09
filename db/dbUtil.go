package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
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

	var BasicQuestionVar commons.BasicQuestion
	var BDivisionVar commons.BDivision
	var BDivisionsVar []commons.BDivision
	retr_stmt, err := db.Query("Select division, divisionName from basicquestion GROUP BY division")
	commons.CheckErr(err)

	for retr_stmt.Next() {
		var division commons.LabVal
		var QuestionsVar []commons.BQuestion

		err := retr_stmt.Scan(&BDivisionVar.Division.Value, &BDivisionVar.Division.Label)
		commons.CheckErr(err)
		division = commons.LabVal{
			Label: BDivisionVar.Division.Label,
			Value: BDivisionVar.Division.Value,
		}
		retr_stmt1, err := db.Query("Select basicQuestionId, basicQuestion, bSubTypeId from basicquestion where division = '" + BDivisionVar.Division.Value + "'")
		commons.CheckErr(err)
		for retr_stmt1.Next() {
			var QuestionVar commons.BQuestion

			err := retr_stmt1.Scan(&QuestionVar.BqId, &QuestionVar.BqText, &QuestionVar.BSubType)
			commons.CheckErr(err)

			QuestionVar = commons.BQuestion{
				BSubType: QuestionVar.BSubType,
				BqId:     QuestionVar.BqId,
				BqText:   QuestionVar.BqText,
				Answer:   " ",
				AnswerId: nil,
			}
			QuestionsVar = append(QuestionsVar, QuestionVar)
		}
		BDivisionVar = commons.BDivision{
			Division: division,
			Ques:     QuestionsVar,
		}

		BDivisionsVar = append(BDivisionsVar, BDivisionVar)

	}
	BasicQuestionVar = commons.BasicQuestion{
		RfpId:                " ",
		RfpName:              " ",
		TravelAgencyMasterId: " ",
		Division:             BDivisionsVar,
	}
	b, err := json.Marshal(BasicQuestionVar)
	commons.CheckErr(err)
	return string(b)
}

func RfpBasicAns(RfpBasicQ *commons.BasicQuestion) string {
	db = GetDB()
	var rfpid int64

	if RfpBasicQ.RfpId == "" {
		insrtstmt, err := db.Prepare(`INSERT INTO rfpmaster SET rfpName = ?, travelAgencyMasterId = ?, clientTypeMasterId = ?, active = ?, completionStatus = ?, createDate = ?`)
		commons.CheckErr(err)
		createDate := time.Now()
		fmt.Println("rfp.TravelAgencyMasterId", RfpBasicQ.TravelAgencyMasterId)
		res, err := insrtstmt.Exec(RfpBasicQ.RfpName, RfpBasicQ.TravelAgencyMasterId, "4", "1", "0", createDate.String())
		commons.CheckErr(err)
		rfpid, err = res.LastInsertId()
		for i := range RfpBasicQ.Division {
			for j := range RfpBasicQ.Division[i].Ques {
				if RfpBasicQ.Division[i].Ques[j].BqId == "5" {
					for k := range RfpBasicQ.Division[i].Ques[j].AnswerId {
						insrtstmt1, err := db.Prepare(`INSERT INTO basicrfpinfo SET rfpId = ?, travelAgencyMasterId = ?, basicQuestionId = ?, answer = ?, basicAnswerId = ?`)
						commons.CheckErr(err)
						res1, err := insrtstmt1.Exec(rfpid, RfpBasicQ.TravelAgencyMasterId, RfpBasicQ.Division[i].Ques[j].BqId, RfpBasicQ.Division[i].Ques[j].AnswerId[k].Label, RfpBasicQ.Division[i].Ques[j].AnswerId[k].Value)
						commons.CheckErr(err)
						fmt.Println(res1.LastInsertId)

					}
				} else {
					insrtstmt1, err := db.Prepare(`INSERT INTO basicrfpinfo SET rfpId = ?, travelAgencyMasterId = ?, basicQuestionId = ?, answer = ?`)
					commons.CheckErr(err)
					res1, err := insrtstmt1.Exec(rfpid, RfpBasicQ.TravelAgencyMasterId, RfpBasicQ.Division[i].Ques[j].BqId, RfpBasicQ.Division[i].Ques[j].Answer)
					commons.CheckErr(err)
					fmt.Println(res1.LastInsertId)

				}
				fmt.Println("rfpid", string(rfpid))

			}
		}
		fmt.Println("rfpid", rfpid)
		var RfpDataVar commons.RfpData
		RfpDataVar = commons.RfpData{
			RfpId:   strconv.Itoa(int(rfpid)),
			RfpName: RfpBasicQ.RfpName,
		}
		b, err := json.Marshal(RfpDataVar)
		commons.CheckErr(err)
		return string(b)
	} else {
		updatestmt, err := db.Prepare(`UPDATE rfpmaster SET rfpName = ?, travelAgencyMasterId = ?, clientTypeMasterId = ?, active = ?, completionStatus = ? where rfpId = '` + RfpBasicQ.RfpId + "'")
		commons.CheckErr(err)
		res1, err := updatestmt.Exec(RfpBasicQ.RfpName, RfpBasicQ.TravelAgencyMasterId, "4", "1", "0")
		commons.CheckErr(err)
		fmt.Println(res1)

		count := 0
		//length := len(RfpBasicQ.Ques)
		for i := range RfpBasicQ.Division {
			for j := range RfpBasicQ.Division[i].Ques {
				if CheckDuplicateBasic(RfpBasicQ.RfpId, RfpBasicQ.Division[i].Ques[j].BqId) {
					if RfpBasicQ.Division[i].Ques[j].BqId == "5" {
						del, err := db.Exec("delete from basicrfpinfo where rfpId = '" + RfpBasicQ.RfpId + "' and basicQuestionId = '" + RfpBasicQ.Division[i].Ques[j].BqId + "'")
						commons.CheckErr(err)
						fmt.Println(del)

						for k := range RfpBasicQ.Division[i].Ques[j].AnswerId {
							insrtstmt1, err := db.Prepare(`INSERT INTO basicrfpinfo SET rfpId = ?, travelAgencyMasterId = ?, basicQuestionId = ?, answer = ?, basicAnswerId = ?`)
							commons.CheckErr(err)
							res1, err := insrtstmt1.Exec(RfpBasicQ.RfpId, RfpBasicQ.TravelAgencyMasterId, RfpBasicQ.Division[i].Ques[j].BqId, RfpBasicQ.Division[i].Ques[j].AnswerId[k].Label, RfpBasicQ.Division[i].Ques[j].AnswerId[k].Value)
							commons.CheckErr(err)
							fmt.Println(res1.LastInsertId)

						}
					} else {
						updatestmt, err := db.Prepare("UPDATE basicrfpinfo SET basicQuestionId = ?, answer = ?, basicAnswerId = ? where rfpId = '" + RfpBasicQ.RfpId + "' and basicQuestionId = '" + RfpBasicQ.Division[i].Ques[j].BqId + "'")
						commons.CheckErr(err)
						res1, err := updatestmt.Exec(RfpBasicQ.Division[i].Ques[j].BqId, RfpBasicQ.Division[i].Ques[j].Answer, RfpBasicQ.Division[i].Ques[j].AnswerId)
						commons.CheckErr(err)
						count++
						fmt.Println(res1, "updated")
					}

				} else {
					if RfpBasicQ.Division[i].Ques[j].BqId == "5" {
						for k := range RfpBasicQ.Division[i].Ques[j].AnswerId {
							insrtstmt1, err := db.Prepare(`INSERT INTO basicrfpinfo SET rfpId = ?, travelAgencyMasterId = ?, basicQuestionId = ?, answer = ?, basicAnswerId = ?`)
							commons.CheckErr(err)
							res1, err := insrtstmt1.Exec(RfpBasicQ.RfpId, RfpBasicQ.TravelAgencyMasterId, RfpBasicQ.Division[i].Ques[j].BqId, RfpBasicQ.Division[i].Ques[j].AnswerId[k].Label, RfpBasicQ.Division[i].Ques[j].AnswerId[k].Value)
							commons.CheckErr(err)
							fmt.Println(res1.LastInsertId)

						}
					} else {
						insrtstmt1, err := db.Prepare(`INSERT INTO basicrfpinfo SET rfpId = ?, travelAgencyMasterId = ?, basicQuestionId = ?, answer = ?, basicAnswerId = ?`)
						commons.CheckErr(err)
						res1, err := insrtstmt1.Exec(RfpBasicQ.RfpId, RfpBasicQ.TravelAgencyMasterId, RfpBasicQ.Division[i].Ques[j].BqId, RfpBasicQ.Division[i].Ques[j].Answer, RfpBasicQ.Division[i].Ques[j].AnswerId)
						commons.CheckErr(err)
						count++
						fmt.Println(res1)
					}
				}
			}
		}
		var RfpDataVar commons.RfpData
		RfpDataVar = commons.RfpData{
			RfpId:   RfpBasicQ.RfpId,
			RfpName: RfpBasicQ.RfpName,
		}
		b, err := json.Marshal(RfpDataVar)
		commons.CheckErr(err)

		return string(b)
	}

	//	result = "all questions are saved"

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
		res1, err := insrtstmt1.Exec(rowid, rfp.Ques[i].QuestionId, rfp.Ques[i].GroupQuestionId, rfp.Ques[i].IsMandatory, rfp.TravelAgencyMasterId)
		commons.CheckErr(err)
		rowid1, err := res1.LastInsertId()
		commons.CheckErr(err)
		fmt.Println("rfp questions updated, id : ", rowid1)

		for j := range rfp.Ques[i].Answer {
			//fmt.Println(".......", hres.Ans[i].QuestionId, hres.Ans[i].Answer[j].AnswerId, hres.Ans[i].Answer[j].Answer, hres.Ans[i].GroupQuestionId, hres.Ans[i].QuestionSubTypeId, hres.ClientTypeMasterId, hres.TravelAgencyMasterId)

			fmt.Println(rfp.Ques[i].Answer[j])

			insrtstmt2, err := db.Prepare(`INSERT INTO rfpquestionchoices SET rfpId = ?, rfpQuestionId = ?, answerMasterId = ?, groupQuestionId = ?`)
			commons.CheckErr(err)
			res2, err := insrtstmt2.Exec(rowid, rowid1, rfp.Ques[i].Answer[j].AnswerId, rfp.Ques[i].GroupQuestionId)
			commons.CheckErr(err)
			rowid2, err := res2.LastInsertId()
			commons.CheckErr(err)
			fmt.Println("rfp questions updated", rowid2)

			insrtstmt3, err := db.Prepare(`INSERT INTO rfpchoicepriority SET rfpQuestionChoiceId = ?, rfpId = ?, priorityNumber = ?`)
			commons.CheckErr(err)
			res3, err := insrtstmt3.Exec(rowid2, rowid, rfp.Ques[i].Answer[j].Priority)
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
	return "Rfp created"
}

func CompanyEditRfp(rfpId, parentId, travelAgencyMasterId string) string {
	db = GetDB()
	var parentCatVar commons.ParentCat
	var questionCatVar commons.QuestionCat
	var questionCatVars []commons.QuestionCat
	var groupQuestionVar commons.GroupQuestion

	var quesMVar commons.QuesM
	//var quesMVars []commons.QuesM
	var ansMVar commons.AnsM
	//	var answersVar commons.Answers
	//var ansMVars []commons.AnsM

	var count string
	var method string
	err := db.QueryRow("Select COUNT(*) from rfpquestion where rfpId = '" + rfpId + "'").Scan(&count)
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

		retr_stmt1, err := db.Query("Select questionMasterId, questionText, groupQuestionId, questionSubTypeId from questionmaster where questionCategoryId = '" + questionCatVar.QuestionCategoryId + "'")
		for retr_stmt1.Next() {
			err := retr_stmt1.Scan(&quesMVar.QuestionId, &quesMVar.QuestionText, &quesMVar.GroupQuestionId, &quesMVar.QuestionSubTypeId)
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
			//	if quesMVar.QuestionSubTypeId != "10" {
			var answersVar commons.Answers
			var rfpQuestionId string
			err = db.QueryRow("Select rfpQuestionId from rfpquestion where rfpId = '" + rfpId + "' and travelAgencyMasterId = '" + travelAgencyMasterId + "' and questionMasterId = '" + quesMVar.QuestionId + "' and groupQuestionId = '0'").Scan(&rfpQuestionId)
			commons.CheckErr(err)

			retr_stmt4, err := db.Query("Select answerMasterId from rfpquestionchoices where rfpQuestionId = '" + rfpQuestionId + "'")
			for retr_stmt4.Next() {
				err := retr_stmt4.Scan(&answersVar.AnswerId)
				commons.CheckErr(err)
				err = db.QueryRow("Select answerText from answermaster where answerMasterId = '" + answersVar.AnswerId + "'").Scan(&answersVar.Answer)

				answersVar = commons.Answers{
					AnswerId: answersVar.AnswerId,
					Answer:   answersVar.Answer,
				}
				answersVars = append(answersVars, answersVar)
			}

			// } else {

			// 	fmt.Println("/////")
			// 	retr_stmt5, err := db.Query("Select groupQuestionMasterId from groupquestion where groupQuestionId = '" + quesMVar.GroupQuestionId + "'")
			// 	commons.CheckErr(err)
			// 	for retr_stmt5.Next() {
			// 		err := retr_stmt5.Scan(&quesMVar.QuestionId)
			// 		commons.CheckErr(err)

			// 		retr_stmt6, err := db.Query("Select answerId, answer, questionSubTypeId, groupQuestionId from clientanswer where questionMasterId = '" + quesMVar.QuestionId + "' and travelAgencyMasterId = '" + travelAgencyMasterId + "' and groupQuestionId = '" + quesMVar.GroupQuestionId + "'")
			// 		commons.CheckErr(err)
			// 		for retr_stmt6.Next() {
			// 			err := retr_stmt6.Scan(&answersVar.AnswerId, &answersVar.Answer, &answersVar.QuestionSubTypeId, &answersVar.GroupQuestionMasterId)
			// 			commons.CheckErr(err)
			// 			answersVar = commons.Answers{
			// 				AnswerId:              answersVar.AnswerId,
			// 				Answer:                answersVar.Answer,
			// 				Priority:              "",
			// 				QuestionSubTypeId:     answersVar.QuestionSubTypeId,
			// 				GroupQuestionMasterId: answersVar.GroupQuestionMasterId,
			// 			}
			// 			answersVars = append(answersVars, answersVar)
			// 		}
			// 	}
			// }
			err = db.QueryRow("Select isMandatory from rfpquestion where rfpId = '" + rfpId + "' and travelAgencyMasterId = '" + travelAgencyMasterId + "' and questionMasterId = '" + quesMVar.QuestionId + "' and groupQuestionId = '0'").Scan(&quesMVar.IsMandatory)

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
		if CheckDuplicate("rfpquestion", "questionMasterId", rfpId, rfp.Ques[i].QuestionId) {

			fmt.Println(rfp.Ques[i].QuestionId, " Question already there in rfpid ", rfpId)
		} else {
			insrtstmt1, err := db.Prepare(`INSERT INTO rfpquestion SET rfpId = ?, questionMasterId = ?, groupQuestionId = ?, isMandatory = ?, travelAgencyMasterId = ?`)
			commons.CheckErr(err)
			res1, err := insrtstmt1.Exec(rfpId, rfp.Ques[i].QuestionId, rfp.Ques[i].GroupQuestionId, rfp.Ques[i].IsMandatory, rfp.TravelAgencyMasterId)
			commons.CheckErr(err)
			rowid1, err := res1.LastInsertId()
			commons.CheckErr(err)
			fmt.Println("rfp questions updated, id : ", rowid1)

			for j := range rfp.Ques[i].Answer {
				//fmt.Println(".......", hres.Ans[i].QuestionId, hres.Ans[i].Answer[j].AnswerId, hres.Ans[i].Answer[j].Answer, hres.Ans[i].GroupQuestionId, hres.Ans[i].QuestionSubTypeId, hres.ClientTypeMasterId, hres.TravelAgencyMasterId)

				fmt.Println(rfp.Ques[i].Answer[j])

				insrtstmt2, err := db.Prepare(`INSERT INTO rfpquestionchoices SET rfpId = ?, rfpQuestionId = ?, answerMasterId = ?, groupQuestionId = ?`)
				commons.CheckErr(err)
				res2, err := insrtstmt2.Exec(rfpId, rowid1, rfp.Ques[i].Answer[j].AnswerId, rfp.Ques[i].GroupQuestionId)
				commons.CheckErr(err)
				rowid2, err := res2.LastInsertId()
				commons.CheckErr(err)
				fmt.Println("rfp questions updated", rowid2)

				insrtstmt3, err := db.Prepare(`INSERT INTO rfpchoicepriority SET rfpQuestionChoiceId = ?, rfpId = ?, priorityNumber = ?`)
				commons.CheckErr(err)
				res3, err := insrtstmt3.Exec(rowid2, rfpId, rfp.Ques[i].Answer[j].Priority)
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
		cities := GetCity(rfp.RfpId)
		hotels = ListHotels(cities)
	} else {
		hotels = "Rfp Saved as draft"
	}

	return hotels
}

func GetCity(rfpId string) []string {

	db = GetDB()
	retr_stmt, err := db.Query("Select basicAnswerId from basicrfpinfo where rfpId = '" + rfpId + "' and basicQuestionId = 5")
	commons.CheckErr(err)
	var s []string
	var s1 string
	for retr_stmt.Next() {
		err = retr_stmt.Scan(&s1)
		s = append(s, s1)
	}
	fmt.Println(s)
	return s
}

func ListHotels(Cities []string) string {
	db = GetDB()
	fmt.Println("hotels...")
	var ListHotelsVar []commons.ListHotel
	for i := range Cities {
		retr_stmt, err := db.Query("Select hotelsMasterId, hotelName, rstarRating, cityLocalityId, cityMasterId, distanceFromCity from hotelsmaster where cityMasterId = '" + Cities[i] + "'")
		commons.CheckErr(err)
		var ListHotelVar commons.ListHotel

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
	}

	b, err := json.Marshal(ListHotelsVar)
	commons.CheckErr(err)
	return string(b)
}

func RfpSend(rfp *commons.RfpSend) string {
	db = GetDB()

	for i := range rfp.Hotels {
		insrtstmt, err := db.Prepare(`INSERT INTO rfphotelmapping SET rfpId = ?, travelAgencyMasterId = ?, hotelMasterId = ?, slabId = ?, accepted = ?, status = ?`)
		fmt.Println("/////")
		commons.CheckErr(err)
		res, err := insrtstmt.Exec(rfp.RfpId, rfp.TravelAgencyMasterId, rfp.Hotels[i], "0", "0", "1")
		fmt.Println(".....,,,,")
		updatestmt, err := db.Prepare(`UPDATE rfpmaster SET completionStatus = ? where rfpId = '` + rfp.RfpId + "'")
		commons.CheckErr(err)
		res1, err := updatestmt.Exec("1")
		commons.CheckErr(err)
		fmt.Println(res1.LastInsertId)
		commons.CheckErr(err)
		fmt.Println(res.LastInsertId)
	}
	return ""
}

func RfpFullView(RfpId string) string {
	db = GetDB()

	//var BasicQuestionVar commons.BasicQuestion
	var BDivisionVar commons.BDivision
	var BDivisionsVar []commons.BDivision
	retr_stmt, err := db.Query("Select division, divisionName from basicquestion GROUP BY division")
	commons.CheckErr(err)

	for retr_stmt.Next() {
		var division commons.LabVal
		var QuestionsVar []commons.BQuestion

		err := retr_stmt.Scan(&BDivisionVar.Division.Value, &BDivisionVar.Division.Label)
		commons.CheckErr(err)
		division = commons.LabVal{
			Label: BDivisionVar.Division.Label,
			Value: BDivisionVar.Division.Value,
		}
		var LabValVar commons.LabVal
		var LabValsVar []commons.LabVal
		retr_stmt1, err := db.Query("Select basicQuestionId, basicQuestion, bSubTypeId from basicquestion where division = '" + BDivisionVar.Division.Value + "'")
		commons.CheckErr(err)
		for retr_stmt1.Next() {
			var QuestionVar commons.BQuestion

			err := retr_stmt1.Scan(&QuestionVar.BqId, &QuestionVar.BqText, &QuestionVar.BSubType)
			commons.CheckErr(err)
			if QuestionVar.BqId != "5" {
				err = db.QueryRow("Select answer from basicrfpinfo where basicQuestionId = '" + QuestionVar.BqId + "' and rfpId = '" + RfpId + "'").Scan(&QuestionVar.Answer)
				QuestionVar = commons.BQuestion{
					BSubType: QuestionVar.BSubType,
					BqId:     QuestionVar.BqId,
					BqText:   QuestionVar.BqText,
					Answer:   QuestionVar.Answer,
					AnswerId: nil,
				}
			} else {
				retr_stmt2, err := db.Query("Select answer, basicAnswerId from basicrfpinfo where basicQuestionId = '" + QuestionVar.BqId + "' and rfpId = '" + RfpId + "'")
				commons.CheckErr(err)
				for retr_stmt2.Next() {
					err := retr_stmt2.Scan(&LabValVar.Label, &LabValVar.Value)
					commons.CheckErr(err)

					LabValVar = commons.LabVal{
						Label: LabValVar.Label,
						Value: LabValVar.Value,
					}
					LabValsVar = append(LabValsVar, LabValVar)
				}
				QuestionVar = commons.BQuestion{
					BSubType: QuestionVar.BSubType,
					BqId:     QuestionVar.BqId,
					BqText:   QuestionVar.BqText,
					Answer:   " ",
					AnswerId: LabValsVar,
				}
			}

			QuestionsVar = append(QuestionsVar, QuestionVar)
		}
		BDivisionVar = commons.BDivision{
			Division: division,
			Ques:     QuestionsVar,
		}

		BDivisionsVar = append(BDivisionsVar, BDivisionVar)

	}
	// populating questions by category

	//var count string
	method := ""

	var ParentCatsVar []commons.ParentCat
	retr_stmt7, err := db.Query("Select questionCategoryId, questionCategory from questioncategory where parentId = '0'")
	for retr_stmt7.Next() {

		var parentCatVar commons.ParentCat
		var questionCatVar commons.QuestionCat
		questionCatVars := make([]commons.QuestionCat, 0)
		var groupQuestionVar commons.GroupQuestion

		var quesMVar commons.QuesM
		//var quesMVars []commons.QuesM
		//var ansMVar commons.AnsM

		//var ansMVars []commons.AnsM
		err = retr_stmt7.Scan(&parentCatVar.QuestionCategoryParentId, &parentCatVar.QuestionCategoryParent)

		retr_stmt8, err := db.Query("Select qm.questionCategoryParentId, qm.questionCategoryId, qq.questionCategory from questionmaster as qm JOIN questioncategory as qq ON qm.questionCategoryId = qq.questionCategoryId where qq.parentId = '" + parentCatVar.QuestionCategoryParentId + "' GROUP BY qq.questionCategoryId")
		commons.CheckErr(err)

		for retr_stmt8.Next() {

			err := retr_stmt8.Scan(&parentCatVar.QuestionCategoryParentId, &questionCatVar.QuestionCategoryId, &questionCatVar.QuestionCategory)
			commons.CheckErr(err)

			quesMVars := make([]commons.QuesM, 0)

			retr_stmt1, err := db.Query("Select questionMasterId, questionText, groupQuestionId, questionSubTypeId from questionmaster where questionCategoryId = '" + questionCatVar.QuestionCategoryId + "'")
			for retr_stmt1.Next() {
				err := retr_stmt1.Scan(&quesMVar.QuestionId, &quesMVar.QuestionText, &quesMVar.GroupQuestionId, &quesMVar.QuestionSubTypeId)
				commons.CheckErr(err)
				var temp string
				//var ansMVars []commons.AnsM
				var answersVar commons.Answers
				var answersVars []commons.Answers
				var groupQuestionVars []commons.GroupQuestion

				// retr_stmt2, err := db.Query("Select answerMasterId, answerText from answermaster where questionMasterId = '" + quesMVar.QuestionId + "'")
				// for retr_stmt2.Next() {
				// 	err := retr_stmt2.Scan(&ansMVar.AnswerId, &ansMVar.AnswerText)
				// 	commons.CheckErr(err)
				// 	ansMVar = commons.AnsM{
				// 		AnswerId:   ansMVar.AnswerId,
				// 		AnswerText: ansMVar.AnswerText,
				// 	}
				// 	ansMVars = append(ansMVars, ansMVar)
				// }

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

				err = db.QueryRow("select rfpQuestionId from rfpquestion where questionMasterId = '" + quesMVar.QuestionId + "' and rfpId = '" + RfpId + "'").Scan(&temp)
				if len(temp) != 0 {
					fmt.Println("temp ", temp, quesMVar.QuestionId, RfpId)
					retr_stmt4, err := db.Query("Select answerMasterId from rfpquestionchoices where rfpQuestionId = '" + temp + "'")
					commons.CheckErr(err)
					for retr_stmt4.Next() {
						err := retr_stmt4.Scan(&answersVar.AnswerId)
						fmt.Println("answersVar.AnswerId", answersVar.AnswerId)
						commons.CheckErr(err)
						err = db.QueryRow("Select answerText from answermaster where answerMasterId = '" + answersVar.AnswerId + "'").Scan(&answersVar.Answer)
						answersVar = commons.Answers{
							AnswerId:              answersVar.AnswerId,
							Answer:                answersVar.Answer,
							Priority:              "",
							QuestionSubTypeId:     "",
							GroupQuestionMasterId: "",
						}
						answersVars = append(answersVars, answersVar)
					}
					fmt.Println("///////", answersVars)

					// if quesMVar.QuestionSubTypeId != "10" {

					// 	retr_stmt4, err := db.Query("Select answerId, answer, questionSubTypeId, groupQuestionId from clientanswer where questionMasterId = '" + quesMVar.QuestionId + "' and travelAgencyMasterId = '" + travelAgencyMasterId + "' and groupQuestionId = '0'")
					// 	commons.CheckErr(err)
					// 	for retr_stmt4.Next() {
					// 		err := retr_stmt4.Scan(&answersVar.AnswerId, &answersVar.Answer, &answersVar.QuestionSubTypeId, &answersVar.GroupQuestionMasterId)
					// 		commons.CheckErr(err)
					// 		answersVar = commons.Answers{
					// 			AnswerId:              answersVar.AnswerId,
					// 			Answer:                answersVar.Answer,
					// 			Priority:              "",
					// 			QuestionSubTypeId:     answersVar.QuestionSubTypeId,
					// 			GroupQuestionMasterId: answersVar.GroupQuestionMasterId,
					// 		}
					// 		answersVars = append(answersVars, answersVar)
					// 	}

					// } else {

					// 	fmt.Println("/////")
					// 	retr_stmt5, err := db.Query("Select groupQuestionMasterId from groupquestion where groupQuestionId = '" + quesMVar.GroupQuestionId + "'")
					// 	commons.CheckErr(err)
					// 	for retr_stmt5.Next() {
					// 		err := retr_stmt5.Scan(&quesMVar.QuestionId)
					// 		commons.CheckErr(err)

					// 		retr_stmt6, err := db.Query("Select answerId, answer, questionSubTypeId, groupQuestionId from clientanswer where questionMasterId = '" + quesMVar.QuestionId + "' and travelAgencyMasterId = '" + travelAgencyMasterId + "' and groupQuestionId = '" + quesMVar.GroupQuestionId + "'")
					// 		commons.CheckErr(err)
					// 		for retr_stmt6.Next() {
					// 			err := retr_stmt6.Scan(&answersVar.AnswerId, &answersVar.Answer, &answersVar.QuestionSubTypeId, &answersVar.GroupQuestionMasterId)
					// 			commons.CheckErr(err)
					// 			answersVar = commons.Answers{
					// 				AnswerId:              answersVar.AnswerId,
					// 				Answer:                answersVar.Answer,
					// 				Priority:              "",
					// 				QuestionSubTypeId:     answersVar.QuestionSubTypeId,
					// 				GroupQuestionMasterId: answersVar.GroupQuestionMasterId,
					// 			}
					// 			answersVars = append(answersVars, answersVar)
					// 		}
					// 	}
					// }
					quesMVar = commons.QuesM{
						QuestionId:        quesMVar.QuestionId,
						QuestionText:      quesMVar.QuestionText,
						QuestionSubTypeId: quesMVar.QuestionSubTypeId,
						GroupQuestionId:   quesMVar.GroupQuestionId,
						IsMandatory:       "",
						TabColumn:         groupQuestionVars,
						ConcatAns:         nil,
						Answer:            answersVars,
					}
					quesMVars = append(quesMVars, quesMVar)
				}
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
			QuestionCategoryParentId: parentCatVar.QuestionCategoryParentId,
			QuestionCategoryParent:   parentCatVar.QuestionCategoryParent,
			QuesCategory:             questionCatVars,
		}
		ParentCatsVar = append(ParentCatsVar, parentCatVar)
	}
	RfpFullViewVar := commons.RfpFullView{
		Basic:   BDivisionsVar,
		RfpQues: ParentCatsVar,
	}
	b, err := json.Marshal(RfpFullViewVar)
	commons.CheckErr(err)
	return string(b)
}

func GetRfpResponse(RfpId, HotelId string) string {
	db = GetDB()
	var ParentCatsVar []commons.ParentCat
	retr_stmt7, err := db.Query("Select questionCategoryId, questionCategory from questioncategory where parentId = '0'")
	commons.CheckErr(err)
	for retr_stmt7.Next() {

		var parentCatVar commons.ParentCat
		var questionCatVar commons.QuestionCat
		questionCatVars := make([]commons.QuestionCat, 0)
		var groupQuestionVar commons.GroupQuestion

		var quesMVar commons.QuesM
		//var quesMVars []commons.QuesM
		//var ansMVar commons.AnsM

		//var ansMVars []commons.AnsM
		err = retr_stmt7.Scan(&parentCatVar.QuestionCategoryParentId, &parentCatVar.QuestionCategoryParent)

		retr_stmt8, err := db.Query("Select qm.questionCategoryParentId, qm.questionCategoryId, qq.questionCategory from questionmaster as qm JOIN questioncategory as qq ON qm.questionCategoryId = qq.questionCategoryId where qq.parentId = '" + parentCatVar.QuestionCategoryParentId + "' GROUP BY qq.questionCategoryId")
		commons.CheckErr(err)

		for retr_stmt8.Next() {

			err := retr_stmt8.Scan(&parentCatVar.QuestionCategoryParentId, &questionCatVar.QuestionCategoryId, &questionCatVar.QuestionCategory)
			commons.CheckErr(err)

			quesMVars := make([]commons.QuesM, 0)

			retr_stmt1, err := db.Query("Select questionMasterId, questionText, groupQuestionId, questionSubTypeId from questionmaster where questionCategoryId = '" + questionCatVar.QuestionCategoryId + "'")
			for retr_stmt1.Next() {
				err := retr_stmt1.Scan(&quesMVar.QuestionId, &quesMVar.QuestionText, &quesMVar.GroupQuestionId, &quesMVar.QuestionSubTypeId)
				commons.CheckErr(err)
				var temp string
				//var ansMVars []commons.AnsM
				var answersVar commons.Answers
				var answersVars []commons.Answers
				var groupQuestionVars []commons.GroupQuestion

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
				err = db.QueryRow("select rfpQuestionId from rfpquestion where questionMasterId = '" + quesMVar.QuestionId + "' and rfpId = '" + RfpId + "'").Scan(&temp)
				if len(temp) != 0 {
					fmt.Println("temp ", temp, quesMVar.QuestionId, RfpId)
					retr_stmt4, err := db.Query("Select answer, answerId from clientanswer where questionMasterId = '" + quesMVar.QuestionId + "' and travelAgencyMasterId = '" + HotelId + "'  and groupQuestionId = '0'")
					commons.CheckErr(err)
					for retr_stmt4.Next() {
						err := retr_stmt4.Scan(&answersVar.Answer, &answersVar.AnswerId)
						fmt.Println("answersVar.AnswerId", answersVar.AnswerId)
						commons.CheckErr(err)
						//	err = db.QueryRow("Select answerText from answermaster where answerMasterId = '" + answersVar.AnswerId + "'").Scan(&answersVar.Answer)
						answersVar = commons.Answers{
							AnswerId:              answersVar.AnswerId,
							Answer:                answersVar.Answer,
							Priority:              "",
							QuestionSubTypeId:     "",
							GroupQuestionMasterId: "",
						}
						answersVars = append(answersVars, answersVar)
					}
					fmt.Println("///////", answersVars)

					quesMVar = commons.QuesM{
						QuestionId:        quesMVar.QuestionId,
						QuestionText:      quesMVar.QuestionText,
						QuestionSubTypeId: quesMVar.QuestionSubTypeId,
						GroupQuestionId:   quesMVar.GroupQuestionId,
						IsMandatory:       "",
						TabColumn:         groupQuestionVars,
						ConcatAns:         nil,
						Answer:            answersVars,
					}
					quesMVars = append(quesMVars, quesMVar)
				}
			}
			questionCatVar = commons.QuestionCat{
				QuestionCategoryId: questionCatVar.QuestionCategoryId,
				QuestionCategory:   questionCatVar.QuestionCategory,
				Ques:               quesMVars,
			}
			questionCatVars = append(questionCatVars, questionCatVar)

		}
		parentCatVar = commons.ParentCat{

			Method:                   "",
			TravelAgencyMasterId:     "",
			QuestionCategoryParentId: parentCatVar.QuestionCategoryParentId,
			QuestionCategoryParent:   parentCatVar.QuestionCategoryParent,
			QuesCategory:             questionCatVars,
		}
		ParentCatsVar = append(ParentCatsVar, parentCatVar)
	}
	b, err := json.Marshal(ParentCatsVar)
	commons.CheckErr(err)
	return string(b)
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

		var Answer []commons.Answers
		for retr_stmt1.Next() {
			var InterestVar commons.Answers

			err = retr_stmt1.Scan(&InterestVar.AnswerId)
			commons.CheckErr(err)

			err := db.QueryRow("Select answerText from answermaster where answerMasterId = '" + InterestVar.AnswerId + "'").Scan(&InterestVar.Answer)
			commons.CheckErr(err)

			InterestVar = commons.Answers{
				AnswerId: InterestVar.AnswerId,
				Answer:   InterestVar.Answer,
			}
			Answer = append(Answer, InterestVar)
		}

		quesVar = commons.Ques{
			QuestionId:      quesVar.QuestionId,
			QuestionText:    quesVar.QuestionText,
			GroupQuestionId: quesVar.GroupQuestionId,
			Answer:          Answer,
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

func ListRfpPublished(travelAgencyMasterId string) string {

	db = GetDB()
	retr_stmt, err := db.Query("Select rfpId, rfpName, createDate, completionStatus from rfpmaster where travelAgencyMasterId = '" + travelAgencyMasterId + "' and completionStatus > '0'")
	commons.CheckErr(err)

	var createDate string
	var RfpPublishedVar commons.RfpPublished
	var RfpDetVar commons.RfpDet
	var RfpDetsVar []commons.RfpDet
	var LabValVar commons.LabVal
	var LabValsVar []commons.LabVal
	for retr_stmt.Next() {
		err := retr_stmt.Scan(&RfpDetVar.RfpId, &RfpDetVar.Rfp, &createDate, &RfpDetVar.Connected)
		commons.CheckErr(err)

		err = db.QueryRow("Select COUNT(hotelMasterId) from rfphotelmapping where rfpId = '" + RfpDetVar.RfpId + "'").Scan(&RfpDetVar.NoOfHotels)
		err = db.QueryRow("Select COUNT(hotelMasterId) from rfphotelmapping where rfpId = '" + RfpDetVar.RfpId + "' and slabId != '0'").Scan(&RfpDetVar.NoOfQuotes)
		retr_stmt1, err := db.Query("Select answer, basicAnswerId from basicrfpinfo where rfpId = '" + RfpDetVar.RfpId + "' and basicQuestionId = '5'")
		for retr_stmt1.Next() {
			err := retr_stmt1.Scan(&LabValVar.Label, &LabValVar.Value)
			commons.CheckErr(err)
			LabValVar = commons.LabVal{
				Value: LabValVar.Value,
				Label: LabValVar.Label,
			}
			LabValsVar = append(LabValsVar, LabValVar)
		}
		if RfpDetVar.Connected == "3" {
			RfpDetVar.Connected = "connected"
		} else {
			RfpDetVar.Connected = "pending"
		}
		RfpDetVar = commons.RfpDet{
			Location:   LabValsVar,
			MinPrice:   "",
			NoOfHotels: RfpDetVar.NoOfHotels,
			NoOfQuotes: RfpDetVar.NoOfQuotes,
			Connected:  RfpDetVar.Connected,
			Rfp:        RfpDetVar.Rfp,
			RfpId:      RfpDetVar.RfpId,
		}
		RfpDetsVar = append(RfpDetsVar, RfpDetVar)
	}
	RfpPublishedVar = commons.RfpPublished{
		RfpDets:              RfpDetsVar,
		TravelAgencyMasterId: travelAgencyMasterId,
	}
	b, err := json.Marshal(RfpPublishedVar)
	commons.CheckErr(err)
	return string(b)
}

func ListRfpQuotes(rfpId string) string {
	db = GetDB()
	var ListQuotesVar commons.ListQuotes
	var HotelVar commons.Hotel
	var HotelsVar []commons.Hotel
	err := db.QueryRow("Select rfpName from rfpmaster where rfpId = '" + rfpId + "'").Scan(&ListQuotesVar.RfpName)
	commons.CheckErr(err)
	err = db.QueryRow("Select answer from basicrfpinfo where rfpId = '" + rfpId + "' and basicQuestionId = 13").Scan(&HotelVar.RoomPerMonth)
	commons.CheckErr(err)
	retr_stmt, err := db.Query("Select hotelMasterId, accepted from rfphotelmapping where rfpId = '" + rfpId + "' and slabId != '0'")
	commons.CheckErr(err)

	for retr_stmt.Next() {
		err := retr_stmt.Scan(&HotelVar.HotelId, &HotelVar.Status)
		commons.CheckErr(err)
		var cityId string
		err = db.QueryRow("Select hotelName, cityMasterId from hotelsmaster where hotelsMasterId = '"+HotelVar.HotelId+"'").Scan(&HotelVar.Hotel, cityId)
		err = db.QueryRow("Select cityName from citymaster where cityMasterId = '" + cityId + "'").Scan(&HotelVar.Location)

		if HotelVar.Status == "1" {
			HotelVar.Status = "accepted"
		} else {
			HotelVar.Status = "rejected"
		}

		HotelVar = commons.Hotel{
			Hotel:           HotelVar.Hotel,
			HotelId:         HotelVar.HotelId,
			Location:        HotelVar.Location,
			MaxPrice:        "677777",
			MinPrice:        "9999",
			ProposalMatched: "all",
			Status:          HotelVar.Status,
			RoomPerMonth:    HotelVar.RoomPerMonth,
		}
		HotelsVar = append(HotelsVar, HotelVar)

	}

	ListQuotesVar = commons.ListQuotes{
		RfpId:   rfpId,
		RfpName: ListQuotesVar.RfpName,
		Hotels:  HotelsVar,
	}
	b, err := json.Marshal(ListQuotesVar)
	commons.CheckErr(err)
	return string(b)
}

func AcceptQuote(RfpId, HotelId string) string {
	db = GetDB()

	updatestmt, err := db.Prepare(`UPDATE rfphotelmapping SET accepted = ? where rfpId = '` + RfpId + "' and hotelMasterId = '" + HotelId + "'")
	commons.CheckErr(err)
	res, err := updatestmt.Exec("1")
	commons.CheckErr(err)
	fmt.Println(res)

	updatestmt, err = db.Prepare(`UPDATE rfpmaster SET completionStatus = ? where rfpId = '` + RfpId + "'")
	commons.CheckErr(err)
	res, err = updatestmt.Exec("3")
	commons.CheckErr(err)
	fmt.Println(res)

	var travelAgencyMasterId string
	err = db.QueryRow("Select travelAgencyMasterId from rfpmaster where rfpId = '" + RfpId + "'").Scan(&travelAgencyMasterId)
	commons.CheckErr(err)

	updatestmt, err = db.Prepare(`UPDATE hoteltravelagents SET approved = ? where hotelsMasterId = '` + HotelId + "' and travelAgencyMasterId = '" + travelAgencyMasterId + "'")
	commons.CheckErr(err)
	res, err = updatestmt.Exec("1")
	commons.CheckErr(err)
	fmt.Println(res)

	return "true"
}

//=====================================================================================
// **********************hotel(stage - 3)****************************
//=====================================================================================

func ListRfpByHotel(HotelId string) string {

	db = GetDB()
	retr_stmt, err := db.Query("Select rfpId, travelAgencyMasterId, accepted, slabId from rfphotelmapping where hotelMasterId = '" + HotelId + "' and status = '1'")
	commons.CheckErr(err)

	var rfpId, travelAgencyMasterId, slabId string
	var RfpRecievedVar commons.RfpRecieved
	var CompaniesVar commons.Companies
	var Companies []commons.Companies
	var Locations []commons.LabVal
	for retr_stmt.Next() {
		err := retr_stmt.Scan(&rfpId, &travelAgencyMasterId, &CompaniesVar.Status, slabId)
		commons.CheckErr(err)
		var completion string
		err = db.QueryRow("Select completionStatus from rfpmaster where rfpId = '" + rfpId + "'").Scan(&completion)
		err = db.QueryRow("Select travelAgencyName from travelagencymaster where travelAgencyMasterId = '" + travelAgencyMasterId + "'").Scan(&CompaniesVar.Company.Value)
		err = db.QueryRow("Select rfpName from rfpmaster where rfpId = '" + rfpId + "'").Scan(&CompaniesVar.Rfp.Value)
		err = db.QueryRow("Select answer from basicrfpinfo where rfpId = '" + rfpId + "' and basicQuestionId = 8").Scan(&CompaniesVar.RoomsYear)
		err = db.QueryRow("Select answer from basicrfpinfo where rfpId = '" + rfpId + "' and basicQuestionId = 7").Scan(&CompaniesVar.TravelPerYear)
		err = db.QueryRow("Select answer from basicrfpinfo where rfpId = '" + rfpId + "' and basicQuestionId = 6").Scan(&CompaniesVar.TravelPerMonth)
		fmt.Println("..////....")
		retr_stmt1, err := db.Query("Select basicAnswerId from basicrfpinfo where rfpId = '" + rfpId + "' and basicQuestionId = 5")
		fmt.Println("..////....")

		for retr_stmt1.Next() {
			var cityId, city string
			fmt.Println("../\\///....")

			err := retr_stmt1.Scan(&cityId)
			commons.CheckErr(err)
			fmt.Println(cityId)

			fmt.Println("../\\///....")

			fmt.Println(cityId + "......")
			err = db.QueryRow("Select cityName from citymaster where cityMasterId = '" + cityId + "'").Scan(&city)
			location := commons.LabVal{
				Label: city,
				Value: cityId,
			}
			Locations = append(Locations, location)

		}
		//need to work

		company := commons.LabVal{
			Label: CompaniesVar.Company.Value,
			Value: travelAgencyMasterId,
		}
		rfp := commons.LabVal{
			Label: CompaniesVar.Company.Value,
			Value: CompaniesVar.Rfp.Value,
		}

		if completion == "3" && CompaniesVar.Status == "1" {
			CompaniesVar.Status = "accepted"
		} else if completion == "3" && CompaniesVar.Status == "0" {
			CompaniesVar.Status = "closed"
		} else if completion != "3" && slabId != "0" && CompaniesVar.Status == "0" {
			CompaniesVar.Status = "quoted"
		} else if completion != "3" && slabId == "0" && CompaniesVar.Status == "0" {
			CompaniesVar.Status = "pending"
		}

		CompaniesVar = commons.Companies{
			Company:         company,
			Rfp:             rfp,
			Status:          CompaniesVar.Status,
			RoomsYear:       CompaniesVar.RoomsYear,
			Location:        Locations,
			ProposalMatched: "all",
			TravelPerYear:   CompaniesVar.TravelPerYear,
			TravelPerMonth:  CompaniesVar.TravelPerMonth,
		}
		Companies = append(Companies, CompaniesVar)

	}
	RfpRecievedVar = commons.RfpRecieved{
		Comp: Companies,
	}
	b, err := json.Marshal(RfpRecievedVar)
	commons.CheckErr(err)
	return string(b)

}

func ListSlabs(HotelId string) string {
	db = GetDB()
	retr_stmt, err := db.Query("Select hotelTariffSlabsId, slabName from hoteltariffslabs where hotelsMasterId = '" + HotelId + "'")
	commons.CheckErr(err)
	var LabValsVar []commons.LabVal
	var LabValVar commons.LabVal
	for retr_stmt.Next() {
		err := retr_stmt.Scan(&LabValVar.Value, &LabValVar.Label)
		commons.CheckErr(err)
		LabValVar = commons.LabVal{
			Label: LabValVar.Label,
			Value: LabValVar.Value,
		}
		LabValsVar = append(LabValsVar, LabValVar)
	}
	b, err := json.Marshal(LabValsVar)
	commons.CheckErr(err)
	return string(b)
}

func SendQuote(hotelId, slabId, rfpId string) string {
	db = GetDB()
	//update slabId where rfpId is and hotelId is
	updatestmt, err := db.Prepare(`UPDATE rfphotelmapping SET slabId = ?, viewed = ? where rfpId = '` + rfpId + "' and hotelMasterId = '" + hotelId + "'")
	commons.CheckErr(err)
	res, err := updatestmt.Exec(slabId, "1")
	commons.CheckErr(err)
	var travelAgencyMasterId string
	err = db.QueryRow("Select travelAgencyMasterId from rfpmaster where rfpId = '" + rfpId + "'").Scan(&travelAgencyMasterId)
	id := CheckExistence("hoteltravelagents", "hotelTravelAgentsId", []string{"hotelsMasterId", "travelAgencyMasterId"}, []string{hotelId, travelAgencyMasterId})
	if id != "0" {
		updatestmt2, err := db.Prepare(`UPDATE hoteltravelagents SET hotelTariffSlabsId = ?, approved = ? where hotelsMasterId = '` + hotelId + "' and travelAgencyMasterId = '" + travelAgencyMasterId + "'")
		commons.CheckErr(err)
		res2, err := updatestmt2.Exec(id, "0")
		commons.CheckErr(err)
		fmt.Println(res2)
	} else {
		insrtstmt, err := db.Prepare(`INSERT INTO hoteltravelagents SET hotelsMasterId = ?, hotelTariffSlabsId = ?, travelAgencyMasterId = ?, approved = ?`)
		commons.CheckErr(err)
		res, err := insrtstmt.Exec(hotelId, slabId, travelAgencyMasterId, "0")
		commons.CheckErr(err)
		fmt.Println(res)

	}

	updatestmt1, err := db.Prepare(`UPDATE rfpmaster SET completionStatus = ? where rfpId = '` + rfpId + "'")
	commons.CheckErr(err)
	res1, err := updatestmt1.Exec("2")
	fmt.Println(res, "/", res1)

	commons.CheckErr(err)
	return "true"
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

func CheckDuplicateBasic(RfpId, BqId string) bool {
	db = GetDB()
	var columnvalue string
	err := db.QueryRow("Select COUNT(*) from basicrfpinfo where rfpId = '" + RfpId + "' and basicQuestionId = '" + BqId + "'").Scan(&columnvalue)
	commons.CheckErr(err)
	if columnvalue != "0" {
		return true
	} else {
		return false
	}
}

func CheckExistence(table, hotelTravelAgentsId string, column, value []string) string {
	db = GetDB()
	id := "0"
	err := db.QueryRow("Select " + hotelTravelAgentsId + " from " + table + commons.CreateCondStr(column, value)).Scan(&id)
	commons.CheckErr(err)
	return id
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
		db, err = sql.Open("mysql", "root:@/company_policy?parseTime=true&charset=utf8")
		//	db, err = sql.Open("mysql", "sriram:sriram123@tcp(127.0.0.1:3306)/hotnix_dev?parseTime=true&charset=utf8")

		commons.CheckErr(err)
	}

	return db
}
