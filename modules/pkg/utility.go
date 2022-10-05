package pkg

import "time"

func FormatDateQuery(timeFormat, date string) (string, error) {
	t := time.Now()
	dateQuery := ""
	if date == "" {
		year, month, day := t.Date()
		dateQuery = time.Date(year, month, day, 0, 0, 0, 0, t.Location()).Format(timeFormat)
	} else {
		start, err := time.Parse("20060102150405", date)
		if err != nil {
			return "", err
		}
		dateQuery = start.Format(timeFormat)
	}
	return dateQuery, nil
}
