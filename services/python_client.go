package services

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

// PythonRecognitionRequest to Python server
type PythonRecognitionRequest struct {
	ClientImageData string `json:"client_image_data"` // Base64 encoded image from client
	DBImagePath     string `json:"db_image_path"`     // Path to the image file on the Python server's side
}

// PythonServerClientInterface defines the interface for Python server communication.
type PythonServerClientInterface interface {
	SendToPythonServer(payload PythonRecognitionRequest) (map[string]interface{}, error)
}

// pythonClientImpl is a concrete implementation of PythonServerClientInterface.
type pythonClientImpl struct{}

// NewPythonClient creates a new instance of PythonServerClientInterface.
func NewPythonClient() PythonServerClientInterface {
	return &pythonClientImpl{}
}

// SendToPythonServer connects to the Python TCP server, sends the payload, and returns the response.
func (p *pythonClientImpl) SendToPythonServer(payload PythonRecognitionRequest) (map[string]interface{}, error) {
	pythonServerAddr := os.Getenv("PYTHON_SERVER_ADDRESS")
	if pythonServerAddr == "" {
		pythonServerAddr = "127.0.0.1:5000" // Default to localhost if not set
	}
	conn, err := net.Dial("tcp", pythonServerAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Python server: %w", err)
	}
	defer conn.Close()

	// Marshal payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Send payload to Python server with a newline delimiter
	_, err = conn.Write(append(payloadBytes, '\n'))
	if err != nil {
		return nil, fmt.Errorf("failed to send payload to Python server: %w", err)
	}

	// Read response from Python server
	decoder := json.NewDecoder(conn)
	var pythonResponse map[string]interface{}
	if err := decoder.Decode(&pythonResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response from Python server: %w", err)
	}

	return pythonResponse, nil
}
