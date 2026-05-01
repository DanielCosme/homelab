#!/bin/sh

set -eu

echo HOST: $CLUSTER_HOST
echo USER: $LINUX_ADMIN
echo WORKER: $HYDRA_WORKER

NODE_TOKEN=$(ssh $LINUX_ADMIN@$CLUSTER_HOST "sudo cat /var/lib/rancher/k3s/server/node-token")
echo $NODE_TOKEN

# curl -sfL https://get.k3s.io | K3S_URL=https://hydra-0:6443 K3S_TOKEN=<node_token> sh -s -
