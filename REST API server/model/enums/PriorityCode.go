package model

// PriorityCode represents a priority: LOW | MEDIUM | HIGH
// @Enum
type PriorityCode = int

// PriorityCodes the list of priority values
// @EnumValuesFor: PriorityCode
type priorityCode struct {
	// Undefined [0]
	UNDEFINED PriorityCode `value:"0"`

	// No priority [1]
	NONE PriorityCode `value:"1"`

	// Low priority [2]
	LOW PriorityCode `value:"2"`

	// Medium priority [3]
	MEDIUM PriorityCode `value:"3"`

	// High priority [4]
	HIGH PriorityCode `value:"4"`

	IsValid func(int) bool
	String  func(int) string
}

var PriorityCodes = &priorityCode{
	UNDEFINED: 0, // Undefined [0]
	NONE:      1, // No priority [1]
	LOW:       2, // Low priority [2]
	MEDIUM:    3, // Medium priority [3]
	HIGH:      4, // High priority [4]

	IsValid: isValidPriorityCode,
	String:  stringPriorityCode,
}

func isValidPriorityCode(code int) bool {
	return code >= 0 && code <= 4
}

var priorityCodes = []string{
	"UNDEFINED",
	"NONE",
	"LOW",
	"MEDIUM",
	"HIGH",
}

func stringPriorityCode(code int) string {
	if isValidPriorityCode(code) {
		return priorityCodes[code]
	} else {
		return "UNKNOWN"
	}
}
