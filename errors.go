package hanetai

import "fmt"

type ServerError struct {
	Code    int
	Message string

	// Person is set if the error is caused by a duplicated
	Person *Person
}

func (se ServerError) Error() string {
	return fmt.Sprintf("hanet (%d): %s", se.Code, se.Message)
}

const (
	errCodeUnsupported      = -404
	errCodePersonImgInvalid = -5010
	errCodeEmployeeIsExists = -9005
	errCodeInvalidImage     = -9006
	errCodeDuplicatedImage  = -9007
)

// IsRetriable checks if a given error is an Hanet retriable error
func IsRetriable(err error) bool {
	if err == nil {
		return false
	}

	if e, ok := err.(*ServerError); ok {
		switch e.Code {
		case errCodeUnsupported, errCodePersonImgInvalid, errCodeEmployeeIsExists, errCodeInvalidImage, errCodeDuplicatedImage:
			return false
		}
	}

	return true
}
