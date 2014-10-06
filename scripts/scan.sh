#!/usr/bin/bash

# Stop the script on first error
set -e

DEST_DIR=$1
DEST=$2
TIF_DEST=${DEST_DIR}/${DEST}.tiff
PDF_DEST=${DEST_DIR}/${DEST}.pdf
TBN_DEST=${DEST_DIR}/.thumb/${DEST}.png

scanimage --resolution 200 --mode color --format tiff > ${TIF_DEST}
tiff2pdf ${TIF_DEST} > ${PDF_DEST}
`dirname $0`/thumb.sh ${PDF_DEST} ${TBN_DEST}
rm ${TIF_DEST}
