import socket
import json
import base64
import cv2
import sys

def test_face_check(image_path):
    # Load image
    img = cv2.imread(image_path)
    if img is None:
        print("Gambar tidak ditemukan!")
        return
    # Encode to base64
    _, buffer = cv2.imencode('.jpg', img)
    img_base64 = base64.b64encode(buffer).decode('utf-8')

    # Prepare payload
    payload = {
        "action": "check_face",
        "client_image_data": img_base64
    }

    # Send to server
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.connect(('127.0.0.1', 5000))
        data = json.dumps(payload) + '\n'
        s.sendall(data.encode('utf-8'))
        response = s.recv(4096).decode('utf-8')
        print("Response:", json.loads(response))

def test_face_compare(client_image_path, db_image_path):
    # Load client image
    img = cv2.imread(client_image_path)
    if img is None:
        print("Gambar client tidak ditemukan!")
        return
    _, buffer = cv2.imencode('.jpg', img)
    img_base64 = base64.b64encode(buffer).decode('utf-8')

    # Prepare payload
    payload = {
        "action": "compare_faces",
        "client_image_data": img_base64,
        "db_image_path": db_image_path
    }

    # Send to server
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.connect(('127.0.0.1', 5000))
        data = json.dumps(payload) + '\n'
        s.sendall(data.encode('utf-8'))
        response = s.recv(4096).decode('utf-8')
        print("Response:", json.loads(response))

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python test_client.py check <image_path>")
        print("       python test_client.py compare <client_image> <db_image>")
        sys.exit(1)

    action = sys.argv[1]
    if action == "check" and len(sys.argv) == 3:
        test_face_check(sys.argv[2])
    elif action == "compare" and len(sys.argv) == 4:
        test_face_compare(sys.argv[2], sys.argv[3])
    else:
        print("Invalid arguments")