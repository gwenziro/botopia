package dto

import (
	"time"

	"github.com/gwenziro/botopia/internal/domain/finance"
	"github.com/gwenziro/botopia/internal/utils"
)

// FinanceRecordDTO adalah DTO untuk record keuangan
type FinanceRecordDTO struct {
	UniqueCode    string    `json:"uniqueCode"`
	Date          time.Time `json:"date"`
	DateFormatted string    `json:"dateFormatted"`
	Description   string    `json:"description"`
	Amount        float64   `json:"amount"`
	AmountText    string    `json:"amountText"`
	Category      string    `json:"category"`
	PaymentMethod string    `json:"paymentMethod,omitempty"`
	StorageMedia  string    `json:"storageMedia"`
	Notes         string    `json:"notes"`
	ProofURL      string    `json:"proofUrl"`
	HasProof      bool      `json:"hasProof"`
	Type          string    `json:"type"`
	TypeText      string    `json:"typeText"`
}

// FromFinanceRecord mengkonversi domain model ke DTO
func FromFinanceRecord(record *finance.FinanceRecord) *FinanceRecordDTO {
	if record == nil {
		return nil
	}

	typeText := "pengeluaran"
	if record.Type == finance.TypeIncome {
		typeText = "pemasukan"
	}

	return &FinanceRecordDTO{
		UniqueCode:    record.UniqueCode,
		Date:          record.Date,
		DateFormatted: utils.FormatDateID(record.Date),
		Description:   record.Description,
		Amount:        record.Amount,
		AmountText:    utils.FormatMoney(record.Amount),
		Category:      record.Category,
		PaymentMethod: record.PaymentMethod,
		StorageMedia:  record.StorageMedia,
		Notes:         record.Notes,
		ProofURL:      record.ProofURL,
		HasProof:      record.ProofURL != "" && record.ProofURL != "-",
		Type:          string(record.Type),
		TypeText:      typeText,
	}
}
