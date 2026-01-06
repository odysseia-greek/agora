package stomion

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	pb "github.com/odysseia-greek/agora/eupalinos/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type QueueService interface {
	WaitForHealthyState() bool
	StreamQueueUpdates(ctx context.Context, in *pb.MessageUpdate) (pb.Eupalinos_StreamQueueUpdatesClient, error)
	EnqueueMessage(ctx context.Context, in *pb.Epistello) (*pb.EnqueueResponse, error)
	DequeueMessage(ctx context.Context, in *pb.ChannelInfo) (*pb.Epistello, error)
	GetQueueLength(ctx context.Context, in *pb.ChannelInfo) (*pb.QueueLength, error)
	EnqueueMessageBytes(ctx context.Context, in *pb.EpistelloBytes) (*pb.EnqueueResponse, error)
	DequeueMessageBytes(ctx context.Context, in *pb.ChannelInfo) (*pb.EpistelloBytes, error)
}

const (
	DEFAULTADDRESS string = "localhost:50060"
)

type QueueServiceImpl struct {
	Version     string
	DiexodosMap []*Diexodos
	mu          sync.Mutex
	pb.UnimplementedEupalinosServer
	Addresses []string // Addresses for each replica
	Streaming bool
	SavePath  string
	TLSConfig *tls.Config
}

type QueueServiceClient struct {
	Impl QueueService
}

type QueueClient struct {
	queue pb.EupalinosClient
}

func NewEupalinosClient(address string) (*QueueClient, error) {
	if address == "" {
		address = DEFAULTADDRESS
	}
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to queue service: %w", err)
	}
	client := pb.NewEupalinosClient(conn)
	return &QueueClient{queue: client}, nil
}

func (q *QueueClient) WaitForHealthyState() bool {
	timeout := 30 * time.Second
	checkInterval := 1 * time.Second
	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		response, err := q.Health(context.Background(), &pb.HealthRequest{})
		if err == nil && response.Healthy {
			return true
		}

		time.Sleep(checkInterval)
	}

	return false
}

func (q *QueueClient) Health(ctx context.Context, request *pb.HealthRequest) (*pb.HealthResponse, error) {
	return q.queue.Health(ctx, request)
}

func (q *QueueClient) StreamQueueUpdates(ctx context.Context, in *pb.MessageUpdate) (pb.Eupalinos_StreamQueueUpdatesClient, error) {
	return q.queue.StreamQueueUpdates(ctx)
}

func (q *QueueClient) EnqueueMessage(ctx context.Context, request *pb.Epistello) (*pb.EnqueueResponse, error) {
	return q.queue.EnqueueMessage(ctx, request)
}

func (q *QueueClient) EnqueueMessageBytes(ctx context.Context, request *pb.EpistelloBytes) (*pb.EnqueueResponse, error) {
	return q.queue.EnqueueMessageBytes(ctx, request)
}

func (q *QueueClient) DequeueMessage(ctx context.Context, request *pb.ChannelInfo) (*pb.Epistello, error) {
	return q.queue.DequeueMessage(ctx, request)
}

func (q *QueueClient) DequeueMessageBytes(ctx context.Context, request *pb.ChannelInfo) (*pb.EpistelloBytes, error) {
	return q.queue.DequeueMessageBytes(ctx, request)
}

func (q *QueueClient) GetQueueLength(ctx context.Context, request *pb.ChannelInfo) (*pb.QueueLength, error) {
	return q.queue.GetQueueLength(ctx, request)
}
