package handlers

import (
	"encoding/json"
	"log"
	"net"
	"net/http"


	"go-face-auth/database/repository"


	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for now. In production, you should restrict this.
		return true
	},
}

// RequestPayload from WebSocket client
type RecognitionRequest struct {
	EmployeeID int    `json:"employee_id"` // Changed from UserID to EmployeeID
	ImageData  string `json:"image_data"`  // Base64 encoded image from client
}

// PythonRequestPayload to Python server
type PythonRecognitionRequest struct {
	ClientImageData string `json:"client_image_data"` // Base64 encoded image from client
	DBImagePath     string `json:"db_image_path"`     // Path to the image file on the Python server's side
}

// FaceRecognitionWebSocketHandler handles WebSocket connections for face recognition.
func FaceRecognitionWebSocketHandler(c *gin.Context) {
	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection to WebSocket: %v", err)
		return
	}
	defer wsConn.Close()

	log.Println("Face Recognition WebSocket client connected.")

	// Connect to Python TCP server
	pythonServerAddr := "127.0.0.1:5000" // Python server address
	pyConn, err := net.Dial("tcp", pythonServerAddr)
	if err != nil {
		log.Printf("Failed to connect to Python server at %s: %v", pythonServerAddr, err)
		wsConn.WriteMessage(websocket.TextMessage, []byte("Error: Could not connect to face recognition service."))
		return
	}
	defer pyConn.Close()

	log.Printf("Connected to Python server at %s.", pythonServerAddr)

	// Goroutine to read from Python server and send to WebSocket client
	go func() {
		decoder := json.NewDecoder(pyConn) // Use JSON decoder for Python responses
		for {
			var pythonResponse map[string]interface{}
			if err := decoder.Decode(&pythonResponse); err != nil {
				log.Printf("Error decoding JSON from Python server: %v", err)
				wsConn.Close()
				return
			}

			status, ok := pythonResponse["status"].(string)
			if !ok {
				log.Printf("Invalid 'status' in Python response: %v", pythonResponse)
				continue
			}

			var clientResponse gin.H
			switch status {
			case "recognized":
				clientResponse = gin.H{"status": "success", "message": "Wajah berhasil dikenali."}
			case "unrecognized":
				clientResponse = gin.H{"status": "failure", "message": "Wajah tidak cocok."}
			case "error":
				clientResponse = gin.H{"status": "error", "message": pythonResponse["message"]}
			default:
				clientResponse = gin.H{"status": "error", "message": "Status tidak diketahui dari server."}
			}

			responseBytes, _ := json.Marshal(clientResponse)
			if err := wsConn.WriteMessage(websocket.TextMessage, responseBytes); err != nil {
				log.Printf("Error writing to WebSocket client: %v", err)
				return
			}
		}
	}()

	// Read from WebSocket client and send to Python server
	for {
		messageType, p, err := wsConn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from WebSocket client: %v", err)
			break
		}

		if messageType == websocket.TextMessage {
			var req RecognitionRequest
			if err := json.Unmarshal(p, &req); err != nil {
				log.Printf("Error unmarshaling WebSocket request: %v", err)
				wsConn.WriteMessage(websocket.TextMessage, []byte("Error: Invalid JSON request."))
				continue
			}

			// Get the face image path for the given employee ID
			faceImages, err := repository.GetFaceImagesByEmployeeID(req.EmployeeID)
			if err != nil {
				log.Printf("Error getting face image from DB for employee %d: %v", req.EmployeeID, err)
				wsConn.WriteMessage(websocket.TextMessage, []byte("Error: Could not retrieve employee face image."))
				continue
			}

			if len(faceImages) == 0 {
				log.Printf("No face images found for employee %d in DB.", req.EmployeeID)
				wsConn.WriteMessage(websocket.TextMessage, []byte("Error: No registered face images for this employee."))
				continue
			}

			// For simplicity, use the first face image found
			dbImagePath := faceImages[0].ImagePath

			// Prepare payload for Python server
			pythonPayload := PythonRecognitionRequest{
				ClientImageData: req.ImageData, // Send raw base64 from client
				DBImagePath:     dbImagePath,   // Send path to DB image
			}
			payloadBytes, err := json.Marshal(pythonPayload)
			if err != nil {
				log.Printf("Error marshaling Python payload: %v", err)
				wsConn.WriteMessage(websocket.TextMessage, []byte("Error: Internal server error."))
				continue
			}

			// Send payload to Python server
			_, err = pyConn.Write(payloadBytes)
			if err != nil {
				log.Printf("Error writing to Python server: %v", err)
				break
			}
			// Add a newline to delimit JSON messages for Python
			_, err = pyConn.Write([]byte("\n"))
			if err != nil {
				log.Printf("Error writing newline to Python server: %v", err)
				break
			}

		} else {
			log.Printf("Received non-text message from WebSocket: %d", messageType)
			wsConn.WriteMessage(websocket.TextMessage, []byte("Error: Expected text JSON message."))
		}
	}
	log.Println("Face Recognition WebSocket client disconnected.")
}
