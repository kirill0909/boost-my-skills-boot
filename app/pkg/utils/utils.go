package utils

import (
	"fmt"
	"time"
)

// getFormatedTime returns time in format: 2023/12/31 15:08:21
func GetFormatedTime() string {
	now := time.Now()
	date := []string{
		fmt.Sprintf("%d", now.Year()),
		fmt.Sprintf("%d", now.Month()),
		fmt.Sprintf("%d", now.Day()),
	}

	time := []string{
		fmt.Sprintf("%d", now.Hour()),
		fmt.Sprintf("%d", now.Minute()),
		fmt.Sprintf("%d", now.Second()),
	}

	formatedDate := ""
	for _, d := range date {
		if len(d) == 1 {
			d = fmt.Sprintf("0%s", d)
		}

		if d == date[len(date)-1] || d == fmt.Sprintf("0%s", date[len(date)-1]) {
			formatedDate += fmt.Sprintf("%s", d)
			break
		}

		formatedDate += fmt.Sprintf("%s/", d)
	}

	formatedTime := ""
	for _, t := range time {
		if len(t) == 1 {
			t = fmt.Sprintf("0%s", t)
		}

		if t == time[len(time)-1] || t == fmt.Sprintf("0%s", time[len(time)-1]) {
			formatedTime += fmt.Sprintf("%s", t)
			break
		}

		formatedTime += fmt.Sprintf("%s:", t)
	}

	return fmt.Sprintf("%s %s", formatedDate, formatedTime)
}
