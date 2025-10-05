#!/bin/bash

# Setup test HTTP servers for IPVS round robin testing
# Run this script in separate terminals or as background processes

echo "Setting up test HTTP servers for IPVS testing..."

# Server 1 - Port 9001
echo "Starting server on port 9001..."
python3 -c "
import http.server
import socketserver
import sys

class MyHandler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()
        self.wfile.write(b'<h1>Backend Server 1 (Port 9001)</h1><p>Request handled by backend 1</p>')

PORT = 9001
with socketserver.TCPServer(('', PORT), MyHandler) as httpd:
    print(f'Server running on port {PORT}')
    httpd.serve_forever()
" &

# Server 2 - Port 9002
echo "Starting server on port 9002..."
python3 -c "
import http.server
import socketserver
import sys

class MyHandler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()
        self.wfile.write(b'<h1>Backend Server 2 (Port 9002)</h1><p>Request handled by backend 2</p>')

PORT = 9002
with socketserver.TCPServer(('', PORT), MyHandler) as httpd:
    print(f'Server running on port {PORT}')
    httpd.serve_forever()
" &

# Server 3 - Port 9003
echo "Starting server on port 9003..."
python3 -c "
import http.server
import socketserver
import sys

class MyHandler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()
        self.wfile.write(b'<h1>Backend Server 3 (Port 9003)</h1><p>Request handled by backend 3</p>')

PORT = 9003
with socketserver.TCPServer(('', PORT), MyHandler) as httpd:
    print(f'Server running on port {PORT}')
    httpd.serve_forever()
" &

echo "All test servers started!"
echo "Test servers are running on:"
echo "  - http://127.0.0.1:9001 (Backend 1)"
echo "  - http://127.0.0.1:9002 (Backend 2)"
echo "  - http://127.0.0.1:9003 (Backend 3)"
echo ""
echo "Press Ctrl+C to stop all servers"
wait
