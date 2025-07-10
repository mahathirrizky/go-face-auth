package helper

import (
	"fmt"
	"time"
)

// ParseTime parses a time string (HH:MM:SS) into a time.Time object for a given date and location.
func ParseTime(date time.Time, timeStr string, loc *time.Location) (time.Time, error) {
	layout := "15:04:05"
	parsedTime, err := time.Parse(layout, timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time %s: %w", timeStr, err)
	}
	// Combine the date from 'date' with the time from 'parsedTime' in the specified location
	return time.Date(date.Year(), date.Month(), date.Day(), parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), 0, loc), nil
}

// CalculateShiftDuration calculates the duration of a shift in minutes.
// It handles shifts that cross midnight.
func CalculateShiftDuration(startTimeStr, endTimeStr string) (time.Duration, error) {
	// Use a dummy date for parsing to calculate duration
	today := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

	start, err := ParseTime(today, startTimeStr, time.UTC)
	if err != nil {
		return 0, err
	}

	end, err := ParseTime(today, endTimeStr, time.UTC)
	if err != nil {
		return 0, err
	}

	// If end time is before start time, it means the shift crosses midnight
	if end.Before(start) {
		end = end.Add(24 * time.Hour)
	}

	return end.Sub(start), nil
}

// IsTimeWithinShift checks if a given time falls within a shift's start and end times,
// considering a grace period for the start time.
// It handles shifts that cross midnight and uses the provided location for shift times.
func IsTimeWithinShift(checkTime time.Time, shiftStartTimeStr, shiftEndTimeStr string, gracePeriodMinutes int, loc *time.Location) (bool, error) {
	// Use the date of checkTime for parsing shift times in the specified location
	shiftStart, err := ParseTime(checkTime, shiftStartTimeStr, loc)
	if err != nil {
		return false, err
	}
	shiftEnd, err := ParseTime(checkTime, shiftEndTimeStr, loc)
	if err != nil {
		return false, err
	}

	// Adjust shiftEnd if it crosses midnight
	if shiftEnd.Before(shiftStart) {
		shiftEnd = shiftEnd.Add(24 * time.Hour)
		// If checkTime is before shiftStart, it means checkTime is on the next day relative to shiftStart
		if checkTime.Before(shiftStart) {
			checkTime = checkTime.Add(24 * time.Hour)
		}
	}

	// Apply grace period to shift start time
	shiftStartWithGrace := shiftStart.Add(-time.Duration(gracePeriodMinutes) * time.Minute)

	// Check if the time is within the adjusted shift period
	return (checkTime.After(shiftStartWithGrace) || checkTime.Equal(shiftStartWithGrace)) && checkTime.Before(shiftEnd), nil
}
