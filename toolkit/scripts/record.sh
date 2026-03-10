#!/bin/bash
set -e

# Usage: record.sh <session_id> <output_path> <duration_limit>

SESSION_ID=$1
OUTPUT_PATH=$2
DURATION_LIMIT=$3

if [ -z "$SESSION_ID" ] || [ -z "$OUTPUT_PATH" ]; then
    echo "Error: Missing required arguments."
    echo "Usage: $0 <session_id> <output_path> [duration_limit]"
    exit 1
fi

echo "Starting recording for session $SESSION_ID"
echo "Output path: $OUTPUT_PATH"

mkdir -p "$OUTPUT_PATH"

# Simulate recording by sleeping and occasionally logging
# Real implementation would launch ROS/lerobot recording commands here
echo "Recording started at $(date)"

COUNT=0
LIMIT=${DURATION_LIMIT:-60}

# Handle termination gracefully
trap "echo 'Recording stopped gracefully'; exit 0" SIGINT SIGTERM

while [ $COUNT -lt $LIMIT ]; do
    sleep 1
    COUNT=$((COUNT+1))
    
    # Simulate some stdout telemetry
    if [ $((COUNT % 5)) -eq 0 ]; then
        echo "{\"timestamp\": \"$(date)\", \"status\": \"recording\", \"frames_captured\": $((COUNT * 30))}"
    fi
done

echo "Recording completed successfully after $LIMIT seconds."
exit 0
