#!/bin/bash

IPFS="env IPFS_PATH={{ ipfs_path }} ipfs"
SNAPSHOT_DIR="{{ snapshot_publish_path }}"
STATE_PATH="{{ snapshot_publish_state_path }}"
ADD_PARAMS="-w -r --nocopy --fscache --chunker size-1048576 --hash blake2b-256 --quieter --dereference-args"
# blake2b-256; performance (with large datasets)
# chunker: 1 MiB is IPFS max blocksize (more performance)
PUBLISH_KEY="{{ snapshot_publish_keyname }}"

echo "Publishing most recent IPFS snapshot"

OLD_HASH_FILE="$STATE_PATH/old_hash"
OLD_HASH=`cat $OLD_HASH_FILE`
if [ -n "$OLD_HASH" ]; then
	echo "Found old snapshot: $OLD_HASH"
else
	echo "Old hash not found, not updating!"
fi

# Add dir to IPFS
NEW_HASH=`$IPFS add $ADD_PARAMS $SNAPSHOT_DIR`
STATUS=$?

if [[ $STATUS == "0" ]]; then
    echo "Success creating snapshot: $NEW_HASH"

    echo "Writing $NEW_HASH to $OLD_HASH_FILE"
    echo $NEW_HASH > $OLD_HASH_FILE

	if [ -n "$OLD_HASH" ]; then
	    echo "Replacing pin for $OLD_HASH with $NEW_HASH"
	    $IPFS pin update $OLD_HASH $NEW_HASH
	fi

    $IPFS name publish --key=$PUBLISH_KEY $NEW_HASH
else
    echo "Error creating new snapshot!"
    exit -1
fi


