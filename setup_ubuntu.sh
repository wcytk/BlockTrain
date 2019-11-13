#!/bin/bash

# Ubuntu 下的安装
# Setup for Ubuntu

# 安装方法，执行bash setup.sh
# Installation: exec "bash setup.sh"

# 首先探测当前用户
# Detect user first
user=$( whoami )

# 判断当前用户身份是否是root
# Detect whether user is root
# user=$(env | grep USER | cut -d "=" -f 2 | head -1)
if [ "$user" == "root" ]; then
	root=$user
	act=""
else
	root=$( sudo -s whoami )
	act="sudo -s"
fi

if [ "$root" == "root" ]; then
	echo "You are in root mode!"
else
	echo "You are not in root mode!"
	echo "Please use \"sudo su\" to change into root mode, or \"sudo bash setup_ubuntu.sh\""
	exit 1
fi

# 请以root身份运行(su sudo)
# Please run it in root(su sudo)

# 获取当前工作目录
# Obtain the directory you are working on
DIR=$( pwd )

$act apt-get update -y
$act apt-get install tar -y
$act apt-get install wget -y
$act apt-get install dpkg -y
$act apt-get install python3 -y
$act apt-get install python3-pip -y
$act apt-get install python3-venv -y

python3 -m pip install --upgrade pip -i https://pypi.tuna.tsinghua.edu.cn/simple

if [ "$user" == "root" ]; then
	profile="/root/.bashrc"
else
	profile="/home/$user/.bashrc"
fi

# 检测系统是否存在go语言环境并进行安装
# Detect for go environment and set it up
source "$profile" &&
if command -v go > /dev/null; then
	echo "Your go environment has been installed!"
else
	echo "$profile"
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
		chown -R $1 /usr/local/go
	else
		echo "Directory $GODIR already exists ..."
	fi
	GOARCH=$(dpkg --print-architecture)
	echo "export GO111MODULE=on" >> $profile
	echo "export GOROOT=/usr/local/go" >> $profile
	echo "export GOOS=linux" >> $profile
	echo "export GOARCH=$GOARCH" >> $profile
	echo "export GOPATH=$GOPATH" >> $profile
	echo "export GOBIN=\$GOROOT/bin/" >> $profile
	echo "export GOTOOLS=\$GOROOT/pkg/tool" >> $profile
	echo "export PATH=\$PATH:\$GOBIN:\$GOTOOLS" >> $profile
	if [ "$profile" != "/root/.bashrc" ]; then
		rootPath="/root/.bashrc"
		echo "export GO111MODULE=on" >> $rootPath
		echo "export GOROOT=/usr/local/go" >> $rootPath
		echo "export GOOS=linux" >> $rootPath
		echo "export GOARCH=$GOARCH" >> $rootPath
		echo "export GOPATH=$GOPATH" >> $rootPath
		echo "export GOBIN=\$GOROOT/bin/" >> $rootPath
		echo "export GOTOOLS=\$GOROOT/pkg/tool" >> $rootPath
		echo "export PATH=\$PATH:\$GOBIN:\$GOTOOLS" >> $rootPath
	fi
fi

# 检测和安装ipfs环境
# Detect and install ipfs environment
source $profile && $act bash $DIR/install_ipfs.sh

# Create python virtual environment
# 创建python虚拟环境
python3 -m venv $DIR/train/venv

# -i选项是使用中国的加速镜像站，如果不需要可以去除
# -i is for mirror in China, if you don't need it, just delete from "-i"
source $DIR/train/venv/bin/activate && python -m pip install --upgrade pip setuptools && pip install -r $DIR/train/requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple

# 清除安装包
# Cleaning up package
echo "Cleaning up ..."
$act rm -rf /root/Python-3.6.6.tar.xz
$act rm -rf /root/go1.13.linux-amd64.tar.gz

# 更新环境变量
# Update environment
source $profile