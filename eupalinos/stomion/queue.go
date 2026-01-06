package stomion

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/odysseia-greek/agora/plato/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/google/uuid"
	pb "github.com/odysseia-greek/agora/eupalinos/proto"
)

const saveInterval = 5 * time.Minute // Set the interval for saving the state to disk

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
	diexodos.MessageQueue[internalMessage.Id] = *internalMessage
	diexodos.LastMessageReceived = time.Now() // Update LastMessageReceived

	// Update statistics
	diexodos.MessagesProcessed.Add(1)
	diexodos.MessagesEnqueued.Add(1)

	// Add the task update to the channel
	update := pb.MessageUpdate{
		Operation: pb.Operation_ENQUEUE,
		Message:   internalMessage,
	}

	if q.Streaming {
		diexodos.MessageUpdateCh <- update
	}

	// Return the generated ID in the response
	response := &pb.EnqueueResponse{Id: messageID}

	return response, nil
}

// EnqueueMessage handles message enqueueing
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
	diexodos.MessageQueue[internalMessage.Id] = *internalMessage
	diexodos.LastMessageReceived = time.Now() // Update LastMessageReceived

	// Update statistics
	diexodos.MessagesProcessed.Add(1)
	diexodos.MessagesEnqueued.Add(1)

	// Add the task update to the channel
	update := pb.MessageUpdate{
		Operation: pb.Operation_ENQUEUE,
		Message:   internalMessage,
	}

	if q.Streaming {
		diexodos.MessageUpdateCh <- update
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

	// dequeue first
	var internalMessage *pb.InternalEpistello
	for _, msg := range diexodos.MessageQueue {
		internalMessage = &msg
		break
	}
	delete(diexodos.MessageQueue, internalMessage.Id)

	// stats + streaming...
	diexodos.MessagesProcessed.Add(1)
	diexodos.MessagesDequeued.Add(1)

	if q.Streaming {
		diexodos.MessageUpdateCh <- pb.MessageUpdate{
			Operation: pb.Operation_DEQUEUE,
			Message:   internalMessage,
		}
	}

	// must be string payload
	p, ok := internalMessage.Payload.(*pb.InternalEpistello_Data)
	if !ok {
		return nil, fmt.Errorf("message %s is not a string payload (use DequeueBytes)", internalMessage.Id)
	}

	return &pb.Epistello{
		Id:      internalMessage.Id,
		Channel: internalMessage.Channel,
		Data:    p.Data,
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

	// dequeue first
	var internalMessage *pb.InternalEpistello
	for _, msg := range diexodos.MessageQueue {
		internalMessage = &msg
		break
	}
	delete(diexodos.MessageQueue, internalMessage.Id)

	// stats + streaming...
	diexodos.MessagesProcessed.Add(1)
	diexodos.MessagesDequeued.Add(1)

	if q.Streaming {
		diexodos.MessageUpdateCh <- pb.MessageUpdate{
			Operation: pb.Operation_DEQUEUE,
			Message:   internalMessage,
		}
	}

	// must be string payload
	p, ok := internalMessage.Payload.(*pb.InternalEpistello_BytesData)
	if !ok {
		return nil, fmt.Errorf("message %s is not a string payload (use DequeueBytes)", internalMessage.Id)
	}

	return &pb.EpistelloBytes{
		Id:      internalMessage.Id,
		Channel: internalMessage.Channel,
		Data:    p.BytesData,
	}, nil
}

// StreamQueueUpdates handles bidirectional streaming for task updates between Eupalinos pods
func (q *QueueServiceImpl) StreamQueueUpdates(stream pb.Eupalinos_StreamQueueUpdatesServer) error {
	// Receive task update requests from other replicas
	for {
		updatedMessage, err := stream.Recv()
		if err != nil {
			return err
		}

		logging.Info(fmt.Sprintf("recieved message: %v", updatedMessage.Message))

		// Process the received task update
		// Update the task queue based on the task operation received from other replicas
		q.mu.Lock()
		defer q.mu.Unlock()
		diexodos := q.findOrCreateDiexodos(updatedMessage.Message.Channel)
		if diexodos == nil {
			return fmt.Errorf("channel not found: %s", updatedMessage.Message.Channel)
		}

		if updatedMessage.Operation == pb.Operation_ENQUEUE {
			diexodos.MessageQueue[updatedMessage.Message.Id] = *updatedMessage.Message
		} else if updatedMessage.Operation == pb.Operation_DEQUEUE {
			delete(diexodos.MessageQueue, updatedMessage.Message.Id)
		}

		// Update the last message received time for the Diexodos
		diexodos.LastMessageReceived = time.Now()

		// Broadcast the task update to all replicas (excluding the sender)
		for _, replica := range q.DiexodosMap {
			if replica.Name != diexodos.Name {
				replica.MessageUpdateCh <- *updatedMessage
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

	return &pb.QueueLength{Length: length}, nil
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
		MessageQueue:        make(map[string]pb.InternalEpistello),
		MessageUpdateCh:     make(chan pb.MessageUpdate),
		LastMessageReceived: now,
		LastStatsResetTime:  now,
	}

	// Initialize atomic counters
	d.MessagesProcessed.Store(0)
	d.MessagesEnqueued.Store(0)
	d.MessagesDequeued.Store(0)

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

// StartAutoSave starts a goroutine to periodically save the state of the message queues to disk
func (q *QueueServiceImpl) StartAutoSave() {
	go func() {
		for {
			q.SaveStateToDisk()
			time.Sleep(saveInterval)
		}
	}()
}

// SaveStateToDisk saves the state of the message queues to disk
func (q *QueueServiceImpl) SaveStateToDisk() {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Create a map to store the serialized message queues
	state := make(map[string]map[string]pb.InternalEpistello)

	// Convert the message queues to a map
	for _, diexodos := range q.DiexodosMap {
		state[diexodos.Name] = diexodos.MessageQueue
	}

	// Serialize the state to JSON
	data, err := json.Marshal(state)
	if err != nil {
		log.Printf("error serializing state: %v", err)
		return
	}

	// Create the directory if it doesn't exist
	dir := filepath.Dir(q.SavePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Printf("error creating directory: %v", err)
			return
		}
	}

	// Write the data to the file
	err = os.WriteFile(q.SavePath, data, 0644)
	if err != nil {
		log.Printf("error writing state to disk: %v", err)
	}
}

// LoadStateFromDisk loads the state of the message queues from disk on startup
func (q *QueueServiceImpl) LoadStateFromDisk() {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Check if the file exists
	if _, err := os.Stat(q.SavePath); os.IsNotExist(err) {
		// File does not exist, initialize the message queues with empty maps
		for _, diexodos := range q.DiexodosMap {
			diexodos.MessageQueue = make(map[string]pb.InternalEpistello)
		}
		return
	}

	// Read the data from the file
	data, err := os.ReadFile(q.SavePath)
	if err != nil {
		log.Printf("error reading state from disk: %v", err)
		return
	}

	// Deserialize the data
	state := make(map[string]map[string]pb.InternalEpistello)
	err = json.Unmarshal(data, &state)
	if err != nil {
		log.Printf("error deserializing state: %v", err)
		return
	}

	// Convert the map to message queues
	for _, diexodos := range q.DiexodosMap {
		queue, ok := state[diexodos.Name]
		if ok {
			diexodos.MessageQueue = queue
		} else {
			diexodos.MessageQueue = make(map[string]pb.InternalEpistello)
		}
	}
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
						if err := stream.Send(&updatedTask); err != nil {
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
