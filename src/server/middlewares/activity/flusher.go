package activity

import (
	"moncaveau/database"
	"time"
)

func startActivityFlusher(buffer *ActivityDoubleBuffer, flushInterval time.Duration, serverStopChan chan struct{}) {
	logger.Infof("Starting activity flusher with interval: %v", flushInterval)

	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logger.Info("Flusher tick: Swapping buffers")
			buffer.SwapBuffers()

		case sessions := <-buffer.flushChan:
			if len(sessions) > 0 {
				logger.Infof("Flushing %d session updates to database", len(sessions))
				if err := database.FlushActivityUpdate(sessions); err != nil {
					logger.Errorf("Failed to flush activity to database: %v", err)
				} else {
					logger.Info("Successfully flushed session updates to database")
				}
			} else {
				logger.Info("No session updates to flush")
			}

		case <-serverStopChan:
			logger.Info("Received stop signal, stopping activity flusher...")
			return
		}
	}
}
