package finance

// Configuration menyimpan konfigurasi untuk fitur keuangan
type Configuration struct {
	Year              int
	Month             int
	StorageMedias     []string
	PaymentMethods    []string
	ExpenseCategories []string
	IncomeCategories  []string
}
