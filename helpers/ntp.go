package helpers

import (
	"log"
	"time"

	"github.com/beevik/ntp"
)

const host = "0.beevik-ntp.pool.ntp.org"

// GetNTPTime queries the defined NTP host and returns the time object
func GetNTPTime() time.Time {
	time, err := ntp.Time(host)
	if err != nil {
		log.Fatal(err)
	}

	return time
}
