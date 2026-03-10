#!/bin/bash
set -e

# Usage: upload_local.sh <job_id> <dataset_path> <platform_url>

JOB_ID=$1
DATASET_PATH=$2
PLATFORM_URL=$3

if [ -z "$JOB_ID" ] || [ -z "$DATASET_PATH" ] || [ -z "$PLATFORM_URL" ]; then
    echo "Error: Missing required arguments."
    echo "Usage: $0 <job_id> <dataset_path> <platform_url>"
    exit 1
fi

echo "Starting upload for job $JOB_ID"
echo "Dataset path: $DATASET_PATH"
echo "Target URL: $PLATFORM_URL"

# Simulate an upload process that takes some time
echo "Uploading chunk 1..."
sleep 2

# Simulate an intermittent failure for robustness demonstration
# In a real app, this would be actual curl/HTTP client logic to the platform
if [ $((RANDOM % 4)) -eq 0 ]; then
    echo "Network error during upload chunk 2!" >&2
    exit 1
fi

echo "Uploading chunk 2..."
sleep 2

echo "Uploading chunk 3..."
sleep 2

echo "Upload completed successfully."
exit 0
