#!/bin/sh

curl -s -XPUT http://127.0.0.1:9200/_snapshot/{{ snapshot_name }}/snapshot_`date +'%y%m%d_%H%M'` | jq -e '.accepted' > /dev/null
