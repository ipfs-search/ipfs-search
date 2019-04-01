#!/bin/bash

while true
do
	echo "Script stopped, starting again"
	ipfs log tail | \
        jq -r 'if .Operation == "handleAddProvider" then .Tags.key else empty end' |
        sed '0~1000 s/$/\x00/g' | \
        xargs -0 -I % bash -c "echo \"%\" | sort -u | xargs -L1 ipfs-search a"
done

