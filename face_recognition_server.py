import socket
import threading
import json
import base64
import numpy as np
import cv2
import os
import logging
import urllib.request
from deepface import DeepFace

# --- Logging Configuration ---
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - [PYTHON] %(levelname)s - %(message)s',
    datefmt='%Y-%m-%d %H:%M:%S'
)
logger = logging.getLogger(__name__)

# --- Configuration ---
HOST = os.getenv('PYTHON_SERVER_HOST', '127.0.0.1')
PORT = int(os.getenv('PYTHON_SERVER_PORT', '5000'))

# --- Image Resizing Helper ---
def resize_image(image, max_width=640):
    (h, w) = image.shape[:2]
    if w > max_width:
        r = max_width / float(w)
        dim = (max_width, int(h * r))
        resized = cv2.resize(image, dim, interpolation=cv2.INTER_AREA)
        return resized
    return image

# --- Model Preloading ---
def preload_models():
    """Preload DeepFace models to reduce latency on first request."""
    logger.info("Preloading DeepFace models... This may take a moment.")
    try:
        # Perform a dummy verification to load models into memory
        # We can use small dummy black images
        dummy_img = np.zeros((200, 200, 3), dtype=np.uint8)
        DeepFace.build_model('VGG-Face') 
        # DeepFace.extract_faces(dummy_img, enforce_detection=False)
        logger.info("DeepFace models preloaded successfully.")
    except Exception as e:
        logger.error(f"Error preloading models: {e}")

# --- TCP Server Logic ---
def handle_client(conn, addr):
    logger.info(f"Connected by {addr}")
    buffer = ""
    try:
        while True:
            data = conn.recv(4096).decode('utf-8')
            if not data:
                break
            
            buffer += data
            
            while '\n' in buffer:
                message, buffer = buffer.split('\n', 1)
                
                try:
                    payload = json.loads(message)
                    action = payload.get("action", "compare_faces") # Default to compare_faces
                    client_image_data_b64 = payload.get("client_image_data")

                    result = {}
                    if not client_image_data_b64:
                        result = {"status": "error", "message": "No client image data provided."}
                    elif action == "check_face":
                        result = process_face_check_request(client_image_data_b64)
                    elif action == "compare_faces":
                        db_image_path = payload.get("db_image_path")
                        if not db_image_path:
                            result = {"status": "error", "message": "No database image path provided for comparison."}
                        else:
                            result = process_face_recognition_request(client_image_data_b64, db_image_path)
                    else:
                        result = {"status": "error", "message": f"Unknown action: {action}"}

                    response_json = json.dumps(result)
                    conn.sendall(response_json.encode('utf-8') + b'\n')

                except json.JSONDecodeError as e:
                    logger.error(f"JSON Decode Error: {e} - Message: {message}")
                    error_response = {"status": "error", "message": f"Invalid JSON: {e}"}
                    conn.sendall(json.dumps(error_response).encode('utf-8') + b'\n')
                except Exception as e:
                    logger.error(f"Error processing payload: {e}")
                    error_response = {"status": "error", "message": f"Processing error: {e}"}
                    conn.sendall(json.dumps(error_response).encode('utf-8') + b'\n')
            
    except ConnectionResetError:
        logger.warning(f"Client {addr} forcibly closed the connection.")
    except Exception as e:
        logger.error(f"Error handling client {addr}: {e}")
    finally:
        logger.info(f"Client {addr} disconnected.")
        conn.close()

# --- Face Check Logic ---
def process_face_check_request(client_image_b64):
    try:
        # Decode the image
        client_image_bytes = base64.b64decode(client_image_b64)
        np_arr_client = np.frombuffer(client_image_bytes, np.uint8)
        client_img = cv2.imdecode(np_arr_client, cv2.IMREAD_COLOR)

        if client_img is None:
            return {"status": "error", "message": "Could not decode client image."}

        # Resize image for faster processing
        client_img = resize_image(client_img)

        # Convert to RGB
        rgb_client_img = cv2.cvtColor(client_img, cv2.COLOR_BGR2RGB)

        # Get face representations. If no face, list is empty.
        embeddings = DeepFace.represent(rgb_client_img, enforce_detection=False)
        logger.info(f"Face check: Found {len(embeddings)} face(s).")

        if not embeddings:
            return {"status": "no_face_found", "message": "No face was found in the provided image."}
        
        # Optional: Check for multiple faces
        if len(embeddings) > 1:
            return {"status": "multiple_faces_found", "message": f"Multiple faces ({len(embeddings)}) were found. Please provide an image with only one face."}

        # Check for spoofing using DeepFace extract_faces
        try:
            logger.info("Starting spoof check...")
            face_objs = DeepFace.extract_faces(img=rgb_client_img, anti_spoofing=True)
            logger.info(f"Face objects: {len(face_objs)}")
            is_real = all(face_obj["is_real"] for face_obj in face_objs)
            logger.info(f"Spoof detection: is_real = {is_real}")
            if not is_real:
                return {"status": "spoof_detected", "message": "Spoofing detected in the image."}
        except Exception as e:
            logger.error(f"Spoof check error: {e}")
            # Continue if spoof check fails

        return {"status": "face_found", "message": "A single face was successfully found and no spoofing detected."}

    except Exception as e:
        logger.error(f"Error during face check processing: {e}")
        return {"status": "error", "message": f"Processing error: {str(e)}"}

# --- Face Recognition Logic with DeepFace library ---
def process_face_recognition_request(client_image_b64, db_image_path):
    try:
        # Decode the unknown image (from websocket)
        client_image_bytes = base64.b64decode(client_image_b64)
        np_arr_client = np.frombuffer(client_image_bytes, np.uint8)
        client_img = cv2.imdecode(np_arr_client, cv2.IMREAD_COLOR)

        if client_img is None:
            return {"status": "error", "message": "Could not decode client image."}

        # Resize client image
        client_img = resize_image(client_img)

        # Convert to RGB
        rgb_client_img = cv2.cvtColor(client_img, cv2.COLOR_BGR2RGB)

        # Load and process known image (from database path)
        try:
            if db_image_path.startswith("http://") or db_image_path.startswith("https://"):
                req = urllib.request.urlopen(db_image_path)
                arr = np.asarray(bytearray(req.read()), dtype=np.uint8)
                db_img = cv2.imdecode(arr, -1)
            else:
                db_img = cv2.imread(db_image_path)
                
            if db_img is None:
                return {"status": "error", "message": f"Could not load database image from {db_image_path}"}
        except Exception as e:
            return {"status": "error", "message": f"Could not load database image from {db_image_path}: {e}"}

        # Resize DB image
        db_img = resize_image(db_img)

        rgb_db_img = cv2.cvtColor(db_img, cv2.COLOR_BGR2RGB)

        # Verify faces with anti-spoofing enabled
        try:
            result = DeepFace.verify(rgb_client_img, rgb_db_img, anti_spoofing=True, threshold=0.5, enforce_detection=False)
            logger.info(f"Verification result: {result}")

            if result['verified']:
                return {"status": "recognized", "message": "Face recognized!"}
            else:
                return {"status": "unrecognized", "message": "Face not recognized."}

        except Exception as e:
            logger.error(f"Verification error: {e}")
            return {"status": "error", "message": f"Face verification failed: {str(e)}"}

    except Exception as e:
        logger.error(f"Error during face recognition processing: {e}")
        return {"status": "error", "message": f"Processing error: {str(e)}"}

def main():
    preload_models() # Preload models before starting server
    
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        s.bind((HOST, PORT))
        s.listen()
        logger.info(f"Python server listening on {HOST}:{PORT}")
        while True:
            conn, addr = s.accept()
            threading.Thread(target=handle_client, args=(conn, addr)).start()

if __name__ == "__main__":
    main()