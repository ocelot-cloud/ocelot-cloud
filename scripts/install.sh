#!/bin/bash

if [ "$(id -u)" = "0" ]; then
    echo "This script should not be run as root" 1>&2
    exit 1
fi

echo "Installing required debian packages"
sudo apt-get update
# gcc is required to compile Go applications with sqlite3 support
sudo apt-get install -y wget git curl sqlite3 gcc
# The libraries are needed by cypress: https://docs.cypress.io/guides/getting-started/installing-cypress#UbuntuDebian
sudo apt-get install -y libgtk2.0-0 libgtk-3-0 libgbm-dev libnotify-dev libnss3 libxss1 libasound2

echo "Installing docker"
curl -fsSL https://get.docker.com | sudo sh
sudo usermod -aG docker $USER

echo "Installing go"
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz -O go.tar.gz
sudo tar -C /usr/local -xzf go.tar.gz
rm go.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> $HOME/.bashrc
mkdir -p "$HOME/.go-cache"
echo "export GOPATH=$HOME/.go-cache" >> $HOME/.bashrc

echo "Installing node.js"
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
source $HOME/.nvm/nvm.sh
nvm install 18.10
nvm use 18.10
echo 'export PATH=$PATH:/usr/local/lib/node_modules/.bin' >> $HOME/.bashrc

source $HOME/.bashrc
echo "Please reboot the system manually to complete the installation."