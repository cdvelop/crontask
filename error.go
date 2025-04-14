package crontask

import (
	"strconv"
)

// errMessage representa un error simple para TinyGo
type errMessage struct {
	message string
}

func (e *errMessage) Error() string {
	return e.message
}

func newErr(args ...any) *errMessage {

	var result string
	var space string

	// Check if we have at least one argument
	if len(args) == 0 {
		return &errMessage{}
	}

	// Process remaining arguments
	for argNumber, arg := range args {
		switch v := arg.(type) {
		case string:
			if v == "" {
				continue
			}
			result += space + v
		case []string:
			for _, s := range v {
				if s == "" {
					continue
				}
				result += space + s
				space = " "
			}
		// Other cases remain the same
		case rune:
			if v == ':' {
				result += ":"
				continue
			}
			result += space + string(v)
		case int:
			result += space + strconv.Itoa(v)
		case float64:
			result += space + strconv.FormatFloat(v, 'f', -1, 64)
		case bool:
			result += space + strconv.FormatBool(v)
		case error:
			result += space + v.Error()
		default:
			result += space + "error not supported arg number: " + strconv.Itoa(argNumber)
		}
		space = " "
	}

	return &errMessage{
		message: result,
	}
}
