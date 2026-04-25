#!/bin/sh

set -eu

echo HOST: $CLUSTER_HOST
echo USER: $LINUX_ADMIN
echo WORKER: $HYDRA_WORKER

NODE_TOKEN=$(ssh $LINUX_ADMIN@$CLUSTER_HOST "sudo cat /var/lib/rancher/k3s/server/node-token")
echo $NODE_TOKEN
# curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="agent --server https://k3s.example.com --token mypassword" sh -s -
# curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="agent" K3S_TOKEN="mypassword" sh -s - --server https://k3s.example.com
# curl -sfL https://get.k3s.io | K3S_URL=https://k3s.example.com sh -s - agent --token mypassword
# curl -sfL https://get.k3s.io | K3S_URL=https://k3s.example.com K3S_TOKEN=mypassword sh -s - # agent is assumed because of K3S_URL
