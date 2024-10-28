#!/bin/bash
while true; do
  echo "Connecting to WebSocket server..."
  wscat -c ws://localhost:8080/ws
  echo "Disconnected. Reconnecting in 2 seconds..."
  sleep 2
done
