#!/bin/bash
MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

gpg --batch --use-agent --decrypt $MY_DIR/passphrase.gpg
