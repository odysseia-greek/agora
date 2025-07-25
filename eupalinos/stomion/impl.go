package stomion

import (
	"context"
	"crypto/tls"
	"fmt"
	pb "github.com/odysseia-greek/agora/eupalinos/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
	"time"
)

type QueueService interface {
	WaitForHealthyState() bool
	StreamQueueUpdates(ctx context.Context, in *pb.ChannelInfo) (pb.Eupalinos_StreamQueueUpdatesClient, error)
	EnqueueMessage(ctx context.Context, in *pb.Epistello) (*pb.EnqueueResponse, error)
	DequeueMessage(ctx context.Context, in *pb.ChannelInfo) (*pb.Epistello, error)
	GetQueueLength(ctx context.Context, in *pb.ChannelInfo) (*pb.QueueLength, error)
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
		return nil, fmt.Errorf("failed to connect to tracing service: %w", err)
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
