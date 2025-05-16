package config

// GetMockDataMaster mengembalikan data master dummy untuk testing
func GetMockDataMaster() map[string][]string {
	return map[string][]string{
		"ExpenseCategories": {
			"Makanan", "Transportasi", "Belanja", "Hiburan", "Kesehatan",
			"Pendidikan", "Listrik", "Internet", "Air", "Sewa", "Lainnya",
		},
		"IncomeCategories": {
			"Gaji", "Bonus", "Hadiah", "Investasi", "Penjualan", "Lainnya",
		},
		"PaymentMethods": {
			"Tunai", "Kartu Debit", "Kartu Kredit", "Transfer Bank", "E-Wallet",
		},
		"StorageMedias": {
			"Dompet", "Bank BCA", "Bank Mandiri", "Bank BNI", "Gopay", "OVO", "Dana",
		},
	}
}

// PopulateWithMockData mengisi konfigurasi dengan data dummy jika kosong
func PopulateWithMockData(config *any) {
	// Implementation would depend on your exact config structure
	// This function would be called when no real data is available
}
