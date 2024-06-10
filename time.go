package invoke

import "time"

// Time is a package-level variable representing a time handler.
var Time timeHandler

// timeHandler is a struct for time-related operations.
type timeHandler struct{}

// Unix returns the Unix time (seconds since January 1, 1970).
func (timeHandler) Unix() int64 {
	return time.Now().Unix()
}

// UnixNano returns the Unix time in nanoseconds.
func (timeHandler) UnixNano() int64 {
	return time.Now().UnixNano()
}

// Date returns a time corresponding to a specific date.
func (timeHandler) Date(year, month, day, hour, min, sec, nsec int) time.Time {
	return time.Date(year, time.Month(month), day, hour, min, sec, nsec, time.Local)
}

// Parse parses a string into a time.Time object.
func (timeHandler) Parse(layout, value string) (time.Time, error) {
	return time.Parse(layout, value)
}

// Format formats a time.Time object into a string.
func (timeHandler) Format(t time.Time, layout string) string {
	return t.Format(layout)
}

// Add adds a duration to a time.Time object.
func (timeHandler) Add(t time.Time, d time.Duration) time.Time {
	return t.Add(d)
}

// Sub subtracts two time.Time objects to get the duration between them.
func (timeHandler) Sub(t1, t2 time.Time) time.Duration {
	return t1.Sub(t2)
}

// Equal checks if two time.Time objects are equal.
func (timeHandler) Equal(t1, t2 time.Time) bool {
	return t1.Equal(t2)
}

// Before checks if one time is before another.
func (timeHandler) Before(t1, t2 time.Time) bool {
	return t1.Before(t2)
}

// After checks if one time is after another.
func (timeHandler) After(t1, t2 time.Time) bool {
	return t1.After(t2)
}

// Hour returns the hour component of a time.Time object.
func (timeHandler) Hour(t time.Time) int {
	return t.Hour()
}

// Minute returns the minute component of a time.Time object.
func (timeHandler) Minute(t time.Time) int {
	return t.Minute()
}

// Second returns the second component of a time.Time object.
func (timeHandler) Second(t time.Time) int {
	return t.Second()
}

// Day returns the day component of a time.Time object.
func (timeHandler) Day(t time.Time) int {
	return t.Day()
}

// Month returns the month component of a time.Time object.
func (timeHandler) Month(t time.Time) int {
	return int(t.Month())
}

// Year returns the year component of a time.Time object.
func (timeHandler) Year(t time.Time) int {
	return t.Year()
}

// Weekday returns the day of the week (Sunday=0, Monday=1, etc.) of a time.Time object.
func (timeHandler) Weekday(t time.Time) int {
	return int(t.Weekday())
}

// Truncate truncates a time.Time object to a specific precision.
func (timeHandler) Truncate(t time.Time, d time.Duration) time.Time {
	return t.Truncate(d)
}

// Round rounds a time.Time object to the nearest specified duration.
func (timeHandler) Round(t time.Time, d time.Duration) time.Time {
	return t.Round(d)
}

// UTC converts a time.Time object to UTC.
func (timeHandler) UTC(t time.Time) time.Time {
	return t.UTC()
}

// Location returns the location (time zone) of a time.Time object.
func (timeHandler) Location(t time.Time) *time.Location {
	return t.Location()
}

// LoadLocation loads a time zone location by name.
func (timeHandler) LoadLocation(name string) (*time.Location, error) {
	return time.LoadLocation(name)
}

// In changes the time zone of a time.Time object.
func (timeHandler) In(t time.Time, loc *time.Location) time.Time {
	return t.In(loc)
}

// FormatRFC3339 formats a time.Time object in RFC 3339 format.
func (timeHandler) FormatRFC3339(t time.Time) string {
	return t.Format(time.RFC3339)
}

// FormatRFC1123 formats a time.Time object in RFC 1123 format.
func (timeHandler) FormatRFC1123(t time.Time) string {
	return t.Format(time.RFC1123)
}

// FormatUnixDate formats a time.Time object in Unix date format.
func (timeHandler) FormatUnixDate(t time.Time) string {
	return t.Format(time.UnixDate)
}

// FormatKitchen formats a time.Time object in kitchen format.
func (timeHandler) FormatKitchen(t time.Time) string {
	return t.Format(time.Kitchen)
}

// AddDate adds a number of years, months, and days to a time.Time object.
func (timeHandler) AddDate(t time.Time, years, months, days int) time.Time {
	return t.AddDate(years, months, days)
}

// SubtractDate subtracts a number of years, months, and days from a time.Time object.
func (timeHandler) SubtractDate(t time.Time, years, months, days int) time.Time {
	return t.AddDate(-years, -months, -days)
}

// BeginningOfDay returns the start of the day for a given time.Time object.
func (timeHandler) BeginningOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// EndOfDay returns the end of the day for a given time.Time object.
func (timeHandler) EndOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 23, 59, 59, 999999999, t.Location())
}

// BeginningOfMonth returns the start of the month for a given time.Time object.
func (timeHandler) BeginningOfMonth(t time.Time) time.Time {
	y, m, _ := t.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth returns the end of the month for a given time.Time object.
func (th timeHandler) EndOfMonth(t time.Time) time.Time {
	return th.BeginningOfMonth(t).AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// BeginningOfYear returns the start of the year for a given time.Time object.
func (timeHandler) BeginningOfYear(t time.Time) time.Time {
	y, _, _ := t.Date()
	return time.Date(y, 1, 1, 0, 0, 0, 0, t.Location())
}

// EndOfYear returns the end of the year for a given time.Time object.
func (th timeHandler) EndOfYear(t time.Time) time.Time {
	return th.BeginningOfYear(t).AddDate(1, 0, 0).Add(-time.Nanosecond)
}

// BeginningOfWeek returns the start of the week (Sunday) for a given time.Time object.
func (timeHandler) BeginningOfWeek(t time.Time) time.Time {
	weekday := t.Weekday()
	return t.In(time.UTC).AddDate(0, 0, -int(weekday))
}

// EndOfWeek returns the end of the week (Saturday) for a given time.Time object.
func (th timeHandler) EndOfWeek(t time.Time) time.Time {
	return th.BeginningOfWeek(t).AddDate(0, 0, 6).Add(23*time.Hour + 59*time.Minute + 59*time.Second + 999999999*time.Nanosecond)
}

// BeginningOfHour returns the start of the hour for a given time.Time object.
func (timeHandler) BeginningOfHour(t time.Time) time.Time {
	return t.Truncate(time.Hour)
}

// EndOfHour returns the end of the hour for a given time.Time object.
func (th timeHandler) EndOfHour(t time.Time) time.Time {
	return th.BeginningOfHour(t).Add(time.Hour).Add(-time.Nanosecond)
}

// BeginningOfMinute returns the start of the minute for a given time.Time object.
func (timeHandler) BeginningOfMinute(t time.Time) time.Time {
	return t.Truncate(time.Minute)
}

// EndOfMinute returns the end of the minute for a given time.Time object.
func (th timeHandler) EndOfMinute(t time.Time) time.Time {
	return th.BeginningOfMinute(t).Add(time.Minute).Add(-time.Nanosecond)
}

// BeginningOfSecond returns the start of the second for a given time.Time object.
func (th timeHandler) BeginningOfSecond(t time.Time) time.Time {
	return t.Truncate(time.Second)
}

// EndOfSecond returns the end of the second for a given time.Time object.
func (th timeHandler) EndOfSecond(t time.Time) time.Time {
	return th.BeginningOfSecond(t).Add(time.Second).Add(-time.Nanosecond)
}

// BeginningOfDayUTC returns the start of the day in UTC for a given time.Time object.
func (th timeHandler) BeginningOfDayUTC(t time.Time) time.Time {
	return th.BeginningOfDay(th.In(t, time.UTC))
}

// EndOfDayUTC returns the end of the day in UTC for a given time.Time object.
func (th timeHandler) EndOfDayUTC(t time.Time) time.Time {
	return th.EndOfDay(th.In(t, time.UTC))
}

// BeginningOfMonthUTC returns the start of the month in UTC for a given time.Time object.
func (th timeHandler) BeginningOfMonthUTC(t time.Time) time.Time {
	return th.BeginningOfMonth(th.In(t, time.UTC))
}

// EndOfMonthUTC returns the end of the month in UTC for a given time.Time object.
func (th timeHandler) EndOfMonthUTC(t time.Time) time.Time {
	return th.EndOfMonth(th.In(t, time.UTC))
}

// BeginningOfYearUTC returns the start of the year in UTC for a given time.Time object.
func (th timeHandler) BeginningOfYearUTC(t time.Time) time.Time {
	return th.BeginningOfYear(th.In(t, time.UTC))
}

// EndOfYearUTC returns the end of the year in UTC for a given time.Time object.
func (th timeHandler) EndOfYearUTC(t time.Time) time.Time {
	return th.EndOfYear(th.In(t, time.UTC))
}

// BeginningOfWeekUTC returns the start of the week (Sunday) in UTC for a given time.Time object.
func (th timeHandler) BeginningOfWeekUTC(t time.Time) time.Time {
	return th.BeginningOfWeek(th.In(t, time.UTC))
}

// EndOfWeekUTC returns the end of the week (Saturday) in UTC for a given time.Time object.
func (th timeHandler) EndOfWeekUTC(t time.Time) time.Time {
	return th.EndOfWeek(th.In(t, time.UTC))
}

// BeginningOfHourUTC returns the start of the hour in UTC for a given time.Time object.
func (th timeHandler) BeginningOfHourUTC(t time.Time) time.Time {
	return th.BeginningOfHour(th.In(t, time.UTC))
}

// EndOfHourUTC returns the end of the hour in UTC for a given time.Time object.
func (th timeHandler) EndOfHourUTC(t time.Time) time.Time {
	return th.EndOfHour(th.In(t, time.UTC))
}

// BeginningOfMinuteUTC returns the start of the minute in UTC for a given time.Time object.
func (th timeHandler) BeginningOfMinuteUTC(t time.Time) time.Time {
	return th.BeginningOfMinute(th.In(t, time.UTC))
}

// EndOfMinuteUTC returns the end of the minute in UTC for a given time.Time object.
func (th timeHandler) EndOfMinuteUTC(t time.Time) time.Time {
	return th.EndOfMinute(th.In(t, time.UTC))
}

// BeginningOfSecondUTC returns the start of the second in UTC for a given time.Time object.
func (th timeHandler) BeginningOfSecondUTC(t time.Time) time.Time {
	return th.BeginningOfSecond(th.In(t, time.UTC))
}

// EndOfSecondUTC returns the end of the second in UTC for a given time.Time object.
func (th timeHandler) EndOfSecondUTC(t time.Time) time.Time {
	return th.BeginningOfSecond(th.In(t, time.UTC))
}
