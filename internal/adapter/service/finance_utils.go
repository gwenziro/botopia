// New file for utility functions
package service

// contains checks if a string exists in a slice
func contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}
