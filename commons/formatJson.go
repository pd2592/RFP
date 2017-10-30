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
	Answer              []Answers       `json:"answer,omitempty"`
}

//store questions details
type QuestionCat struct {
	QuestionCategoryId string  `json:"questionCategoryId,omitempty"`
	QuestionCategory   string  `json:"questionCategory,omitempty"`
	Ques               []QuesM `json:"ques,omitempty"`
}

//getting questions category wise
type ParentCat struct {
	Method                   string        `json:"method,omitempty"`
	TravelAgencyMasterId     string        `json:"travelAgencyMasterId,omitempty"`
	QuestionCategoryParentId string        `json:"questionCategoryParentId,omitempty"`
	QuestionCategoryParent   string        `json:"questionCategoryParent,omitempty"`
	QuesCategory             []QuestionCat `json:"quesCategory,omitempty"`
}

//getting question sub category wise
type ParentSubCat struct {
	QuestionCategoryParentId string  `json:"questionCategoryParentId,omitempty"`
	QuestionCategoryParent   string  `json:"questionCategoryParent,omitempty"`
	QuestionCategoryId       string  `json:"questionCategoryId,omitempty"`
	QuestionCategory         string  `json:"questionCategory,omitempty"`
	Ques                     []QuesM `json:"ques,omitempty,omitempty"`
}

//store answer
type Answers struct {
	AnswerId              string `json:"answerId,omitempty"`
	Answer                string `json:"answer,omitempty"`
	Priority              string `json:"priority,omitempty"`
	QuestionSubTypeId     string `json:"questionSubTypeId,omitempty"`
	GroupQuestionMasterId string `json:"groupQuestionMasterId,omitempty"`
}

//store question id and answer details
type Anss struct {
	QuestionId        string    `json:"questionId,omitempty"`
	QuestionSubTypeId string    `json:"questionSubTypeId,omitempty"`
	GroupQuestionId   string    `json:"groupQuestionId,omitempty"`
	Answer            []Answers `json:"answer,omitempty"`
}

//store hotel response
type HotelRes struct {
	TravelAgencyMasterId string `json:"travelAgencyMasterId,omitempty"`
	ClientTypeMasterId   string `json:"clientTypeMasterId,omitempty"`
	Ans                  []Anss `json:"ans,omitempty"`
}

type ParentCategory struct {
	QuestionCategoryId string `json:"questionCategoryId,omitempty"`
	QuestionCategory   string `json:"questionCategory,omitempty"`
}

//rfp answers choices table : rfpQuestionChoices
type RfpAnsChoice struct {
	AnswerId string `json:"answerId,omitempty"`
	Priority string `json:"priority,omitempty"`
}

//rfp question choices, table : rfpQuestions
type RfpQuesChoice struct {
	QuestionMasterId string         `json:"questionMasterId,omitempty"`
	GroupQuestionId  string         `json:"groupQuestionId,omitempty"`
	IsMandatory      string         `json:"isMandatory,omitempty"`
	Interests        []RfpAnsChoice `json:"interests,omitempty"`
}

//customized questions
type RfpCustomise struct {
	CustomiseQuestionId string `json:"customiseQuestionId,omitempty"`
	QuestionText        string `json:"questionText,omitempty"`
	QuestionCategoryId  string `json:"questionCategoryId,omitempty"`
}

// create rfp, table : rfpMaster
type Rfp struct {
	RfpName              string          `json:"rfpName,omitempty"`
	RfpId                string          `json:"rfpId,omitempty"`
	TravelAgencyMasterId string          `json:"travelAgencyMasterId,omitempty"`
	RefRfpId             string          `json:"refRfpId,omitempty"`
	Status               string          `json:"status,omitempty"`
	Ques                 []RfpQuesChoice `json:"ques,omitempty"`
	CustomiseQues        []RfpCustomise  `json:"customiseQues,omitempty"`
}

//list hotels to send rfp
type ListHotel struct {
	HotelName        string `json:"hotelName,omitempty"`
	HotelId          string `json:"hotelId,omitempty"`
	Star             string `json:"star,omitempty"`
	Locality         string `json:"locality,omitempty"`
	City             string `json:"city,omitempty"`
	DistanceFromCity string `json:"distanceFromCity,omitempty"`
}

//preview RFP
type RfpView struct {
	RfpID string `json:"rfpId,omitempty"`
	Ques  []Ques `json:"ques,omitempty"`
}
type Ques struct {
	QuestionId      string     `json:"questionId,omitempty"`
	QuestionText    string     `json:"questionText,omitempty"`
	GroupQuestionId string     `json:"groupQuestionId,omitempty"`
	Interests       []Interest `json:"interests,omitempty"`
}
type Interest struct {
	Answer   string `json:"answer,omitempty"`
	AnswerId string `json:"answerId,omitempty"`
}

//sending rfp to listed hotels
type RfpSend struct {
	Hotels               []string `json:"hotels,omitempty"`
	RfpId                string   `json:"rfpId,omitempty"`
	TravelAgencyMasterId string   `json:"travelAgencyMasterId,omitempty"`
}

//Creating and saving basic question
type BasicQuestion struct {
	RfpId                string      `json:"rfpId,omitempty"`
	RfpName              string      `json:"rfpName,omitempty"`
	TravelAgencyMasterId string      `json:"travelAgencyMasterId,omitempty"`
	Ques                 []BQuestion `json:"ques,omitempty"`
}

type BQuestion struct {
	BSubType string `json:"bSubType,omitempty"`
	BqId     string `json:"bqId,omitempty"`
	BqText   string `json:"bqText,omitempty"`
	Divison  string `json:"divison,omitempty"`
	Answer   string `json:"answer,omitempty"`
	AnswerId string `json:"answerId,omitempty"`
}

func UnmarshalRFPBasic(jsonStr string) *BasicQuestion {
	res := &BasicQuestion{}
	err := json.Unmarshal([]byte(jsonStr), res)
	CheckErr(err)
	//fmt.Println(res)
	return res
}

func UnmarshalQuestion(jsonStr string) *ParentCat {
	res := &ParentCat{}
	err := json.Unmarshal([]byte(jsonStr), res)
	CheckErr(err)
	//fmt.Println(res)
	return res
}

//hotel response
func UnmarshalResponse(jsonStr string) *HotelRes {
	res := &HotelRes{}
	err := json.Unmarshal([]byte(jsonStr), res)
	CheckErr(err)
	//fmt.Println(res)
	return res
}

func UnmarshalRFP(jsonStr string) *Rfp {
	res := &Rfp{}
	err := json.Unmarshal([]byte(jsonStr), res)
	CheckErr(err)
	//fmt.Println(res)
	return res
}

func UnmarshalRFPSend(jsonStr string) *RfpSend {
	res := &RfpSend{}
	err := json.Unmarshal([]byte(jsonStr), res)
	CheckErr(err)
	//fmt.Println(res)
	return res
}
