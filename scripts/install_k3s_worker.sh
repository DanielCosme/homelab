#!/bin/sh

set -eu

echo HOST: $CLUSTER_HOST
echo USER: $LINUX_ADMIN
echo WORKER: $HYDRA_WORKER

NODE_TOKEN=$(ssh $LINUX_ADMIN@$CLUSTER_HOST "sudo cat /var/lib/rancher/k3s/server/node-token")
echo $NODE_TOKEN

# curl -sfL https://get.k3s.io | K3S_URL=https://hydra-0:6443 K3S_TOKEN=<node_token> sh -s -
curl -sfL https://get.k3s.io | K3S_URL=https://hydra-0:6443 K3S_TOKEN=K106b6d0a91ca6467cef926389ece48b9a39b8ed571c82fe86a8a3cd55882e9c827::server:3f608bf732e3853cad1bfbb7121df16c sh -s -
