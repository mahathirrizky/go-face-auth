package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

// PythonRecognitionRequest to Python server
type PythonRecognitionRequest struct {
	ClientImageData string `json:"client_image_data"` // Base64 encoded image from client
	DBImagePath     string `json:"db_image_path"`     // Path to the image file on the Python server's side
	Action          string `json:"action,omitempty"`  // Action to perform: "compare_faces" (default) or "check_face"
}

// PythonServerClientInterface defines the interface for Python server communication.
type PythonServerClientInterface interface {
	SendToPythonServer(payload PythonRecognitionRequest) (map[string]interface{}, error)
	Close() error
}

// pythonClientImpl is a concrete implementation of PythonServerClientInterface.
type pythonClientImpl struct {
	addr string
	conn net.Conn
	mu   sync.Mutex
}

// NewPythonClient creates a new instance of PythonServerClientInterface.
func NewPythonClient() PythonServerClientInterface {
	addr := os.Getenv("PYTHON_SERVER_ADDRESS")
	if addr == "" {
		addr = "127.0.0.1:5000" // Default to localhost if not set
	}
	return &pythonClientImpl{
		addr: addr,
	}
}

// getConn returns an existing connection or creates a new one safely.
func (p *pythonClientImpl) getConn() (net.Conn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// If connection exists, return it
	// Note: It's hard to know if a TCP conn is closed without writing to it.
	// We will try to use it, and if it fails, we will reconnect in the SendToPythonServer method.
	if p.conn != nil {
		return p.conn, nil
	}

	// Dial new connection
	conn, err := net.DialTimeout("tcp", p.addr, 2*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Python server at %s: %w", p.addr, err)
	}

	p.conn = conn
	return p.conn, nil
}

// closeConn closes the connection and sets it to nil.
func (p *pythonClientImpl) closeConn() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.conn != nil {
		p.conn.Close()
		p.conn = nil
	}
}

// Close explicitly closes the connection.
func (p *pythonClientImpl) Close() error {
	p.closeConn()
	return nil
}

// SendToPythonServer connects to the Python TCP server, sends the payload, and returns the response.
func (p *pythonClientImpl) SendToPythonServer(payload PythonRecognitionRequest) (map[string]interface{}, error) {
	// Try up to 2 times (retry once on connection failure)
	const maxRetries = 2
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		conn, err := p.getConn()
		if err != nil {
			lastErr = err
			// If we can't even get a connection, maybe wait a bit or just retry?
			// For now, let's treat dial error as immediate failure if it persists.
			// But since getConn dials, maybe we should just return error?
			// Actually getConn returns error if Dial fails. So we might fail here.
			// Let's try to backoff slightly if it's the first try
			if i == 0 {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			return nil, lastErr
		}

		// Marshal payload to JSON
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}

		// Set a deadline for write and read
		// This prevents hanging if the server is stuck
		conn.SetDeadline(time.Now().Add(10 * time.Second))

		// Send payload to Python server with a newline delimiter
		_, err = conn.Write(append(payloadBytes, '\n'))
		if err != nil {
			// Write failed. Connection is likely bad.
			log.Printf("Failed to write to Python server: %v. Reconnecting...", err)
			p.closeConn() // Force close and nil
			lastErr = err
			continue // Retry with new connection
		}

		// Read response from Python server
		decoder := json.NewDecoder(conn)
		var pythonResponse map[string]interface{}
		if err := decoder.Decode(&pythonResponse); err != nil {
			// Read failed.
			log.Printf("Failed to decode response from Python server: %v. Reconnecting...", err)
			p.closeConn()
			lastErr = err
			continue
		}

		// Reset deadline
		conn.SetDeadline(time.Time{})

		return pythonResponse, nil
	}

	return nil, fmt.Errorf("failed to communicate with Python server after retries: %w", lastErr)
}
