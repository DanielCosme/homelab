#!/bin/sh

set -euo pipefail

echo CLUSTER NAME: $CLUSTER_NAME

flux --kubeconfig ~/.kube/$CLUSTER_NAME \
	bootstrap gitea \
	--token-auth \
	--hostname=$GITEA_HOST \
	--owner=daniel \
	--repository=homelab \
	--private=false \
	--branch=main \
	--personal=true \
	--path=./clusters/hydra
