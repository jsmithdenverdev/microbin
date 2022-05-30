package main

import "time"

type pasteExpiration = string

const (
	pasteExpirationOneMin  pasteExpiration = "1min"
	pasteExpirationTenMin  pasteExpiration = "10min"
	pasteExpirationOneHour pasteExpiration = "1hour"
	pasteExpirationOneDay  pasteExpiration = "1day"
	pasteExpirationOneWeek pasteExpiration = "1week"
	pasteExpirationNever   pasteExpiration = "never"
)

var (
	expirations = []pasteExpiration{
		pasteExpirationOneMin,
		pasteExpirationTenMin,
		pasteExpirationOneHour,
		pasteExpirationOneDay,
		pasteExpirationOneWeek,
		pasteExpirationNever,
	}

	expirationDuration = map[pasteExpiration]time.Duration{
		pasteExpirationOneMin:  time.Minute,
		pasteExpirationTenMin:  time.Minute * 10,
		pasteExpirationOneHour: time.Hour,
		pasteExpirationOneDay:  time.Hour * 24,
		pasteExpirationOneWeek: time.Hour * 24 * 7,
		pasteExpirationNever:   -1,
	}
)

func expirationIsValid(s string) bool {
	for _, expiration := range expirations {
		if s == expiration {
			return true
		}
	}

	return false
}
