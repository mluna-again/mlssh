#! /usr/bin/env bash

git pull
go build || exit
sudo mv mlssh /usr/local/bin || exit
sudo setcap CAP_NET_BIND_SERVICE=+eip /usr/local/bin/mlssh
sudo systemctl restart mlssh
