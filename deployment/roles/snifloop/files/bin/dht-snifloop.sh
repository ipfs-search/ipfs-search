#!/bin/bash

while true
do
	echo "Script stopped, starting again"
	ipfs log tail | jq -r 'if .Operation == "handleAddProvider" then .Tags.key else empty end' | uniq | xargs -L 1 ipfs-search a
done

