#!/bin/bash

# All images were downlaoded from this website: https://cupcake.nilssonlee.se/
# It has creative commons license, so I can use it for this project
file_path="./imageList.txt"

# Set the download directory
download_path="./images"

# Create the download directory if it doesn't exist
mkdir -p "$download_path"

# Loop through each line in the text file
while IFS= read -r url; do
  # Extract filename from URL
  filename="${url##*/}"

  # Download the image
  wget -O "$download_path/$filename" "$url"

  # Check for successful download
  if [[ $? -eq 0 ]]; then
    echo "Downloaded $filename"
  else
    echo "Error downloading $filename"
  fi
done < "$file_path"

echo "Finished downloading images!"
