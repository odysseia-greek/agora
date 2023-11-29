package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/eupalinos/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	pb "github.com/odysseia-greek/agora/eupalinos/proto"
)

// Diexodos represents a task queue
type Diexodos struct {
	LastMessageReceived time.Time
	Name                string
	InternalID          string
	MessageQueue        map[string]pb.InternalEpistello
	MessageUpdateCh     chan pb.MessageUpdate // Channel for task updates to be broadcasted
}

// EupalinosHandler is the gRPC server handling task queue operations
type EupalinosHandler struct {
	DiexodosMap []*Diexodos // Slice of Diexodos representing different queues
	Config      *config.Config
	pb.UnimplementedEupalinosServer
	mu sync.Mutex // Mutex to protect the task queue
}

const saveInterval = 5 * time.Minute // Set the interval for saving the state to disk

// StartBroadcasting starts a goroutine to handle broadcasting of task updates to replicas
func (s *EupalinosHandler) StartBroadcasting() {
	go func() {
		for {
			for _, channel := range s.DiexodosMap {
				select {
				case updatedTask := <-channel.MessageUpdateCh:
					// Broadcast the task update to all replicas (excluding the sender)
					for _, address := range s.Config.Addresses {
						conn, err := grpc.Dial(address, grpc.WithInsecure())
						if err != nil {
							log.Printf("error connecting to replica %s: %v", address, err)
							continue
						}
						defer conn.Close()

						client := pb.NewEupalinosClient(conn)
						stream, err := client.StreamQueueUpdates(context.Background())
						if err != nil {
							log.Printf("error creating stream to replica %s: %v", address, err)
							continue
						}
						if err := stream.Send(&updatedTask); err != nil {
							log.Printf("error sending task update to replica %s: %v", address, err)
						}
					}
				}
			}
		}
	}()
}

// EnqueueMessage handles message enqueueing
func (s *EupalinosHandler) EnqueueMessage(ctx context.Context, message *pb.Epistello) (*pb.EnqueueResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to get metadata from context")
	}

	var traceID string
	headerValue := md.Get(config.TRACING_KEY)
	if len(headerValue) > 0 {
		traceID = headerValue[0]
	}

	log.Printf("recieved enqueue: %v with traceId: %s", message, traceID)
	// Find or create the Diexodos with the matching channel name
	diexodos := s.findOrCreateDiexodos(message.Channel)

	// Generate a unique ID for the message
	messageID := uuid.New().String()

	// Set the ID field in the message
	internalMessage := &pb.InternalEpistello{
		Id:      messageID,
		Channel: message.Channel,
		Data:    message.Data,
		Traceid: traceID,
	}

	// Process the received message (e.g., enqueue)
	diexodos.MessageQueue[internalMessage.Id] = *internalMessage
	diexodos.LastMessageReceived = time.Now() // Update LastMessageReceived

	// Add the task update to the channel
	update := pb.MessageUpdate{
		Operation: pb.Operation_ENQUEUE,
		Message:   internalMessage,
	}

	if s.Config.Streaming {
		diexodos.MessageUpdateCh <- update
	}

	// Return the generated ID in the response
	response := &pb.EnqueueResponse{Id: messageID}

	return response, nil
}

// DequeueMessage handles message dequeueing from the specified channel
func (s *EupalinosHandler) DequeueMessage(ctx context.Context, channelInfo *pb.ChannelInfo) (*pb.Epistello, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find the Diexodos with the specified channel name
	diexodos := s.findDiexodosByChannel(channelInfo.Name)

	if diexodos == nil {
		return nil, fmt.Errorf("channel not found")
	}

	// Check if the message queue is empty
	if len(diexodos.MessageQueue) == 0 {
		return nil, fmt.Errorf("message queue is empty")
	}

	// Dequeue the first message from the message queue
	var internalMessage *pb.InternalEpistello

	for _, msg := range diexodos.MessageQueue {
		internalMessage = &msg
		break
	}

	// Remove the dequeued message from the message queue
	delete(diexodos.MessageQueue, internalMessage.Id)

	// Add the task update to the channel
	update := pb.MessageUpdate{
		Operation: pb.Operation_DEQUEUE,
		Message:   internalMessage,
	}

	if s.Config.Streaming {
		diexodos.MessageUpdateCh <- update
	}
	message := &pb.Epistello{
		Id:      internalMessage.Id,
		Channel: internalMessage.Channel,
		Data:    internalMessage.Data,
	}

	// Embed the trace ID into the metadata of the response
	responseMd := metadata.New(map[string]string{config.TRACING_KEY: internalMessage.Traceid})
	grpc.SendHeader(ctx, responseMd)
	return message, nil
}

// StreamQueueUpdates handles bidirectional streaming for task updates between Eupalinos pods
func (s *EupalinosHandler) StreamQueueUpdates(stream pb.Eupalinos_StreamQueueUpdatesServer) error {
	// Receive task update requests from other replicas
	for {
		updatedMessage, err := stream.Recv()
		if err != nil {
			return err
		}

		log.Printf("recieved message: %v", updatedMessage.Message)

		// Process the received task update
		// Update the task queue based on the task operation received from other replicas
		s.mu.Lock()
		diexodos := s.findOrCreateDiexodos(updatedMessage.Message.Channel)
		if diexodos == nil {
			s.mu.Unlock()
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
		for _, replica := range s.DiexodosMap {
			if replica.Name != diexodos.Name {
				replica.MessageUpdateCh <- *updatedMessage
			}
		}
		s.mu.Unlock()
	}
}

// GetQueueLength returns the length of the queue for the specified channel
func (s *EupalinosHandler) GetQueueLength(ctx context.Context, channelInfo *pb.ChannelInfo) (*pb.QueueLength, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find the Diexodos with the specified channel name
	diexodos := s.findDiexodosByChannel(channelInfo.Name)

	if diexodos == nil {
		return nil, fmt.Errorf("channel not found")
	}

	// Get the length of the message queue
	length := int32(len(diexodos.MessageQueue))

	log.Printf("length of queue: %d", length)

	return &pb.QueueLength{Length: length}, nil
}

// findOrCreateDiexodos finds the Diexodos with the matching channel name or creates a new one
func (s *EupalinosHandler) findOrCreateDiexodos(channelName string) *Diexodos {
	for _, d := range s.DiexodosMap {
		if d.Name == channelName {
			return d
		}
	}

	// If Diexodos with the given channel name does not exist, create a new one
	d := &Diexodos{
		Name:                channelName,
		InternalID:          uuid.New().String(),
		MessageQueue:        make(map[string]pb.InternalEpistello),
		MessageUpdateCh:     make(chan pb.MessageUpdate),
		LastMessageReceived: time.Now(),
	}
	s.DiexodosMap = append(s.DiexodosMap, d)

	log.Printf("created channel: %s", channelName)

	return d
}

// findDiexodosByChannel finds the Diexodos with the specified channel name
func (s *EupalinosHandler) findDiexodosByChannel(channelName string) *Diexodos {
	for _, d := range s.DiexodosMap {
		if d.Name == channelName {
			return d
		}
	}
	return nil
}

// StartAutoSave starts a goroutine to periodically save the state of the message queues to disk
func (s *EupalinosHandler) StartAutoSave() {
	go func() {
		for {
			s.SaveStateToDisk()
			time.Sleep(saveInterval)
		}
	}()
}

// SaveStateToDisk saves the state of the message queues to disk
func (s *EupalinosHandler) SaveStateToDisk() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create a map to store the serialized message queues
	state := make(map[string]map[string]pb.InternalEpistello)

	// Convert the message queues to a map
	for _, diexodos := range s.DiexodosMap {
		state[diexodos.Name] = diexodos.MessageQueue
	}

	// Serialize the state to JSON
	data, err := json.Marshal(state)
	if err != nil {
		log.Printf("error serializing state: %v", err)
		return
	}

	// Create the directory if it doesn't exist
	dir := filepath.Dir(s.Config.SavePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Printf("error creating directory: %v", err)
			return
		}
	}

	// Write the data to the file
	err = ioutil.WriteFile(s.Config.SavePath, data, 0644)
	if err != nil {
		log.Printf("error writing state to disk: %v", err)
	}
}

// LoadStateFromDisk loads the state of the message queues from disk on startup
func (s *EupalinosHandler) LoadStateFromDisk() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the file exists
	if _, err := os.Stat(s.Config.SavePath); os.IsNotExist(err) {
		// File does not exist, initialize the message queues with empty maps
		for _, diexodos := range s.DiexodosMap {
			diexodos.MessageQueue = make(map[string]pb.InternalEpistello)
		}
		return
	}

	// Read the data from the file
	data, err := ioutil.ReadFile(s.Config.SavePath)
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
	for _, diexodos := range s.DiexodosMap {
		queue, ok := state[diexodos.Name]
		if ok {
			diexodos.MessageQueue = queue
		} else {
			diexodos.MessageQueue = make(map[string]pb.InternalEpistello)
		}
	}
}
