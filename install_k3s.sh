#!/bin/sh

sudo mkdir -p /etc/rancher/k3s
curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="--disable=helm-controller" sh
sudo systemctl status k3s.service
