package utils

import "os"

func IsDebugMode() bool {
	return os.Getenv("DEBUG_MODE") == "true"
}

func UseFileSystemForFrontend() bool {
	return os.Getenv("USE_FILESYSTEM_FRONTEND") == "true"
}

func IsHttps() bool {
	return isServerStartedWithHttp
}
