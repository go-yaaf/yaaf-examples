package model

// StatusCode represents a general workflow status: PENDING | IN_PROGRESS | DONE ...
// @Enum
type StatusCode = int

// List of workflow status values
// @EnumValuesFor: StatusCode
type statusCode struct {
	// Undefined [0]
	UNDEFINED StatusCode `value:"0"`

	// Flow not started yet [1]
	PENDING StatusCode `value:"1"`

	// Flow in process [2]
	IN_PROCESS StatusCode `value:"2"`

	// Flow completed [3]
	COMPLETED StatusCode `value:"3"`

	// Flow cancelled by user [4]
	CANCELLED StatusCode `value:"4"`

	// Flow automatically cancelled by the system [5]
	AUTO_CANCELLED StatusCode `value:"5"`

	IsValid func(int) bool
	String  func(int) string
}

var StatusCodes = &statusCode{
	UNDEFINED:      0, // Undefined [0]
	PENDING:        1, // Flow not started yet [1]
	IN_PROCESS:     2, // Flow in process [2]
	COMPLETED:      3, // Flow completed [3]
	CANCELLED:      4, // Flow cancelled by user [4]
	AUTO_CANCELLED: 5, // Flow automatically cancelled by the system [5]
	IsValid:        isValidStatusCode,
	String:         stringStatusCode,
}

func isValidStatusCode(code int) bool {
	return code >= 0 && code <= 5
}

var statusCodes = []string{
	"UNDEFINED",
	"PENDING",
	"IN_PROCESS",
	"COMPLETED",
	"CANCELLED",
	"AUTO_CANCELLED",
}

func stringStatusCode(code int) string {
	if isValidStatusCode(code) {
		return statusCodes[code]
	} else {
		return "UNKNOWN"
	}
}
