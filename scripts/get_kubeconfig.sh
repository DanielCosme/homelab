#!/bin/sh

set -eu

echo CLUSTER NAME: $CLUSTER_NAME
echo HOST: $CLUSTER_HOST
echo USER: $LINUX_ADMIN

rsync --rsync-path="sudo rsync" $LINUX_ADMIN@$CLUSTER_HOST:/etc/rancher/k3s/k3s.yaml ~/.kube/$CLUSTER_NAME
sudo sed -i "s/127.0.0.1/$CLUSTER_HOST/g" ~/.kube/$CLUSTER_NAME
kubectl --kubeconfig ~/.kube/$CLUSTER_NAME get pods -A
