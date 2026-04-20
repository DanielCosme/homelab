#!/bin/sh

set -eu

if [ -z "${AGE_KEY}" ]; then
	echo "unbound variable"
fi
if [ ! -f "${AGE_KEY}" ]; then
	echo "Error: ${AGE_KEY} file does not exist"
	exit 1
fi

cat $AGE_KEY | kubectl --kubeconfig ~/.kube/hydra create secret generic sops-age --namespace=flux-system --from-file=age.agekey=/dev/stdin
