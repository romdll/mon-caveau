package activity

import (
	"moncaveau/utils"
	"os"
	"strconv"
	"time"
)

func SetupActivityBuffer() *ActivityDoubleBuffer {
	logger.Info("Setting up activity buffer")

	intervalMsStr := os.Getenv("ACTIVITY_FLUSH_INTERVAL")
	if intervalMsStr == "" {
		logger.Fatal("ACTIVITY_FLUSH_INTERVAL is not set. Please configure it in your environment.")
	}

	intervalMs, err := strconv.Atoi(intervalMsStr)
	if err != nil || intervalMs <= 0 {
		logger.Fatalf("ACTIVITY_FLUSH_INTERVAL must be a positive integer, got: %s", intervalMsStr)
	}

	flushInterval := time.Duration(intervalMs) * time.Millisecond
	logger.Infof("Flush interval set to: %v", flushInterval)

	flushStopChan := make(chan struct{})
	utils.GracefulShutdownRegistry.Register("activity_flusher", flushStopChan)

	buffer := NewActivityDoubleBuffer()

	go startActivityFlusher(buffer, flushInterval, flushStopChan)

	logger.Info("Activity buffer and flusher initialized")
	return buffer
}
