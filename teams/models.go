package teams

type Team struct {
	ID          string `json:"id"`
	FullName    string `json:"fullName"`
	MediumName  string `json:"mediumName"`
	Acronym     string `json:"acronym"`
	Nickname    string `json:"nickname"`
	YearFounded int32  `json:"yearFounded"`
	City        string `json:"city"`
	Country     string `json:"country"`
	Stadium     string `json:"stadium"`
	League      string `json:"league"`
}
