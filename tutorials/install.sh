#! /bin/bash

local _pwd=$(pwd)

# install bcc and dependencies
echo "installing bcc dependencies"
install_package bison build-essential cmake flex git libedit-dev
install_package libllvm3.7 llvm-3.7-dev libclang-3.7-dev python zlib1g-dev libelf-dev
install_package luajit luajit-5.1-dev

# echo installing bcc
cd $DATA_DIR
rm -rf bcc  # be sure that there is not al already git clone of bcc
git clone https://github.com/iovisor/bcc.git
# there is an issue with master bcc and hover. use a previous version of
# bcc to avoid it.
cd bcc
git checkout d31f845e23899ff3b27b7e6f5a18330266fef2f3
mkdir build; cd build
cmake .. -DCMAKE_INSTALL_PREFIX=/usr
make
sudo make install

echo "installing go"
curl -O https://storage.googleapis.com/golang/go1.6.2.linux-amd64.tar.gz
tar -xvf go1.6.2.linux-amd64.tar.gz

sudo rm -rf /usr/local/go
sudo mv go /usr/local

sudo rm -f /usr/bin/go
sudo ln -s /usr/local/go/bin/go /usr/bin

export GOPATH=$HOME/go

# Hover
echo "installing hover dependencies"
go get github.com/vishvananda/netns
go get github.com/willf/bitset
go get github.com/gorilla/mux
# to pull customized fork of netlink
go get github.com/vishvananda/netlink
cd $HOME/go/src/github.com/vishvananda/netlink
git remote add drzaeus77 https://github.com/drzaeus77/netlink
git fetch drzaeus77
git reset --hard drzaeus77/master

echo "installing hover"
go get github.com/iovisor/iomodules/hover
go install github.com/iovisor/iomodules/hover/hoverd

echo "installing iovisor ovn"
go get github.com/netgroup-polito/iovisor-ovn/daemon

cd $_pwd