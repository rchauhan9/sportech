package managers

import "time"

type Manager struct {
	ID          string     `json:"id"`
	FirstName   string     `json:"firstName"`
	MiddleNames *string    `json:"middleNames"`
	LastName    string     `json:"lastName"`
	DateOfBirth time.Time  `json:"dateOfBirth"`
	Nationality string     `json:"nationality"`
	Team        string     `json:"team"`
	Started     time.Time  `json:"started"`
	Ended       *time.Time `json:"ended"`
}

type ManagerDB struct {
	ID      string
	Person  string
	Team    string
	Started time.Time
	Ended   *time.Time
}
