package utils

func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func ContainsArray(slice []string, items []string) bool {
	for _, item := range items {
		if Contains(slice, item) {
			return true
		}
	}
	return false
}
