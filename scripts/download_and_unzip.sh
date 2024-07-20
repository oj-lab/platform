#!/bin/bash

TARGET_PATH=$1
DOWNLOAD_URL=$2
OVERRIDE=$3

# Check if a parameter was provided
if [ -z "$TARGET_PATH" ]
then
  echo "Please provide the TARGET_PATH."
  exit 1
fi

if [ -z "$DOWNLOAD_URL" ]
then
  echo "Please provide the DOWNLOAD_URL."
  exit 1
fi

# Check dist exists
if [ -d "$TARGET_PATH" ]; then
  if [ "$OVERRIDE" != "OVERRIDE=true" ]; then
    echo "$TARGET_PATH already exists. Use OVERRIDE=true to override." >&2
    exit 0
  fi
  echo "\"$TARGET_PATH\" already exists, cleaning..." >&2
  rm -rf "$TARGET_PATH"
else
  mkdir -p "$TARGET_PATH"
fi

# Check zip installed
if ! [ -x "$(command -v zip)" ]; then
  echo 'Error: zip is not installed.' >&2
  exit 1
fi

curl -o /tmp/download.zip -L "$DOWNLOAD_URL"
unzip -q -o /tmp/download.zip -d "$TARGET_PATH"
echo "Unzipped to \"$TARGET_PATH\""
rm /tmp/download.zip
echo "Removed \"/tmp/download.zip\""