#!/bin/sh

CONTEXT=ingest
SCRIPT_FILE=$1
SCRIPT_NAME=${SCRIPT_FILE%.*}
CLEANED_SCRIPT=`grep -v '^\s*//' $SCRIPT_FILE | tr '\n' ' '`
URL="localhost:9200/_scripts/$SCRIPT_NAME/$CONTEXT"

UPLOAD=`curl --fail-with-body -s -X POST $URL -H 'Content-Type: application/json' -d'
{
  "script": {
    "lang": "painless",
    "source": "'"$CLEANED_SCRIPT"'"
  }
}
'`

if [ $? -eq 0 ]; then
  echo "\033[0;32mScript '$SCRIPT_NAME' uploaded succesfully!"
else
  echo "\033[0;31mError uploading '$SCRIPT_NAME':"
  echo "$UPLOAD" | jq
fi

