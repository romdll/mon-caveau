package activity

import (
	"moncaveau/database"
	"moncaveau/utils"
	"sync"
	"time"
)

type ActivityDoubleBuffer struct {
	mu        sync.Mutex
	active    map[string]database.SessionActivity
	inactive  map[string]database.SessionActivity
	flushChan chan map[string]database.SessionActivity
}

func NewActivityDoubleBuffer() *ActivityDoubleBuffer {
	logger.Info("Initializing ActivityDoubleBuffer")
	return &ActivityDoubleBuffer{
		active:    make(map[string]database.SessionActivity),
		inactive:  make(map[string]database.SessionActivity),
		flushChan: make(chan map[string]database.SessionActivity, 1),
	}
}

func (db *ActivityDoubleBuffer) Write(sessionToken string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.active[sessionToken] = database.SessionActivity{
		SessionToken: sessionToken,
		LastActivity: time.Now(),
	}

	logger.Debugf("Updated activity for session: %s", utils.MaskAll(sessionToken, 14))
}

func (db *ActivityDoubleBuffer) SwapBuffers() {
	db.mu.Lock()
	defer db.mu.Unlock()

	logger.Info("Swapping active and inactive buffers")
	db.active, db.inactive = db.inactive, db.active

	logger.Infof("Flushing %d sessions", len(db.inactive))
	db.flushChan <- db.inactive

	db.inactive = make(map[string]database.SessionActivity)
}
