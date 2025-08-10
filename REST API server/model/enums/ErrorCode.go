package model

// ErrorCode represents a general error
// @Enum
type ErrorCode = int

// List of error codes
// @EnumValuesFor: ErrorCode
type errorCode struct {
	// Undefined [0]
	UNDEFINED ErrorCode `value:"0"`

	// General server error [-1]
	GENERAL_ERROR ErrorCode `value:"-1"`

	// Unauthenticated [-2]
	UNAUTHENTICATED ErrorCode `value:"-2"`

	// Unauthorized [-3]
	UNAUTHORIZED ErrorCode `value:"-3"`

	// Not found [-2]
	NOT_FOUND ErrorCode `value:"-10"`
}

var ErrorCodes = &errorCode{
	UNDEFINED:       0,
	GENERAL_ERROR:   -1,
	UNAUTHENTICATED: -2,
	UNAUTHORIZED:    -3,
	NOT_FOUND:       -10,
}
