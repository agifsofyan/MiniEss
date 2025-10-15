package utils

import "time"

type countdown struct {
	t int
	d int
	h int
	m int
	s int
}

var countryTz = map[string]string{
	"Hungary":   "Europe/Budapest",
	"Egypt":     "Africa/Cairo",
	"Indonesia": "Asia/Jakarta",
}

func TimeIn(t time.Time, name string) (time.Time, error) {
	loc, err := time.LoadLocation(countryTz[name])
	if err != nil {
		return time.Time{}, err
	}
	return t.In(loc), nil
}

func IdnTime(t time.Time) (time.Time, error) {
	return TimeIn(t, "Indonesia")
}

func CreateTimeIdn(t time.Time) time.Time {
	return t.UTC().Add(time.Hour * 7)
}

func GetTimeRemaining(t time.Time) countdown {
	currentTime := CreateTimeIdn(time.Now())
	difference := t.Sub(currentTime)

	total := int(difference.Seconds())
	days := int(total / (60 * 60 * 24))
	hours := int(total / (60 * 60) % 24)
	minutes := int(total/60) % 60
	seconds := int(total % 60)

	return countdown{
		t: total,
		d: days,
		h: hours,
		m: minutes,
		s: seconds,
	}
}

func IsAfterOrEqualTimeInWIB(t time.Time, hour int, toleranceMinute int) bool {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// fallback jika tzdata tidak ada di container
		loc = time.FixedZone("WIB", 7*60*60)
	}

	var now time.Time
	if t.IsZero() {
		now = time.Now().In(loc)
	} else {
		now = t.In(loc)
	}

	// buat waktu target jam 09:00 WIB hari ini
	target := time.Date(now.Year(), now.Month(), now.Day(), hour, toleranceMinute, 0, 0, loc)

	// cek apakah sekarang >= 09:00 WIB
	return target.After(now)
}

func IsWeekendWIB(t time.Time) bool {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// fallback timezone manual (UTC+7)
		loc = time.FixedZone("WIB", 7*60*60)
	}

	var now time.Time
	if t.IsZero() {
		now = time.Now().In(loc)
	} else {
		now = t.In(loc)
	}

	weekday := now.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}
