package players

import "time"

type Player struct {
	ID               string     `json:"id"`
	FirstName        string     `json:"firstName"`
	MiddleNames      *string    `json:"middleNames"`
	LastName         string     `json:"lastName"`
	DateOfBirth      time.Time  `json:"dateOfBirth"`
	Nationality      string     `json:"nationality"`
	Team             string     `json:"team"`
	SquadNumber      int32      `json:"squadNumber"`
	GeneralPosition  string     `json:"generalPosition"`
	SpecificPosition *string    `json:"specificPosition"`
	Started          time.Time  `json:"started"`
	Ended            *time.Time `json:"ended"`
}

type PlayerDB struct {
	ID               string
	PersonID         string
	TeamID           string
	SquadNumber      int32
	GeneralPosition  string
	SpecificPosition *string
	Started          time.Time
	Ended            *time.Time
}
