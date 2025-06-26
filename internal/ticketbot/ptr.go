package ticketbot

import "time"

func strToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func intToPtr(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

func timeToPtr(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}
