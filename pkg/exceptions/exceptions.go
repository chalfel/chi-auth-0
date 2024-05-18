package exceptions

import "fmt"

type Exception struct {
	Message    string `json:"message"`
	Constraint string `json:"constraint"`
}

func (e *Exception) Error() string {
	return fmt.Sprintf("error: %s, constraint: %s", e.Message, e.Constraint)
}

func New(message string, constraint string) *Exception {
	return &Exception{
		Message:    message,
		Constraint: constraint,
	}

}
