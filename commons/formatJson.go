package commons

import "encoding/json"

type AnsM struct {
	AnswerId   string `json:"answerId,omitempty"`
	AnswerText string `json:"answerText,omitempty"`
}

type QuesM struct {
	QuestionId        string `json:"questionId,omitempty"`
	QuestionText      string `json:"questionText,omitempty"`
	QuestionSubTypeId string `json:"questionSubTypeId,omitempty"`
	GroupQuestionId   string `json:"groupQuestionId,omitempty"`
	IsMandatory       string `json:"isMandatory,omitempty"`
	ConcatAns         []AnsM `json:"concatAns,omitempty"`
}

type QuestionCat struct {
	QuestionCategoryId string  `json:"questionCategoryId,omitempty"`
	QuestionCategory   string  `json:"questionCategory,omitempty"`
	Ques               []QuesM `json:"ques,omitempty"`
}

type ParentCat struct {
	QuestionCategoryParentId string        `json:"questionCategoryParentId,omitempty"`
	QuestionCategoryParent   string        `json:"questionCategoryParent,omitempty"`
	QuesCategory             []QuestionCat `json:"quesCategory,omitempty"`
}

type Answers struct {
	Answer   string `json:"answer"`
	AnswerID string `json:"answerId"`
	Priority string `json:"priority"`
}

type Anss struct {
	Answer            []Answers
	GroupQuestionID   string `json:"groupQuestionId"`
	QuestionID        string `json:"questionId"`
	QuestionSubTypeID string `json:"questionSubTypeId"`
}

type MyJsonName struct {
	Ans                  []Anss `json:"ans"`
	ClientTypeMasterID   string `json:"clientTypeMasterId"`
	TravelAgencyMasterID string `json:"travelAgencyMasterId"`
}

func UnmarshalQuestion(jsonStr string) *ParentCat {
	res := &ParentCat{}
	err := json.Unmarshal([]byte(jsonStr), res)
	CheckErr(err)
	//fmt.Println(res)
	return res
}
