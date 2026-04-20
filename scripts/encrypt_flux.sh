#!/bin/sh

set -eu

if [ -z "${AGE_KEY}" ]; then
	echo "unbound variable"
fi
if [ ! -f "${AGE_KEY}" ]; then
	echo "Error: ${AGE_KEY} file does not exist"
	exit 1
fi

PUBLIC_KEY=$(age-keygen -y $AGE_KEY)

SECRETS_ENC_PATH=$HYDRA_SECRETS_PATH
mkdir -p $SECRETS_ENC_PATH
for FILE in $SECRETS_FOLDER/*; do
	sops --encrypt --in-place $FILE

  FILENAME="${FILE##*/}"
  DEST=$SECRETS_ENC_PATH/$FILENAME
  echo Moving encrypted file to $DEST
  mv $FILE $DEST
done
