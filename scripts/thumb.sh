#!/usr/bin/bash

set -e

PDF_DEST=$1
TBN_DEST=$2

convert -thumbnail x200 ${PDF_DEST}[0] ${TBN_DEST}
