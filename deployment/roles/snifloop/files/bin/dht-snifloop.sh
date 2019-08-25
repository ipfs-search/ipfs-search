#!/bin/bash

# Tail IPFS log
# Get hashes for handleAddProvider messages
# Filter IPFS README and demo hashes
# Add null character every 100 lines
# Sort and filter unique entries after every null character and call ipfs-search add for every hash

while true
do
	echo "Script stopped, starting again"
	ipfs log tail |
        jq -r 'if .Operation == "handleAddProvider" then .Tags.key else empty end' |
        grep -v QmY5heUM5qgRubMDD1og9fhCPA6QdkMp3QCwd4s7gJsyE7 |
        grep -v QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp |
        sed '0~100 s/$/\x00/g' |
        xargs -0 -I % bash -c "echo \"%\" | sort -u | xargs -L1 ipfs-search a"
done

