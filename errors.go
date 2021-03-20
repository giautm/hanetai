package hanetai

import "fmt"

type ServerError struct {
	Code    int
	Message string
}

func (se ServerError) Error() string {
	return fmt.Sprintf("hanet (%d): %s", se.Code, se.Message)
}

const (
	errCodeUnsupported                  = -404
	errCodePersonImgInvalid             = -5010
	errCodeEmployeeIsExists             = -9005
	errCodeEmployeeRegisterImageInvalid = -9006
)

// IsRetriable checks if a given error is an Hanet retriable error
func IsRetriable(err error) bool {
	if err == nil {
		return false
	}

	if e, ok := err.(*ServerError); ok {
		switch e.Code {
		case errCodeUnsupported:
			return false
		case errCodePersonImgInvalid:
			return false
		case errCodeEmployeeIsExists:
			return false
		case errCodeEmployeeRegisterImageInvalid:
			return false
		}
	}

	return true
}
