#!/bin/bash

# 检测和安装IPFS
# Detect and install IPFS

DIR=$( pwd )

if command -v ipfs >> /dev/null; then
    echo "Your ipfs environment is satisfied!"
else
    echo "Your ipfs envitonment has not been installed!"

    IPFSFILE="/root/go1.12.8.linux-amd64.tar.gz"
	IPFSDIR="/root/go-ipfs"

    if [ ! -f "$IPFSFILE" ]; then
        # 由于国内下载需要代理，所以这里直接给出压缩包，如果有需要可以自行下载
        # The go-ipfs tar.gz has been given here, you can download yourself if needed
        # wget -P /root https://dist.ipfs.io/go-ipfs/v0.4.22/go-ipfs_v0.4.22_linux-amd64.tar.gz 
		cp ./go-ipfs_v0.4.22_linux-amd64.tar.gz /root
	else
		echo "File already exists, using cache ..."
	fi

    if [ ! -d "$IPFSEDIR" ]; then
		mkdir $IPFSDIR
		echo "Directory $IPFSDIR created!"
		tar -xzvf /root/go-ipfs_v0.4.22_linux-amd64.tar.gz -C /root
	else
		echo "Directory already exists ..."
	fi

    # Ipfs 自带的环境安装配置
    # Installation in ipfs
    cd /root/go-ipfs && ./install.sh

    # 加入私有网络（由于这里仅作测试，所以不会将加入共有网络，而是创建具有相同私钥的swarm私有网络）
    # Enter private network(This is only for testing, so we won't enter public ipfs-network)
    
    # 使用go安装ipfs-swarm的密钥产生工具
    # Using go to install ipfs-swarm key generator
    go get -u github.com/Kubuxu/go-ipfs-swarm-key-gen/ipfs-swarm-key-gen

    # 这里直接给出一个私有网络的swarm.key，如果有需要可以自己产生一个新的swarm.key以创建自己的私有网络
    # A private network's swarm.key is been given here, if you need a new private network, you can generate yourself
    # ipfs-swarm-key-gen > /root/.ipfs/swarm.key
    cp $DIR/swarm.key /root/.ipfs
 fi

 # 清除安装包
# Cleaning up package
echo "Cleaning up ..."
rm -rf /root/go-ipfs
rm -rf /root/go-ipfs_v0.4.22_linux-amd64.tar.gz