package microbin

import "time"

type Expiration string

const (
	ExpirationOneMin  Expiration = "1min"
	ExpirationTenMin  Expiration = "10min"
	ExpirationOneHour Expiration = "1hour"
	ExpirationOneDay  Expiration = "1day"
	ExpirationOneWeek Expiration = "1week"
	ExpirationNever   Expiration = "never"
)

func (e Expiration) IsValid() bool {
	for _, expiration := range []Expiration{
		ExpirationOneMin,
		ExpirationTenMin,
		ExpirationOneHour,
		ExpirationOneDay,
		ExpirationOneWeek,
		ExpirationNever,
	} {
		if e == expiration {
			return true
		}
	}

	return false
}

func (e Expiration) ToDuration() time.Duration {
	expirationDuration := map[Expiration]time.Duration{
		ExpirationOneMin:  time.Minute,
		ExpirationTenMin:  time.Minute * 10,
		ExpirationOneHour: time.Hour,
		ExpirationOneDay:  time.Hour * 24,
		ExpirationOneWeek: time.Hour * 24 * 7,
		ExpirationNever:   -1,
	}

	return expirationDuration[e]
}
