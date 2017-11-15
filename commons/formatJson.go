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
	Answer              []Answers       `json:"answer"`
}

//store questions details
type QuestionCat struct {
	QuestionCategoryId string  `json:"questionCategoryId,omitempty"`
	QuestionCategory   string  `json:"questionCategory,omitempty"`
	Ques               []QuesM `json:"ques"`
}

//getting questions category wise
type ParentCat struct {
	Method                   string        `json:"method,omitempty"`
	TravelAgencyMasterId     string        `json:"travelAgencyMasterId,omitempty"`
	QuestionCategoryParentId string        `json:"questionCategoryParentId,omitempty"`
	QuestionCategoryParent   string        `json:"questionCategoryParent,omitempty"`
	QuesCategory             []QuestionCat `json:"quesCategory"`
}

//getting question sub category wise
type ParentSubCat struct {
	QuestionCategoryParentId string  `json:"questionCategoryParentId,omitempty"`
	QuestionCategoryParent   string  `json:"questionCategoryParent,omitempty"`
	QuestionCategoryId       string  `json:"questionCategoryId,omitempty"`
	QuestionCategory         string  `json:"questionCategory,omitempty"`
	Ques                     []QuesM `json:"ques,omitempty"`
}

//store answer
type Answers struct {
	AnswerId              string `json:"answerId"`
	Answer                string `json:"answer"`
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
	QuestionId      string         `json:"questionId,omitempty"`
	GroupQuestionId string         `json:"groupQuestionId,omitempty"`
	IsMandatory     string         `json:"isMandatory,omitempty"`
	Answer          []RfpAnsChoice `json:"answer,omitempty"`
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
	QuestionId      string    `json:"questionId,omitempty"`
	QuestionText    string    `json:"questionText,omitempty"`
	GroupQuestionId string    `json:"groupQuestionId,omitempty"`
	Answer          []Answers `json:"answer,omitempty"`
}
type Answer struct {
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
	Division             []BDivision `json:"division,omitempty"`
}

type BDivision struct {
	Division LabVal      `json:"divison,omitempty"`
	Ques     []BQuestion `json:"ques,omitempty"`
}
type BQuestion struct {
	BSubType string   `json:"bSubType,omitempty"`
	BqId     string   `json:"bqId,omitempty"`
	BqText   string   `json:"bqText,omitempty"`
	Answer   string   `json:"answer"`
	AnswerId []LabVal `json:"answerId"`
}

//List Of rfp recieved by hotel
type RfpRecieved struct {
	Comp []Companies `json:"comp,omitempty"`
}

type LabVal struct {
	Label string `json:"label,omitempty"`
	Value string `json:"value,omitempty"`
}

//List of companies
type Companies struct {
	Company         LabVal   `json:"company,omitempty"`
	Rfp             LabVal   `json:"rfp,omitempty"`
	Status          string   `json:"status,omitempty"`
	RfpStatus       string   `json:"rfpStatus,omitempty"`
	RoomsYear       string   `json:"roomsYear,omitempty"`
	Location        []LabVal `json:"location,omitempty"`
	Slab            LabVal   `json:"slab,omitempty"`
	ProposalMatched string   `json:"proposalMatched,omitempty"`
	TravelPerYear   string   `json:"travelPerYear,omitempty"`
	TravelPerMonth  string   `json:"travelPerMonth,omitempty"`
}

type RfpSent struct { //by company
	RfpList []CompRfpList `json:"rfpList,omitempty"`
}

type CompRfpList struct {
	Rfp             LabVal   `json:"rfp,omitempty"`
	RoomsYear       string   `json:"roomsYear,omitempty"`
	Location        []LabVal `json:"location,omitempty"`
	ProposalMatched string   `json:"proposalMatched,omitempty"`
	TravelPerYear   string   `json:"travelPerYear,omitempty"`
	TravelPerMonth  string   `json:"travelPerMonth,omitempty"`
	SentHotelCount  string   `json:"sentHotelCount,omitempty"`
	CreateDate      string   `json:"createDate,omitempty"`
}

type RfpPublished struct {
	RfpDets              []RfpDet `json:"rfp,omitempty"`
	TravelAgencyMasterId string   `json:"travelAgencyMasterId,omitempty"`
}

type RfpDet struct {
	Location   []LabVal `json:"location,omitempty"`
	MinPrice   string   `json:"minPrice,omitempty"`
	NoOfHotels string   `json:"noOfHotels,omitempty"`
	NoOfQuotes string   `json:"noOfQuotes,omitempty"`
	Connected  string   `json:"connected,omitempty"`
	Rfp        string   `json:"rfp,omitempty"`
	RfpId      string   `json:"rfpId,omitempty"`
}

type HotelInfo struct {
	HotelId        string `json:"hotelId,omitempty"`
	HotelName      string `json:"hotelName,omitempty"`
	Address        string `json:"address,omitempty"`
	State          string `json:"state,omitempty"`
	City           string `json:"city,omitempty"`
	Locality       string `json:"locality,omitempty"`
	PrimaryPhone   string `json:"primaryPhone,omitempty"`
	SecondaryPhone string `json:"secondaryPhonerfpId,omitempty"`
	PrimaryMail    string `json:"primaryMail,omitempty"`
	SecondaryMail  string `json:"secondaryMail,omitempty"`
	Web            string `json:"web,omitempty"`
	CheckIn        string `json:"checkIn,omitempty"`
	CheckOut       string `json:"checkOutrfpId,omitempty"`
}

type RfpData struct {
	RfpId   string `json:"rfpId,omitempty"`
	RfpName string `json:"rfpName,omitempty"`
}

type ListQuotes struct {
	RfpId   string  `json:"rfpId,omitempty"`
	RfpName string  `json:"rfpName,omitempty"`
	Hotels  []Hotel `json:"hotels,omitempty"`
}
type Hotel struct {
	Hotel           string `json:"hotel,omitempty"`
	HotelId         string `json:"hotelId,omitempty"`
	Shortlist       string `json:"shortlist,omitempty"`
	Location        string `json:"location,omitempty"`
	MaxPrice        string `json:"maxPrice,omitempty"`
	MinPrice        string `json:"minPrice,omitempty"`
	ProposalMatched string `json:"proposalMatched,omitempty"`
	Status          string `json:"status,omitempty"`
	RoomPerMonth    string `json:"roomPerMonth,omitempty"`
}

type RfpFullView struct {
	Company CompanyInfo `json:"company,omitempty"`
	Rfp     LabVal      `json:"rfp,omitempty"`
	Basic   []BDivision `json:"basic,omitempty"`
	RfpQues []ParentCat `json:"rfpQues,omitempty"`
}

type CompanyInfo struct {
	TravelAgencyMasterId string `json:"travelAgencyMasterId,omitempty"`
	TravelAgencyMaster   string `json:"travelAgencyMaster,omitempty"`
	Email                string `json:"email,omitempty"`
	Phone                string `json:"phone,omitempty"`
	Mobile               string `json:"mobile,omitempty"`
	Website              string `json:"website,omitempty"`
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
