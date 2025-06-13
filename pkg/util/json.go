package util

func ErrorJSON(message string) map[string]any {
	return map[string]any{
		"error": message,
	}
}

func ResultJSON(message string) map[string]any {
	return map[string]any{
		"result": message,
	}
}
