import socket
import threading
import json
import base64
import numpy as np
import cv2
import face_recognition # Changed from deepface

# --- Image Resizing Helper ---
def resize_image(image, max_width=640):
    (h, w) = image.shape[:2]
    if w > max_width:
        r = max_width / float(w)
        dim = (max_width, int(h * r))
        resized = cv2.resize(image, dim, interpolation=cv2.INTER_AREA)
        return resized
    return image

# --- Configuration ---
HOST = '127.0.0.1'
PORT = 5000

# --- TCP Server Logic (No changes here) ---
def handle_client(conn, addr):
    print(f"[PYTHON] Connected by {addr}")
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
                    client_image_data_b64 = payload.get("client_image_data")
                    db_image_path = payload.get("db_image_path")

                    if not client_image_data_b64:
                        recognition_result = {"status": "error", "message": "No client image data provided."}
                    elif not db_image_path:
                        recognition_result = {"status": "error", "message": "No database image path provided for comparison."}
                    else:
                        # This function is now updated to use face_recognition
                        recognition_result = process_face_recognition_request(client_image_data_b64, db_image_path)

                    response_json = json.dumps(recognition_result)
                    conn.sendall(response_json.encode('utf-8') + b'\n')

                except json.JSONDecodeError as e:
                    print(f"[PYTHON] JSON Decode Error: {e} - Message: {message}")
                    error_response = {"status": "error", "message": f"Invalid JSON: {e}"}
                    conn.sendall(json.dumps(error_response).encode('utf-8') + b'\n')
                except Exception as e:
                    print(f"[PYTHON] Error processing payload: {e}")
                    error_response = {"status": "error", "message": f"Processing error: {e}"}
                    conn.sendall(json.dumps(error_response).encode('utf-8') + b'\n')
            
    except ConnectionResetError:
        print(f"[PYTHON] Client {addr} forcibly closed the connection.")
    except Exception as e:
        print(f"[PYTHON] Error handling client {addr}: {e}")
    finally:
        print(f"[PYTHON] Client {addr} disconnected.")
        conn.close()

# --- Face Recognition Logic with face_recognition library ---
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

        # The face_recognition library uses RGB images, but OpenCV uses BGR. Convert it.
        rgb_client_img = cv2.cvtColor(client_img, cv2.COLOR_BGR2RGB)

        # Get face encoding for the unknown image. It returns a list of faces.
        unknown_encodings = face_recognition.face_encodings(rgb_client_img)
        print(f"[PYTHON] Client image: Found {len(unknown_encodings)} face(s).")
        if not unknown_encodings:
            return {"status": "error", "message": "No face found in the client image."}
        unknown_encoding = unknown_encodings[0] # Use the first face found

        # Process known image (from database path)
        try:
            db_img = face_recognition.load_image_file(db_image_path)
        except FileNotFoundError:
            return {"status": "error", "message": f"Database image not found at {db_image_path}"}
        except Exception as e:
            return {"status": "error", "message": f"Could not load database image from {db_image_path}: {e}"}

        if db_img is None:
            return {"status": "error", "message": "Could not decode database image."}

        # Resize DB image
        db_img = resize_image(db_img)

        rgb_db_img = cv2.cvtColor(db_img, cv2.COLOR_BGR2RGB)
        db_encodings = face_recognition.face_encodings(rgb_db_img)
        print(f"[PYTHON] DB image: Found {len(db_encodings)} face(s).")
        if not db_encodings:
            return {"status": "error", "message": "No face found in the database image."}
        known_encoding = db_encodings[0]

        # Compare the unknown face with the known face
        # Default tolerance is 0.6. Lower value means stricter match.
  
        tolerance = 0.5 # You can adjust this value. Lower value means stricter match.
        results = face_recognition.compare_faces([known_encoding], unknown_encoding, tolerance=tolerance)
        print(f"[PYTHON] Comparison results): {results}")

        if True in results:
            return {"status": "recognized", "message": "Face recognized!"}
        else:
            return {"status": "unrecognized", "message": "Face not recognized."}

    except Exception as e:
        print(f"[PYTHON] Error during face recognition processing: {e}")
        return {"status": "error", "message": f"Processing error: {str(e)}"}

def main():
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        s.bind((HOST, PORT))
        s.listen()
        print(f"[PYTHON] Python server listening on {HOST}:{PORT}")
        while True:
            conn, addr = s.accept()
            threading.Thread(target=handle_client, args=(conn, addr)).start()

if __name__ == "__main__":
    main()