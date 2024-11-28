package utils

import (
	"fmt"
	"log"
	"os"
)

func CreateLogger(packageCaller string) *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf("[MonCaveau] [%s] ", packageCaller), log.LstdFlags|log.Lmicroseconds)
}
