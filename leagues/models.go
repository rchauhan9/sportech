package leagues

type League struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	NumberOfTeams int32  `json:"numberOfTeams"`
	Country       string `json:"country"`
}
