package stomion

import (
	"sync/atomic"
	"time"

	pb "github.com/odysseia-greek/agora/eupalinos/proto"
)

// Diexodos represents a task queue
type Diexodos struct {
	LastMessageReceived time.Time
	Name                string
	InternalID          string
	MessageQueue        map[string]*pb.InternalEpistello
	MessageUpdateCh     chan pb.MessageUpdate // Channel for task updates to be broadcasted

	// Statistics
	MessagesProcessed  atomic.Int64 // Total messages that have been processed (enqueued + dequeued)
	MessagesEnqueued   atomic.Int64 // Total messages enqueued
	MessagesDequeued   atomic.Int64 // Total messages dequeued
	LastStatsResetTime time.Time    // Time when stats were last reset
}

// ResetStats resets all the statistics counters
func (d *Diexodos) ResetStats() {
	d.MessagesProcessed.Store(0)
	d.MessagesEnqueued.Store(0)
	d.MessagesDequeued.Store(0)
	d.LastStatsResetTime = time.Now()
}

func (d *Diexodos) GetCurrentStats() map[string]interface{} {
	queueLength := len(d.MessageQueue)
	return map[string]interface{}{
		"name":               d.Name,
		"queueLength":        queueLength,
		"messagesProcessed":  d.MessagesProcessed.Load(),
		"messagesEnqueued":   d.MessagesEnqueued.Load(),
		"messagesDequeued":   d.MessagesDequeued.Load(),
		"lastMessageTime":    d.LastMessageReceived,
		"lastStatsResetTime": d.LastStatsResetTime,
		"uptime":             time.Since(d.LastStatsResetTime).String(),
	}
}
