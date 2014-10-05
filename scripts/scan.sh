#!/usr/bin/bash

# Stop the script on first error
set -e

DEST=$1
TIF_DEST=${DEST}.tiff
PDF_DEST=${DEST}.pdf
TBN_DEST=.thumb/${DEST}.png

scanimage --resolution 200 --mode color --format tiff > ${TIF_DEST}
tiff2pdf ${TIF_DEST} > ${PDFDEST}
convert -thumbnail x200 ${PDF_DEST} ${TBN_DEST}
rm ${TIF_DEST}
