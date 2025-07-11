package types

import "time"

type TimeDetails struct {
	UpdatedAt time.Time
}

type Ticket struct {
	ID      int    // required
	Summary string // required
	TimeDetails
}
