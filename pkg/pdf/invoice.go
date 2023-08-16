package pdf

import (
	"time"
)

func dateTimeFromTime(x time.Time) string {
	location, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		return "N/A"
	}
	return x.In(location).Format("02.01.2006 15:04:05")
}
