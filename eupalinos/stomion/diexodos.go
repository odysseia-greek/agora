package stomion

import (
	"sync/atomic"
	"time"

	pb "github.com/odysseia-greek/agora/eupalinos/v1"
)

// Diexodos represents a task queue
type Diexodos struct {
	LastMessageReceived time.Time
	Name                string
	InternalID          string
	MessageQueue        map[string]*pb.InternalEpistello
	MessageOrder        []string
	PendingAcks         map[string]struct{}
	MessageUpdateCh     chan *pb.MessageUpdate // Channel for updates broadcast to replicas

	// Statistics
	MessagesProcessed  atomic.Int64 // Total messages that have been processed (enqueued + dequeued)
	MessagesEnqueued   atomic.Int64 // Total messages enqueued
	MessagesDequeued   atomic.Int64 // Total messages dequeued
	MessagesNacked     atomic.Int64 // Total messages returned for another attempt
	LastStatsResetTime time.Time    // Time when stats were last reset
}

// ResetStats resets all the statistics counters
func (d *Diexodos) ResetStats() {
	d.MessagesProcessed.Store(0)
	d.MessagesEnqueued.Store(0)
	d.MessagesDequeued.Store(0)
	d.MessagesNacked.Store(0)
	d.LastStatsResetTime = time.Now()
}

func (d *Diexodos) GetCurrentStats() map[string]interface{} {
	queueLength := len(d.MessageQueue)
	return map[string]interface{}{
		"name":               d.Name,
		"queueLength":        queueLength,
		"pendingAcks":        len(d.PendingAcks),
		"messagesProcessed":  d.MessagesProcessed.Load(),
		"messagesEnqueued":   d.MessagesEnqueued.Load(),
		"messagesDequeued":   d.MessagesDequeued.Load(),
		"messagesNacked":     d.MessagesNacked.Load(),
		"lastMessageTime":    d.LastMessageReceived,
		"lastStatsResetTime": d.LastStatsResetTime,
		"uptime":             time.Since(d.LastStatsResetTime).String(),
	}
}
