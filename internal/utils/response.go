package utils

func JSONError(message string) string {
	return `{"message": "` + message + `"}`
}
