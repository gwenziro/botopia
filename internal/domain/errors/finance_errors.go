package errors

import "fmt"

// ValidationError merepresentasikan error validasi
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

// NewValidationError membuat error validasi baru
func NewValidationError(field, message string) ValidationError {
	return ValidationError{Field: field, Message: message}
}

// RecordNotFoundError merepresentasikan error record tidak ditemukan
type RecordNotFoundError struct {
	Code string
}

func (e RecordNotFoundError) Error() string {
	return fmt.Sprintf("transaksi dengan kode %s tidak ditemukan", e.Code)
}

// NewRecordNotFoundError membuat error record tidak ditemukan
func NewRecordNotFoundError(code string) RecordNotFoundError {
	return RecordNotFoundError{Code: code}
}

// DuplicateProofError merepresentasikan error bukti transaksi sudah ada
type DuplicateProofError struct {
	Code string
}

func (e DuplicateProofError) Error() string {
	return fmt.Sprintf("transaksi dengan kode %s sudah memiliki bukti", e.Code)
}

// NewDuplicateProofError membuat error bukti transaksi sudah ada
func NewDuplicateProofError(code string) DuplicateProofError {
	return DuplicateProofError{Code: code}
}
