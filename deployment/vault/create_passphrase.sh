#!/bin/bash
MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

pwgen -cynC | head -1 | gpg -e -o $MY_DIR/passphrase.gpg
