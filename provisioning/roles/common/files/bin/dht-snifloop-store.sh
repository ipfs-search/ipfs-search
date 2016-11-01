#!/bin/bash

while true
do
	echo "Script stopped, starting again"
	sudo -E -u ipfs ipfs log tail | jq -r 'if .event == "handleAddProvider" then .key else empty end' >> sniffed_hashes.txt	
done

