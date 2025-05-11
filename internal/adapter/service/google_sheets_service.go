package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gwenziro/botopia/internal/domain/repository"
	"github.com/gwenziro/botopia/internal/domain/service"
	"github.com/gwenziro/botopia/internal/infrastructure/config"
	"github.com/gwenziro/botopia/internal/infrastructure/logger"

	"google.golang.org/api/sheets/v4"
)

// GoogleSheetsService implementasi service Google Sheets
type GoogleSheetsService struct {
	googleAPIRepo repository.GoogleAPIRepository
	config        *config.GoogleSheetsConfig
	log           *logger.Logger
}

// NewGoogleSheetsService membuat instance service baru
func NewGoogleSheetsService(
	googleAPIRepo repository.GoogleAPIRepository,
	config *config.Config,
	log *logger.Logger,
) *GoogleSheetsService {
	return &GoogleSheetsService{
		googleAPIRepo: googleAPIRepo,
		config:        config.GoogleSheets,
		log:           log,
	}
}

// Konstanta nama sheet
const (
	ExpenseSheetName = "Pengeluaran"
	IncomeSheetName  = "Pemasukan"
	SummarySheetName = "Konfigurasi"
)

// AddExpense menambahkan data pengeluaran ke spreadsheet
func (s *GoogleSheetsService) AddExpense(
	amount float64,
	category string,
	description string,
	imageURL string,
) (*service.FinanceRecord, error) {
	ctx := context.Background()
	sheetsService, err := s.googleAPIRepo.GetSheetsService(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan sheets service: %v", err)
	}

	// Buat record
	now := time.Now()
	id := fmt.Sprintf("EXP-%s", now.Format("20060102-150405"))
	record := &service.FinanceRecord{
		ID:          id,
		Date:        now,
		Amount:      amount,
		Category:    category,
		Description: description,
		IsIncome:    false,
		ImageURL:    imageURL,
	}

	// Siapkan data untuk spreadsheet
	values := []interface{}{
		id,
		now.Format("2006-01-02 15:04:05"),
		amount,
		category,
		description,
		imageURL,
	}

	// Tambahkan ke spreadsheet
	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{values},
	}

	_, err = sheetsService.Spreadsheets.Values.Append(
		s.config.SpreadsheetID,
		ExpenseSheetName+"!A:F",
		valueRange,
	).ValueInputOption("USER_ENTERED").Do()

	if err != nil {
		return nil, fmt.Errorf("gagal menambah data: %v", err)
	}

	// Update saldo di sheet ringkasan
	s.updateBalance(ctx, sheetsService)

	return record, nil
}

// AddIncome menambahkan data pemasukan ke spreadsheet
func (s *GoogleSheetsService) AddIncome(
	amount float64,
	category string,
	description string,
	imageURL string,
) (*service.FinanceRecord, error) {
	ctx := context.Background()
	sheetsService, err := s.googleAPIRepo.GetSheetsService(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan sheets service: %v", err)
	}

	// Buat record
	now := time.Now()
	id := fmt.Sprintf("INC-%s", now.Format("20060102-150405"))
	record := &service.FinanceRecord{
		ID:          id,
		Date:        now,
		Amount:      amount,
		Category:    category,
		Description: description,
		IsIncome:    true,
		ImageURL:    imageURL,
	}

	// Siapkan data untuk spreadsheet
	values := []interface{}{
		id,
		now.Format("2006-01-02 15:04:05"),
		amount,
		category,
		description,
		imageURL,
	}

	// Tambahkan ke spreadsheet
	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{values},
	}

	_, err = sheetsService.Spreadsheets.Values.Append(
		s.config.SpreadsheetID,
		IncomeSheetName+"!A:F",
		valueRange,
	).ValueInputOption("USER_ENTERED").Do()

	if err != nil {
		return nil, fmt.Errorf("gagal menambah data: %v", err)
	}

	// Update saldo di sheet ringkasan
	s.updateBalance(ctx, sheetsService)

	return record, nil
}

// GetLatestRecords mendapatkan 5 record keuangan terakhir
func (s *GoogleSheetsService) GetLatestRecords(limit int) ([]*service.FinanceRecord, error) {
	ctx := context.Background()
	sheetsService, err := s.googleAPIRepo.GetSheetsService(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan sheets service: %v", err)
	}

	var records []*service.FinanceRecord

	// Ambil 5 data terakhir pengeluaran
	expenseResp, err := sheetsService.Spreadsheets.Values.Get(
		s.config.SpreadsheetID,
		fmt.Sprintf("%s!A2:F", ExpenseSheetName),
	).Do()

	if err != nil {
		return nil, fmt.Errorf("gagal membaca data pengeluaran: %v", err)
	}

	if len(expenseResp.Values) > 0 {
		for i := len(expenseResp.Values) - 1; i >= 0 && len(records) < limit; i-- {
			row := expenseResp.Values[i]
			if len(row) < 5 {
				continue
			}

			date, _ := time.Parse("2006-01-02 15:04:05", row[1].(string))
			amount, _ := strconv.ParseFloat(row[2].(string), 64)

			imageURL := ""
			if len(row) >= 6 {
				imageURL = row[5].(string)
			}

			records = append(records, &service.FinanceRecord{
				ID:          row[0].(string),
				Date:        date,
				Amount:      amount,
				Category:    row[3].(string),
				Description: row[4].(string),
				IsIncome:    false,
				ImageURL:    imageURL,
			})
		}
	}

	// Ambil 5 data terakhir pemasukan
	incomeResp, err := sheetsService.Spreadsheets.Values.Get(
		s.config.SpreadsheetID,
		fmt.Sprintf("%s!A2:F", IncomeSheetName),
	).Do()

	if err != nil {
		return nil, fmt.Errorf("gagal membaca data pemasukan: %v", err)
	}

	if len(incomeResp.Values) > 0 {
		for i := len(incomeResp.Values) - 1; i >= 0 && len(records) < limit; i-- {
			row := incomeResp.Values[i]
			if len(row) < 5 {
				continue
			}

			date, _ := time.Parse("2006-01-02 15:04:05", row[1].(string))
			amount, _ := strconv.ParseFloat(row[2].(string), 64)

			imageURL := ""
			if len(row) >= 6 {
				imageURL = row[5].(string)
			}

			records = append(records, &service.FinanceRecord{
				ID:          row[0].(string),
				Date:        date,
				Amount:      amount,
				Category:    row[3].(string),
				Description: row[4].(string),
				IsIncome:    true,
				ImageURL:    imageURL,
			})
		}
	}

	return records, nil
}

// GetBalance mendapatkan saldo terkini
func (s *GoogleSheetsService) GetBalance() (float64, error) {
	ctx := context.Background()
	sheetsService, err := s.googleAPIRepo.GetSheetsService(ctx)
	if err != nil {
		return 0, fmt.Errorf("gagal mendapatkan sheets service: %v", err)
	}

	// Baca saldo dari sheet ringkasan
	resp, err := sheetsService.Spreadsheets.Values.Get(
		s.config.SpreadsheetID,
		fmt.Sprintf("%s!B1", SummarySheetName),
	).Do()

	if err != nil {
		return 0, fmt.Errorf("gagal membaca data saldo: %v", err)
	}

	if len(resp.Values) == 0 || len(resp.Values[0]) == 0 {
		return 0, nil
	}

	balance, err := strconv.ParseFloat(resp.Values[0][0].(string), 64)
	if err != nil {
		return 0, fmt.Errorf("gagal parse saldo: %v", err)
	}

	return balance, nil
}

// GetSheetURL mendapatkan URL spreadsheet
func (s *GoogleSheetsService) GetSheetURL() string {
	return "https://docs.google.com/spreadsheets/d/" + s.config.SpreadsheetID
}

// updateBalance memperbarui saldo di sheet ringkasan
// ctx parameter digunakan untuk operasi sheets service
func (s *GoogleSheetsService) updateBalance(_ context.Context, sheetsService *sheets.Service) error {
	// Ambil total pemasukan
	incomeResp, err := sheetsService.Spreadsheets.Values.Get(
		s.config.SpreadsheetID,
		fmt.Sprintf("%s!C2:C", IncomeSheetName),
	).Do()
	if err != nil {
		return err
	}

	totalIncome := 0.0
	if len(incomeResp.Values) > 0 {
		for _, row := range incomeResp.Values {
			if len(row) > 0 {
				amount, _ := strconv.ParseFloat(row[0].(string), 64)
				totalIncome += amount
			}
		}
	}

	// Ambil total pengeluaran
	expenseResp, err := sheetsService.Spreadsheets.Values.Get(
		s.config.SpreadsheetID,
		fmt.Sprintf("%s!C2:C", ExpenseSheetName),
	).Do()
	if err != nil {
		return err
	}

	totalExpense := 0.0
	if len(expenseResp.Values) > 0 {
		for _, row := range expenseResp.Values {
			if len(row) > 0 {
				amount, _ := strconv.ParseFloat(row[0].(string), 64)
				totalExpense += amount
			}
		}
	}

	// Hitung saldo
	balance := totalIncome - totalExpense

	// Update saldo di sheet ringkasan
	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{{"Saldo saat ini:", balance}},
	}

	_, err = sheetsService.Spreadsheets.Values.Update(
		s.config.SpreadsheetID,
		fmt.Sprintf("%s!A1:B1", SummarySheetName),
		valueRange,
	).ValueInputOption("USER_ENTERED").Do()

	return err
}
