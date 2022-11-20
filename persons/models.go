package persons

import "time"

type Person struct {
	ID          string
	FirstName   string
	MiddleNames *string
	LastName    string
	DateOfBirth time.Time
	Nationality string
}
