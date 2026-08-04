package stomion

import (
	"context"
	"fmt"
	"time"

	"github.com/odysseia-greek/agora/plato/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/google/uuid"
	pb "github.com/odysseia-greek/agora/eupalinos/v1"
)

func (q *QueueServiceImpl) Health(context.Context, *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{
		Healthy: true,
		Time:    time.Now().String(),
		Version: q.Version,
	}, nil
}

// EnqueueMessage handles message enqueueing
func (q *QueueServiceImpl) EnqueueMessage(ctx context.Context, message *pb.Epistello) (*pb.EnqueueResponse, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	diexodos := q.findOrCreateDiexodos(message.Channel)

	// Generate a unique ID for the message
	messageID := uuid.New().String()

	// Set the ID field in the message
	internalMessage := &pb.InternalEpistello{
		Id:      messageID,
		Channel: message.Channel,
		Payload: &pb.InternalEpistello_Data{
			Data: message.Data,
		},
	}

	// Process the received message (e.g., enqueue)
	diexodos.MessageQueue[internalMessage.Id] = internalMessage
	diexodos.MessageOrder = append(diexodos.MessageOrder, internalMessage.Id)
	diexodos.LastMessageReceived = time.Now() // Update LastMessageReceived

	// Update statistics
	diexodos.MessagesProcessed.Add(1)
	diexodos.MessagesEnqueued.Add(1)

	// Add the task update to the channel
	update := &pb.MessageUpdate{
		Operation: pb.Operation_ENQUEUE,
		Message:   internalMessage,
	}

	if q.Streaming {
		select {
		case diexodos.MessageUpdateCh <- update:
		default:
			logging.Error("MessageUpdateCh full; dropping update")
		}
	}

	// Return the generated ID in the response
	response := &pb.EnqueueResponse{Id: messageID}

	return response, nil
}

// EnqueueMessageBytes handles message enqueueing
func (q *QueueServiceImpl) EnqueueMessageBytes(ctx context.Context, message *pb.EpistelloBytes) (*pb.EnqueueResponse, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	diexodos := q.findOrCreateDiexodos(message.Channel)

	// Generate a unique ID for the message
	messageID := uuid.New().String()

	// Set the ID field in the message
	internalMessage := &pb.InternalEpistello{
		Id: messageID,
		Payload: &pb.InternalEpistello_BytesData{
			BytesData: message.Data,
		},
		Channel: message.Channel,
		Traceid: "",
	}

	// Process the received message (e.g., enqueue)
	diexodos.MessageQueue[internalMessage.Id] = internalMessage
	diexodos.MessageOrder = append(diexodos.MessageOrder, internalMessage.Id)
	diexodos.LastMessageReceived = time.Now() // Update LastMessageReceived

	// Update statistics
	diexodos.MessagesProcessed.Add(1)
	diexodos.MessagesEnqueued.Add(1)

	// Add the task update to the channel
	update := &pb.MessageUpdate{
		Operation: pb.Operation_ENQUEUE,
		Message:   internalMessage,
	}

	if q.Streaming {
		select {
		case diexodos.MessageUpdateCh <- update:
		default:
			logging.Error("MessageUpdateCh full; dropping update")
		}
	}

	// Return the generated ID in the response
	response := &pb.EnqueueResponse{Id: messageID}

	return response, nil
}

func (q *QueueServiceImpl) DequeueMessage(ctx context.Context, channelInfo *pb.ChannelInfo) (*pb.Epistello, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	diexodos := q.findDiexodosByChannel(channelInfo.Name)
	if diexodos == nil {
		return nil, fmt.Errorf("channel not found")
	}
	if len(diexodos.MessageQueue) == 0 {
		return nil, fmt.Errorf("message queue is empty")
	}

	internalMessage, err := reserveMessage(diexodos, channelInfo.AckMode, false)
	if err != nil {
		return nil, err
	}

	diexodos.MessagesProcessed.Add(1)
	if !channelInfo.AckMode {
		q.recordRemoval(diexodos, internalMessage, pb.Operation_DEQUEUE)
	}

	// must be string payload
	p, ok := internalMessage.Payload.(*pb.InternalEpistello_Data)
	if !ok {
		return nil, fmt.Errorf("message %s is not a string payload (use DequeueBytes)", internalMessage.Id)
	}

	return &pb.Epistello{
		Id:       internalMessage.Id,
		Channel:  internalMessage.Channel,
		Data:     p.Data,
		Delivery: &pb.DeliveryMetadata{NackCount: internalMessage.NackCount},
	}, nil
}

// DequeueMessageBytes handles message dequeueing from the specified channel
func (q *QueueServiceImpl) DequeueMessageBytes(ctx context.Context, channelInfo *pb.ChannelInfo) (*pb.EpistelloBytes, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	diexodos := q.findDiexodosByChannel(channelInfo.Name)
	if diexodos == nil {
		return nil, fmt.Errorf("channel not found")
	}
	if len(diexodos.MessageQueue) == 0 {
		return nil, fmt.Errorf("message queue is empty")
	}

	internalMessage, err := reserveMessage(diexodos, channelInfo.AckMode, true)
	if err != nil {
		return nil, err
	}

	diexodos.MessagesProcessed.Add(1)
	if !channelInfo.AckMode {
		q.recordRemoval(diexodos, internalMessage, pb.Operation_DEQUEUE)
	}

	// must be string payload
	p, ok := internalMessage.Payload.(*pb.InternalEpistello_BytesData)
	if !ok {
		return nil, fmt.Errorf("message %s is not a string payload (use DequeueBytes)", internalMessage.Id)
	}

	return &pb.EpistelloBytes{
		Id:       internalMessage.Id,
		Channel:  internalMessage.Channel,
		Data:     p.BytesData,
		Delivery: &pb.DeliveryMetadata{NackCount: internalMessage.NackCount},
	}, nil
}

// AcknowledgeMessage permanently removes a message reserved by an ack-mode dequeue.
func (q *QueueServiceImpl) AcknowledgeMessage(_ context.Context, request *pb.AcknowledgeRequest) (*pb.AcknowledgeResponse, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	diexodos := q.findDiexodosByChannel(request.Channel)
	if diexodos == nil {
		return nil, fmt.Errorf("channel not found")
	}
	if _, pending := diexodos.PendingAcks[request.Id]; !pending {
		return nil, fmt.Errorf("message %s is not pending acknowledgement", request.Id)
	}
	internalMessage, exists := diexodos.MessageQueue[request.Id]
	if !exists {
		return nil, fmt.Errorf("message %s not found", request.Id)
	}

	q.recordRemoval(diexodos, internalMessage, pb.Operation_ACKNOWLEDGE)
	return &pb.AcknowledgeResponse{Acknowledged: true}, nil
}

// NackMessage releases a reserved message and moves it to the end of its queue.
func (q *QueueServiceImpl) NackMessage(_ context.Context, request *pb.NackRequest) (*pb.NackResponse, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	diexodos := q.findDiexodosByChannel(request.Channel)
	if diexodos == nil {
		return nil, fmt.Errorf("channel not found")
	}
	if _, pending := diexodos.PendingAcks[request.Id]; !pending {
		return nil, fmt.Errorf("message %s is not pending acknowledgement", request.Id)
	}
	message, exists := diexodos.MessageQueue[request.Id]
	if !exists {
		return nil, fmt.Errorf("message %s not found", request.Id)
	}

	delete(diexodos.PendingAcks, request.Id)
	message.NackCount++
	moveToBack(diexodos, request.Id)
	diexodos.MessagesProcessed.Add(1)
	diexodos.MessagesNacked.Add(1)
	q.broadcastUpdate(diexodos, message, pb.Operation_NACK)
	return &pb.NackResponse{Requeued: true, NackCount: message.NackCount}, nil
}

func reserveMessage(diexodos *Diexodos, ackMode, bytesPayload bool) (*pb.InternalEpistello, error) {
	if diexodos.PendingAcks == nil {
		diexodos.PendingAcks = make(map[string]struct{})
	}
	ensureMessageOrder(diexodos)
	for _, id := range diexodos.MessageOrder {
		message, exists := diexodos.MessageQueue[id]
		if !exists {
			continue
		}
		if _, pending := diexodos.PendingAcks[id]; pending {
			continue
		}
		_, isBytes := message.Payload.(*pb.InternalEpistello_BytesData)
		if isBytes != bytesPayload {
			continue
		}
		if ackMode {
			diexodos.PendingAcks[id] = struct{}{}
		}
		return message, nil
	}
	return nil, fmt.Errorf("no available message with the requested payload type")
}

func ensureMessageOrder(diexodos *Diexodos) {
	ordered := make(map[string]struct{}, len(diexodos.MessageOrder))
	for _, id := range diexodos.MessageOrder {
		ordered[id] = struct{}{}
	}
	for id := range diexodos.MessageQueue {
		if _, exists := ordered[id]; !exists {
			diexodos.MessageOrder = append(diexodos.MessageOrder, id)
		}
	}
}

func removeFromOrder(diexodos *Diexodos, id string) {
	for index, queuedID := range diexodos.MessageOrder {
		if queuedID == id {
			diexodos.MessageOrder = append(diexodos.MessageOrder[:index], diexodos.MessageOrder[index+1:]...)
			return
		}
	}
}

func moveToBack(diexodos *Diexodos, id string) {
	removeFromOrder(diexodos, id)
	diexodos.MessageOrder = append(diexodos.MessageOrder, id)
}

func (q *QueueServiceImpl) recordRemoval(diexodos *Diexodos, message *pb.InternalEpistello, operation pb.Operation) {
	delete(diexodos.PendingAcks, message.Id)
	delete(diexodos.MessageQueue, message.Id)
	removeFromOrder(diexodos, message.Id)
	diexodos.MessagesDequeued.Add(1)
	q.broadcastUpdate(diexodos, message, operation)
}

func (q *QueueServiceImpl) broadcastUpdate(diexodos *Diexodos, message *pb.InternalEpistello, operation pb.Operation) {
	if !q.Streaming {
		return
	}
	select {
	case diexodos.MessageUpdateCh <- &pb.MessageUpdate{Operation: operation, Message: message}:
	default:
		logging.Error("MessageUpdateCh full; dropping update")
	}
}

// StreamQueueUpdates handles bidirectional streaming for task updates between Eupalinos pods
func (q *QueueServiceImpl) StreamQueueUpdates(stream pb.Eupalinos_StreamQueueUpdatesServer) error {
	// Receive task update requests from other replicas
	for {
		updatedMessage, err := stream.Recv()
		if err != nil {
			return err
		}

		q.mu.Lock()
		diexodos := q.findOrCreateDiexodos(updatedMessage.Message.Channel)
		if updatedMessage.Operation == pb.Operation_ENQUEUE {
			diexodos.MessageQueue[updatedMessage.Message.Id] = updatedMessage.Message
			moveToBack(diexodos, updatedMessage.Message.Id)
		} else if updatedMessage.Operation == pb.Operation_DEQUEUE || updatedMessage.Operation == pb.Operation_ACKNOWLEDGE {
			delete(diexodos.MessageQueue, updatedMessage.Message.Id)
			delete(diexodos.PendingAcks, updatedMessage.Message.Id)
			removeFromOrder(diexodos, updatedMessage.Message.Id)
		} else if updatedMessage.Operation == pb.Operation_NACK {
			diexodos.MessageQueue[updatedMessage.Message.Id] = updatedMessage.Message
			delete(diexodos.PendingAcks, updatedMessage.Message.Id)
			moveToBack(diexodos, updatedMessage.Message.Id)
		}
		diexodos.LastMessageReceived = time.Now()
		q.mu.Unlock()

		// Broadcast the task update to all replicas (excluding the sender)
		for _, replica := range q.DiexodosMap {
			if replica.Name != diexodos.Name {
				replica.MessageUpdateCh <- updatedMessage
			}
		}
	}
}

// GetQueueLength returns the length of the queue for the specified channel
func (q *QueueServiceImpl) GetQueueLength(ctx context.Context, channelInfo *pb.ChannelInfo) (*pb.QueueLength, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Find the Diexodos with the specified channel name
	diexodos := q.findDiexodosByChannel(channelInfo.Name)

	if diexodos == nil {
		return nil, fmt.Errorf("channel not found")
	}

	// Get the length of the message queue
	length := int32(len(diexodos.MessageQueue))

	logging.Info(fmt.Sprintf("length of queue: %d", length))

	return &pb.QueueLength{Length: length, PendingAcknowledgements: int32(len(diexodos.PendingAcks))}, nil
}

// findOrCreateDiexodos finds the Diexodos with the matching channel name or creates a new one
func (q *QueueServiceImpl) findOrCreateDiexodos(channelName string) *Diexodos {
	for _, d := range q.DiexodosMap {
		if d.Name == channelName {
			return d
		}
	}

	// If Diexodos with the given channel name does not exist, create a new one
	now := time.Now()
	d := &Diexodos{
		Name:                channelName,
		InternalID:          uuid.New().String(),
		MessageQueue:        make(map[string]*pb.InternalEpistello),
		MessageOrder:        make([]string, 0),
		PendingAcks:         make(map[string]struct{}),
		MessageUpdateCh:     make(chan *pb.MessageUpdate, 1024),
		LastMessageReceived: now,
		LastStatsResetTime:  now,
	}

	// Initialize atomic counters
	d.MessagesProcessed.Store(0)
	d.MessagesEnqueued.Store(0)
	d.MessagesDequeued.Store(0)
	d.MessagesNacked.Store(0)

	q.DiexodosMap = append(q.DiexodosMap, d)

	logging.Info(fmt.Sprintf("created channel: %s", channelName))

	return d
}

// findDiexodosByChannel finds the Diexodos with the specified channel name
func (q *QueueServiceImpl) findDiexodosByChannel(channelName string) *Diexodos {
	for _, d := range q.DiexodosMap {
		if d.Name == channelName {
			return d
		}
	}
	return nil
}

func (q *QueueServiceImpl) StartBroadcasting() {
	go func() {
		for {
			for _, channel := range q.DiexodosMap {
				select {
				case updatedTask := <-channel.MessageUpdateCh:
					// Broadcast the task update to all replicas (excluding the sender)
					for _, address := range q.Addresses {
						// Use grpc.NewClient instead of the deprecated grpc.Dial
						conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
						if err != nil {
							logging.Error(fmt.Sprintf("error connecting to replica %s: %v", address, err))
							continue
						}

						// Create a client and stream
						client := pb.NewEupalinosClient(conn)
						stream, err := client.StreamQueueUpdates(context.Background())
						if err != nil {
							logging.Error(fmt.Sprintf("error creating stream to replica %s: %v", address, err))
							conn.Close() // Close connection on error
							continue
						}

						// Send the update
						if err := stream.Send(updatedTask); err != nil {
							logging.Error(fmt.Sprintf("error sending task update to replica %s: %v", address, err))
						}

						// Close the connection after sending
						conn.Close()
					}
				}
			}
		}
	}()
}

// PeriodStatsPrint periodically prints statistics about all channels
func (q *QueueServiceImpl) PeriodStatsPrint() {
	const statsInterval = 1 * time.Minute

	go func() {
		ticker := time.NewTicker(statsInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				q.mu.Lock()
				totalChannels := len(q.DiexodosMap)
				totalMessages := 0
				totalProcessed := int64(0)

				logging.Info("===== Eupalinos Queue Statistics =====")
				logging.Info(fmt.Sprintf("Total Channels: %d", totalChannels))

				// Individual channel stats
				for _, channel := range q.DiexodosMap {
					queueLen := len(channel.MessageQueue)
					totalMessages += queueLen
					processed := channel.MessagesProcessed.Load()
					totalProcessed += processed

					logging.Info(fmt.Sprintf("Channel: %s", channel.Name))
					logging.Info(fmt.Sprintf("  Queue Length: %d", queueLen))
					logging.Info(fmt.Sprintf("  Messages Processed: %d", processed))
					logging.Info(fmt.Sprintf("  Messages Enqueued: %d", channel.MessagesEnqueued.Load()))
					logging.Info(fmt.Sprintf("  Messages Dequeued: %d", channel.MessagesDequeued.Load()))
					logging.Info(fmt.Sprintf("  Messages Nacked: %d", channel.MessagesNacked.Load()))
					logging.Info(fmt.Sprintf("  Last Message Time: %s", channel.LastMessageReceived.Format(time.RFC3339)))
					logging.Info(fmt.Sprintf("  Channel Age: %s", time.Since(channel.LastStatsResetTime).Round(time.Second)))
				}

				// Summary stats
				logging.Info("===== Summary =====")
				logging.Info(fmt.Sprintf("Total Messages in Queue: %d", totalMessages))
				logging.Info(fmt.Sprintf("Total Messages Processed: %d", totalProcessed))
				logging.Info("=============================")

				q.mu.Unlock()
			}
		}
	}()
}
