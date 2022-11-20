package stadiums

type Stadium struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Capacity int32  `json:"capacity"`
	City     string `json:"city"`
	Country  string `json:"country"`
}
