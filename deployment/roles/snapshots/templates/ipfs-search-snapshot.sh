#!/bin/sh

curl -s -XPUT http://127.0.0.1:9200/_snapshot/{{ snapshot_name }}/snapshot_`date +'%y%m%d_%H%M'` -H 'Content-Type: application/json' -d '{"indices": "ipfs_*"}' | jq -e '.accepted' > /dev/null
