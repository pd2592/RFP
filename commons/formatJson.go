package commons

import "encoding/json"

type AnsM struct {
	AnswerId   string `json:"answerId,omitempty"`
	AnswerText string `json:"answerText,omitempty"`
}

type QuesM struct {
	QuestionId   string `json:"questionId,omitempty"`
	QuestionText string `json:"questionText,omitempty"`
	ConcatAns    []AnsM `json:"concatAns,omitempty"`
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

func UnmarshalQuestion(jsonStr string) *ParentCat {
	res := &ParentCat{}
	err := json.Unmarshal([]byte(jsonStr), res)
	CheckErr(err)
	//fmt.Println(res)
	return res
}
