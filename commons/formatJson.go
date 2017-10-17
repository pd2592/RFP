package commons

import "encoding/json"

type AnsM struct {
	AnswerId   string `json:"answerId,omitempty"`
	AnswerText string `json:"answerText,omitempty"`
}

type GroupQuestion struct {
	GroupQuestionMasterId string `json:"groupQuestionMasterId,omitempty"`
	GroupQuestionId       string `json:"groupQuestionId,omitempty"`
	QuestionText          string `json:"questionText,omitempty"`
	QuestionSubTypeId     string `json:"questionSubTypeId,omitempty"`
}

type QuesM struct {
	QuestionId          string          `json:"questionId,omitempty"`
	QuestionText        string          `json:"questionText,omitempty"`
	QuestionSubTypeId   string          `json:"questionSubTypeId,omitempty"`
	GroupQuestionId     string          `json:"groupQuestionId,omitempty"`
	ConnectedQuestionId string          `json:"connectedQuestionId,omitempty"`
	AnswerMasterId      string          `json:"answerMasterId,omitempty"`
	IsMandatory         string          `json:"isMandatory,omitempty"`
	TabColumn           []GroupQuestion `json:"tabColumn,omitempty"`
	ConcatAns           []AnsM          `json:"concatAns,omitempty"`
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

type ParentSubCat struct {
	QuestionCategoryParentId string  `json:"questionCategoryParentId"`
	QuestionCategoryParent   string  `json:"questionCategoryParent"`
	QuestionCategoryId       string  `json:"questionCategoryId"`
	QuestionCategory         string  `json:"questionCategory"`
	Ques                     []QuesM `json:"ques,omitempty"`
}

type Answers struct {
	AnswerId string `json:"answerId"`
	Answer   string `json:"answer"`
	Priority string `json:"priority"`
}

type Anss struct {
	QuestionId        string    `json:"questionId"`
	QuestionSubTypeId string    `json:"questionSubTypeId"`
	GroupQuestionId   string    `json:"groupQuestionId"`
	Answer            []Answers `json:"answer"`
}

type HotelRes struct {
	TravelAgencyMasterId string `json:"travelAgencyMasterId"`
	ClientTypeMasterId   string `json:"clientTypeMasterId"`
	Ans                  []Anss `json:"ans"`
}

func UnmarshalQuestion(jsonStr string) *ParentCat {
	res := &ParentCat{}
	err := json.Unmarshal([]byte(jsonStr), res)
	CheckErr(err)
	//fmt.Println(res)
	return res
}

func UnmarshalResponse(jsonStr string) *HotelRes {
	res := &HotelRes{}
	err := json.Unmarshal([]byte(jsonStr), res)
	CheckErr(err)
	//fmt.Println(res)
	return res
}
