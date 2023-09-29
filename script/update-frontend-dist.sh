#!/bin/bash

# Check zip installed
if ! [ -x "$(command -v zip)" ]; then
  echo 'Error: zip is not installed.' >&2
  exit 1
fi

# Check dist exists
if [ -d "frontend/dist" ]; then
  echo 'Info: frontend/dist already exists, cleaning...' >&2
  rm -rf frontend/dist
fi

curl -o dist.zip -L https://github.com/OJ-lab/oj-lab-front/releases/download/v0.0.1/dist.zip
unzip -o dist.zip -d frontend
rm dist.zip