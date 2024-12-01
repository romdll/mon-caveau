package utils

import "os"

func IsDebugMode() bool {
	return os.Getenv("DEBUG_MODE") == "true"
}
