#!/bin/bash

# Ubuntu 下的安装
# Setup for Ubuntu

# 安装方法，在管理员(root)权限下，执行sudo bash setup.sh
# Installation: exec "sudo bash setup.sh" in root

# 判断当前用户身份是否是root
# Detect whether user is root
# user=$(env | grep USER | cut -d "=" -f 2 | head -1)
user=$( whoami )
if [ "$user" == "root" ]; then
	echo "You are in root mode!"
else
	echo "You are not in root mode!"
	echo "Please use \"sudo su\" to change into root mode, or \"sudo bash setup_ubuntu.sh\""
	exit 1
fi

# 请以root身份运行(su sudo)
# Please run it in root(su sudo)

# 获取当前工作目录
# Obtain the directory you working on
DIR=$( pwd )

apt-get update -y
apt-get install tar -y
apt-get install wget -y
apt-get install dpkg -y
apt-get install python3 -y
apt-get install python3-pip -y
apt-get install python3-venv -y

python3 -m pip install --upgrade pip -i https://pypi.tuna.tsinghua.edu.cn/simple

# 检测系统是否存在go语言环境并进行安装
# Detect for go environment and set it up
source /etc/profile &&
if command -v go > /dev/null; then
	echo "Your go environment has been installed!"
else
	echo "Your go environment hasn't been installed!"
	echo "Where do you want to set the GOPATH:"
	read GOPATH
	if [ ! -d "$GOPATH" ]; then
		mkdir -p $GOPATH
		echo "Directory $GOPATH created!"
	else
		echo "Directory $GOPATH already exists ..."
	fi
	echo "Now start installing..."
	GOFILE="/root/go1.13.linux-amd64.tar.gz"
	GODIR="/usr/local/go/"
	if [ ! -f "$GOFILE" ]; then
		wget -P /root https://dl.google.com/go/go1.13.linux-amd64.tar.gz
	else
		echo "File already downloaded, using cache ..."
	fi
	if [ ! -d "$GODIR" ]; then
		mkdir $GODIR
		echo "Directory $GODIR created!"
		tar -xzvf /root/go1.13.linux-amd64.tar.gz -C /usr/local
		chown -R $user /usr/local/go
	else
		echo "Directory $GODIR already exists ..."
	fi
	GOARCH=$(dpkg --print-architecture)
	echo "export GO111MODULE=on" >> /etc/profile
	echo "export GOROOT=/usr/local/go" >> /etc/profile
	echo "export GOOS=linux" >> /etc/profile
	echo "export GOARCH=$GOARCH" >> /etc/profile
	echo "export GOPATH=$GOPATH" >> /etc/profile
	echo "export GOBIN=\$GOROOT/bin/" >> /etc/profile
	echo "export GOTOOLS=\$GOROOT/pkg/tool" >> /etc/profile
	echo "export PATH=\$PATH:\$GOBIN:\$GOTOOLS" >> /etc/profile
	source /etc/profile &&
	(if command -v go > /dev/null; then
		echo "Done! Now please source /etc/profile or restart a bash to use go!"
	else
		echo "Oops! Some issues occurs, try to examine the output or run this scripts again!"
	fi)
fi

# 检测和安装ipfs环境
# Detect and install ipfs environment
source /etc/profile && bash $DIR/install_ipfs.sh

# Create python virtual environment
# 创建python虚拟环境
python3 -m venv $DIR/train/venv

# -i选项是使用中国的加速镜像站，如果不需要可以去除
# -i is for mirror in China, if you don't need it, just delete from "-i"
source $DIR/train/venv/bin/activate && python -m pip install --upgrade pip setuptools && pip install -r $DIR/train/requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple

# 更新环境变量
# Update environment
source /etc/profile

# 清除安装包
# Cleaning up package
echo "Cleaning up ..."
rm -rf /root/Python-3.6.6.tar.xz
rm -rf /root/go1.13.linux-amd64.tar.gz
