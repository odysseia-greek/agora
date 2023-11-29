package app

import (
	"context"
	"encoding/json"
	"github.com/odysseia-greek/agora/eupalinos/config"
	pb "github.com/odysseia-greek/agora/eupalinos/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"
)

const GENERATED_ID string = "eupalinos-test-uuid"

// Create a mock Eupalinos server
func NewMockEupalinosServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle the EnqueueMessage request
		if r.Method == http.MethodPost && r.URL.Path == "/proto.Eupalinos/EnqueueMessage" {
			// Decode the request body into pb.Epistello message
			var requestMessage pb.Epistello
			err := json.NewDecoder(r.Body).Decode(&requestMessage)
			if err != nil {
				http.Error(w, "failed to decode request", http.StatusBadRequest)
				return
			}

			// Here, you can create a UUID or any other response data as needed
			// For example, create a new UUID.
			uuid := GENERATED_ID

			// Create the response message with the generated UUID
			responseMessage := &pb.EnqueueResponse{
				Id: uuid,
			}

			// Marshal the response message into bytes
			responseData, err := json.Marshal(responseMessage)
			if err != nil {
				http.Error(w, "failed to marshal response", http.StatusInternalServerError)
				return
			}

			// Set the Content-Type header to "application/x-protobuf"
			w.Header().Set("Content-Type", "application/x-protobuf")

			// Write the response data (protobuf encoded) to the response writer
			w.Write(responseData)
			return
		}

		// Handle other API endpoints or return 404 for unknown endpoints
		http.NotFound(w, r)
	}))
}

//
//func TestEnqueueMessageTwo(t *testing.T) {
//	// Set up the mock Eupalinos server
//	mockServer := NewMockEupalinosServer()
//	defer mockServer.Close()
//
//	// Create an HTTP client using the mock server URL
//	client := &http.Client{}
//
//	// Create the payload data
//	payload := pb.Epistello{
//		Data:    `{"method":"","category":"","greek":"Ἄβδηρα","translation":"town of Abdera, known for stupidity of inhabitants","chapter":57}`,
//		Channel: "test-channel",
//	}
//
//	// Convert the payload data to JSON
//	payloadJSON, err := json.Marshal(payload)
//	assert.Nil(t, err)
//
//	// Set the URL and method for the request
//	url := mockServer.URL + "/proto.Eupalinos/EnqueueMessage"
//	method := "POST"
//
//	// Create the HTTP request with the JSON payload
//	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadJSON))
//	assert.Nil(t, err)
//	// Send the HTTP request to the mock server
//	resp, err := client.Do(req)
//	assert.Nil(t, err)
//
//	defer resp.Body.Close()
//
//	// Read and decode the response (protobuf encoded) into pb.EnqueueResponse
//	var enqueueResponse pb.EnqueueResponse
//	err = json.NewDecoder(resp.Body).Decode(&enqueueResponse)
//	assert.Nil(t, err)
//
//	// Use assert to check the expected results based on the mock response
//	assert.Equal(t, GENERATED_ID, enqueueResponse.Id)
//}

func TestEnqueueMessage(t *testing.T) {
	// Create a new EupalinosHandler with the desired configuration
	handler := &EupalinosHandler{
		Config: &config.Config{
			Streaming: false, // Set this to true if you want to test streaming behavior
		},
		DiexodosMap: []*Diexodos{},
	}

	// Create a context with the desired metadata (traceID)
	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{config.TRACING_KEY: "your-trace-id"}))

	// Create the payload data
	payload := &pb.Epistello{
		Data:    `{"method":"","category":"","greek":"Ἄβδηρα","translation":"town of Abdera, known for stupidity of inhabitants","chapter":57}`,
		Channel: "test-channel",
	}

	// Call the EnqueueMessage function with the mock request and response writer
	response, err := handler.EnqueueMessage(ctx, payload)

	// Check for errors
	assert.Nil(t, err)

	// Use assert to check the expected results based on the mock response
	assert.True(t, IsUUID(response.Id))
}

func TestDequeueMessage(t *testing.T) {
	channelName := "test-channel"
	data := "{\"method\":\"\",\"category\":\"\",\"greek\":\"Ἄβδηρα\",\"translation\":\"town of Abdera, known for stupidity of inhabitants\",\"chapter\":57}"

	t.Run("ChannelWithOneMessage", func(t *testing.T) {
		// Create a new EupalinosHandler with the desired configuration
		handler := &EupalinosHandler{
			Config: &config.Config{
				Streaming: false, // Set this to true if you want to test streaming behavior
			},
			DiexodosMap: []*Diexodos{
				{
					LastMessageReceived: time.Now(),
					Name:                channelName,
					InternalID:          "",
					MessageQueue: map[string]pb.InternalEpistello{
						channelName: {
							Id:      "testuuid",
							Data:    data,
							Channel: channelName,
							Traceid: GENERATED_ID,
						},
					},
					MessageUpdateCh: nil,
				},
			},
		}

		// Create the payload data
		payload := &pb.ChannelInfo{
			Name: channelName,
		}

		// Call the EnqueueMessage function with the mock request and response writer
		response, err := handler.DequeueMessage(context.Background(), payload)

		// Check for errors
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, response.Channel, channelName)
		assert.Equal(t, response.Data, data)
	})

	t.Run("NoChannel", func(t *testing.T) {
		// Create a new EupalinosHandler with the desired configuration
		handler := &EupalinosHandler{
			Config: &config.Config{
				Streaming: false, // Set this to true if you want to test streaming behavior
			},
			DiexodosMap: []*Diexodos{},
		}

		// Create the payload data
		payload := &pb.ChannelInfo{
			Name: "test-queue",
		}

		// Call the EnqueueMessage function with the mock request and response writer
		response, err := handler.DequeueMessage(context.Background(), payload)

		// Check for errors
		assert.NotNil(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "channel not found")
	})

	t.Run("ChannelWithTwoMessagesLengthChanges", func(t *testing.T) {
		// Create a new EupalinosHandler with the desired configuration
		handler := &EupalinosHandler{
			Config: &config.Config{
				Streaming: false, // Set this to true if you want to test streaming behavior
			},
			DiexodosMap: []*Diexodos{
				{
					LastMessageReceived: time.Now(),
					Name:                channelName,
					InternalID:          "",
					MessageQueue: map[string]pb.InternalEpistello{
						"testuuid": {
							Id:      "testuuid",
							Data:    data,
							Channel: channelName,
							Traceid: GENERATED_ID,
						},
						"testuuid2": {
							Id:      "testuuid2",
							Data:    data,
							Channel: channelName,
							Traceid: GENERATED_ID,
						},
					},
					MessageUpdateCh: nil,
				},
			},
		}

		// Create the payload data
		payload := &pb.ChannelInfo{
			Name: channelName,
		}

		length, err := handler.GetQueueLength(context.Background(), payload)
		assert.Nil(t, err)
		assert.Equal(t, int32(2), length.Length)

		// Call the EnqueueMessage function with the mock request and response writer
		response, err := handler.DequeueMessage(context.Background(), payload)

		// Check for errors
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, channelName, response.Channel)
		assert.Equal(t, data, response.Data)

		lengthAfter, err := handler.GetQueueLength(context.Background(), payload)
		assert.Nil(t, err)
		assert.Equal(t, int32(1), lengthAfter.Length)
	})
}

func IsUUID(s string) bool {
	// Regular expression pattern for UUID (version 1, 2, 3, 4, and 5)
	uuidPattern := `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`
	re := regexp.MustCompile(uuidPattern)
	return re.MatchString(s)
}
