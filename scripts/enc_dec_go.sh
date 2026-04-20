#!/bin/sh

set -eu

if [ -z "${AGE_KEY}" ]; then
	echo "unbound variable"
fi
if [ ! -f "${AGE_KEY}" ]; then
	echo "Error: ${AGE_KEY} file does not exist"
	exit 1
fi

mkdir -p $ENC_SECRETS_PATH
if [ "$1" = "enc" ]; then
	for FILE in $GO_SECRETS/*; do
	  FILENAME="${FILE##*/}"
    age --encrypt \
        --output $ENC_SECRETS_PATH/$FILENAME.age \
        --identity $AGE_KEY \
        $FILE
	done
elif [ "$1" = "dec" ]; then
	mkdir -p $GO_SECRETS
	for FILE in $ENC_SECRETS_PATH/*; do
	  FILENAME="${FILE##*/}" 			# Cut the directories path
	  FILENAME="${FILENAME%.age}" # Cut the .age
		echo $GO_SECRETS/$FILENAME
    age --decrypt \
		    --identity $AGE_KEY \
        $FILE > $GO_SECRETS/$FILENAME
	 done
else
  echo "Error: Invalid argument. Use 'enc' or 'dec'." >&2
  exit 1
fi
