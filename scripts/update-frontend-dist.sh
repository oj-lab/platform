#!/bin/bash

FRONTEND_DIST_PATH=$1

# Check if a parameter was provided
if [ -z "$FRONTEND_DIST_PATH" ]
then
  echo "No argument supplied. Please provide the path to FRONTEND_DIST_PATH."
  exit 1
fi

# Check zip installed
if ! [ -x "$(command -v zip)" ]; then
  echo 'Error: zip is not installed.' >&2
  exit 1
fi

# Check dist exists
if [ -d "$FRONTEND_DIST_PATH" ]; then
  echo "Info: $FRONTEND_DIST_PATH already exists, cleaning..." >&2
  rm -rf "$FRONTEND_DIST_PATH"
else
  mkdir -p "$FRONTEND_DIST_PATH"
fi

curl -o dist.zip -L https://github.com/OJ-lab/oj-lab-front/releases/download/v0.0.2/dist.zip
unzip -o dist.zip -d "$FRONTEND_DIST_PATH"
rm dist.zip