#!/bin/sh

set -eu

echo HOST: $CLUSTER_HOST
echo USER: $LINUX_ADMIN

ssh $LINUX_ADMIN@$CLUSTER_HOST "sudo mkdir -p /etc/rancher/k3s"
rsync --rsync-path="sudo rsync" ./config/k3s/config.yaml $LINUX_ADMIN@$CLUSTER_HOST:/etc/rancher/k3s
ssh $LINUX_ADMIN@$CLUSTER_HOST "sudo curl -sfL https://get.k3s.io | sudo sh"
ssh $LINUX_ADMIN@$CLUSTER_HOST "sudo systemctl status k3s.service"
