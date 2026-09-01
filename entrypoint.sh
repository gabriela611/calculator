#!/bin/sh
# Start Go backend server on port 8081 in background
PORT=8081 /app/server &

# Start Nginx in foreground on port 8080
exec nginx -g 'daemon off;'
