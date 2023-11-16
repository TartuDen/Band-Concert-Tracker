package visualization

import (
	"sort"
	"strconv"
	"strings"
)

func SortDates(dates []string) []string {
	sort.Slice(dates, func(i, j int) bool {
		// Split date strings into day, month, and year parts
		dateParts1 := strings.Split(dates[i], "-")
		dateParts2 := strings.Split(dates[j], "-")

		// Convert day, month, and year parts to integers
		day1, _ := strconv.Atoi(dateParts1[0])
		month1, _ := strconv.Atoi(dateParts1[1])
		year1, _ := strconv.Atoi(dateParts1[2])

		day2, _ := strconv.Atoi(dateParts2[0])
		month2, _ := strconv.Atoi(dateParts2[1])
		year2, _ := strconv.Atoi(dateParts2[2])

		// Compare years
		if year1 != year2 {
			return year1 < year2
		}

		// If years are the same, compare months
		if month1 != month2 {
			return month1 < month2
		}

		// If months are the same, compare days
		return day1 < day2
	})

	return dates
}
