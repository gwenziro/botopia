package errors

import "fmt"

// TimeoutError adalah error ketika operasi melebihi batas waktu
type TimeoutError struct {
	Message string
}

func (e TimeoutError) Error() string {
	return fmt.Sprintf("timeout error: %s", e.Message)
}

// NewTimeoutError membuat timeout error baru
func NewTimeoutError(message string) TimeoutError {
	return TimeoutError{Message: message}
}

// RepositoryError adalah error yang terjadi pada repository
type RepositoryError struct {
	Operation string
	Message   string
	Cause     error
}

func (e RepositoryError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("repository error pada %s: %s - caused by: %v",
			e.Operation, e.Message, e.Cause)
	}
	return fmt.Sprintf("repository error pada %s: %s", e.Operation, e.Message)
}

// NewRepositoryError membuat repository error baru
func NewRepositoryError(operation, message string, cause error) RepositoryError {
	return RepositoryError{
		Operation: operation,
		Message:   message,
		Cause:     cause,
	}
}
