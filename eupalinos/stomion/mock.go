package stomion

import (
	"context"
	"fmt"

	pb "github.com/odysseia-greek/agora/eupalinos/v1"
)

// MockQueueService is an in-memory Eupalinos implementation intended for other
// services' tests. It runs the same queue code as the server without networking.
type MockQueueService struct {
	service *QueueServiceImpl
}

func NewMockQueueService() *MockQueueService {
	return &MockQueueService{service: &QueueServiceImpl{DiexodosMap: []*Diexodos{}}}
}

func (m *MockQueueService) WaitForHealthyState() bool { return true }

func (m *MockQueueService) StreamQueueUpdates(context.Context, *pb.MessageUpdate) (pb.Eupalinos_StreamQueueUpdatesClient, error) {
	return nil, fmt.Errorf("streaming is not supported by MockQueueService")
}

func (m *MockQueueService) EnqueueMessage(ctx context.Context, message *pb.Epistello) (*pb.EnqueueResponse, error) {
	return m.service.EnqueueMessage(ctx, message)
}

func (m *MockQueueService) EnqueueMessageBytes(ctx context.Context, message *pb.EpistelloBytes) (*pb.EnqueueResponse, error) {
	return m.service.EnqueueMessageBytes(ctx, message)
}

func (m *MockQueueService) DequeueMessage(ctx context.Context, channel *pb.ChannelInfo) (*pb.Epistello, error) {
	return m.service.DequeueMessage(ctx, channel)
}

func (m *MockQueueService) DequeueMessageBytes(ctx context.Context, channel *pb.ChannelInfo) (*pb.EpistelloBytes, error) {
	return m.service.DequeueMessageBytes(ctx, channel)
}

func (m *MockQueueService) AcknowledgeMessage(ctx context.Context, request *pb.AcknowledgeRequest) (*pb.AcknowledgeResponse, error) {
	return m.service.AcknowledgeMessage(ctx, request)
}

func (m *MockQueueService) NackMessage(ctx context.Context, request *pb.NackRequest) (*pb.NackResponse, error) {
	return m.service.NackMessage(ctx, request)
}

func (m *MockQueueService) GetQueueLength(ctx context.Context, channel *pb.ChannelInfo) (*pb.QueueLength, error) {
	return m.service.GetQueueLength(ctx, channel)
}

var _ QueueService = (*MockQueueService)(nil)
