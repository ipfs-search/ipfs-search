#!/bin/bash

IPFS=""
SNAPSHOT_DIR=""
KEY_NAME=""

echo "Publishing most recent IPFS snapshot to: $KEY_NAME"

OLD_HASH=""
echo "Latest snapshot: $OLD_HASH"

# Add dir to IPFS
NEW_HASH=`$IPFS add -w --nocopy --fscache $SNAPSHOT_DIR`
STATUS=$?

if [[ $STATUS == "0" ]]; then
    echo "Success creating snapshot: $NEW_HASH"

    echo "Replacing pin for $OLD_HASH with $NEW_HASH"
    $IPFS pin update $OLD_HASH $NEW_HASH

    echo "Publishing to $KEY_NAME"
    $IPFS name publish --key $KEY_NAME $NEW_HASH

    echo "Performing filestore garbage collection"
    $IPFS filestore gc
else
    echo "Error creating new snapshot!"
    exit -1
fi


