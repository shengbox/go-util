package time

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func Today() string {
	return time.Now().Format("2006-01-02")
}

func MsToTimeStr(ms int64) string {
	seconds := ms / 1000
	return SecToStr(seconds)
}

func SecToStr(seconds int64) string {
	minutes := seconds / 60
	seconds %= 60
	if minutes >= 60 {
		hour := minutes / 60
		minutes %= 60
		return fmt.Sprintf("%02d:%02d:%02d", hour, minutes, seconds)
	} else {
		return fmt.Sprintf("%02d:%02d:%02d", 0, minutes, seconds)
	}
}

func TimeStrToSec(timeStr string) int {
	// Split the time string into parts
	parts := strings.Split(timeStr, ":")

	// Check the number of parts
	if len(parts) != 2 && len(parts) != 3 {
		return 0
	}

	// Convert hours, minutes, and seconds to integers
	var hours, minutes, seconds int
	var err error
	if len(parts) == 2 {
		hours = 0
		minutes, err = strconv.Atoi(parts[0])
		if err != nil {
			return 0
		}
		seconds, err = strconv.Atoi(parts[1])
		if err != nil {
			return 0
		}
	} else {
		hours, err = strconv.Atoi(parts[0])
		if err != nil {
			return 0
		}
		minutes, err = strconv.Atoi(parts[1])
		if err != nil {
			return 0
		}
		seconds, err = strconv.Atoi(parts[2])
		if err != nil {
			return 0
		}
	}
	// Calculate total seconds
	totalSeconds := hours*60*60 + minutes*60 + seconds
	return totalSeconds
}
