package services

import "errors"

// Sentinel errors for the service layer.
// These are used to provide structured error handling instead of string comparison.

// Attendance & Face Recognition errors
var (
	ErrFaceNotRecognized        = errors.New("face not recognized")
	ErrNoRegisteredFaceImages   = errors.New("no registered face images for this employee")
	ErrEmployeeNotFound         = errors.New("employee not found")
	ErrCompanyNotFound          = errors.New("failed to retrieve company information")
	ErrOutsideAttendanceLocation = errors.New("you are not within a valid attendance location")
	ErrNoShiftAssigned          = errors.New("employee does not have an assigned shift or division shift")
	ErrNoLocationsConfigured    = errors.New("no valid attendance locations configured for employee or division")
	ErrOutsideShiftHours        = errors.New("cannot check-in for regular attendance outside of shift hours. Use overtime check-in instead")
	ErrAlreadyCheckedOut        = errors.New("anda sudah melakukan check-in dan check-out untuk hari ini")
	ErrOnApprovedLeave          = errors.New("employee is on approved leave")
	ErrTooEarlyForCheckIn       = errors.New("anda tidak dapat absen lebih dari 1.5 jam sebelum jam shift Anda")
	ErrFaceRecognitionUnavailable = errors.New("face recognition service is unavailable")

	// Overtime specific errors
	ErrOvertimeDuringShift      = errors.New("cannot check-in for overtime during regular shift hours")
	ErrAlreadyCheckedInOvertime = errors.New("employee is already checked in for overtime")
	ErrNotCheckedInForOvertime  = errors.New("employee is not currently checked in for overtime")
	ErrMustCheckOutRegular      = errors.New("anda harus check-out dari shift reguler sebelum check-in lembur")

	// General errors
	ErrInvalidTimezone          = errors.New("invalid company timezone configuration")
	ErrFaceImageRetrieval       = errors.New("could not retrieve employee face image")
	ErrAttendanceRetrieval      = errors.New("failed to retrieve attendance record")
	ErrLocationRetrieval        = errors.New("failed to retrieve company attendance locations")
	ErrLeaveCheckFailed         = errors.New("failed to check leave status")
	ErrShiftValidationFailed    = errors.New("failed to validate shift time")
)
