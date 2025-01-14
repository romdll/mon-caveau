#!/bin/bash

cd ../../src/

VERSION="v$(git describe --tags --always --abbrev=7)"

export VPS_IP=192.168.1.100
export VPS_USER=user
export VPS_PASSWORD=password
export LOCAL_FILE=bin/MonCaveau-$VERSION
export REMOTE_PATH=/uploads/$VERSION/

echo "Creating remote directory on VPS..."
echo $VPS_PASSWORD | plink -batch -pw $VPS_PASSWORD $VPS_USER@$VPS_IP "sudo mkdir -p $REMOTE_PATH"

TEMP_PATH=/home/$VPS_USER/MonCaveau-$VERSION
echo "Uploading file to VPS (temporary location)..."
pscp -pw $VPS_PASSWORD $LOCAL_FILE $VPS_USER@$VPS_IP:$TEMP_PATH

echo "Moving file to final destination on VPS..."
echo $VPS_PASSWORD | plink -batch -pw $VPS_PASSWORD $VPS_USER@$VPS_IP "sudo mv $TEMP_PATH $REMOTE_PATH"

echo "Cleaning up temporary file..."
echo $VPS_PASSWORD | plink -batch -pw $VPS_PASSWORD $VPS_USER@$VPS_IP "sudo rm -f $TEMP_PATH"