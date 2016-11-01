#!/bin/sh

#ipfs log tail | jq -r 'if .event == "dhtSentMessage" and .message.type == "ADD_PROVIDER"  then .message.key else empty end'
ipfs log tail | jq -r 'if .event == "handleAddProvider" then .key else empty end' | uniq
