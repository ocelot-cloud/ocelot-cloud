#!/bin/bash

set -e

# Define the filename
filename=myuser_myapp_v1.0.tar.gz

# Remove the users directory if it exists
rm -rf users

# Create a dummy file to upload
echo "hello" > $filename

# Upload the file using curl
response=$(curl -F "file=@$filename" -s -w "%{http_code}" -o /dev/null http://localhost:8082/upload)

# Check if the upload was successful
if [ "$response" -eq 200 ]; then
    echo "File uploaded successfully."
else
    echo "Failed to upload file. HTTP status code: $response"
    # Cleanup and exit if the upload failed
    rm $filename
    exit 1
fi

# Remove the local file after upload
rm $filename

# Download the file using curl
curl -O http://localhost:8082/download/$filename

# Verify the downloaded file exists
if [ -f "$filename" ]; then
    echo "File downloaded successfully."
else
    echo "Failed to download file."
fi
