#!/bin/bash

# 这个脚本用来停止服务器以及关闭训练占用的端口进程
# This script is for stopping server and shutdown the training processes

# 关闭服务器进程
# Stop server's process
echo "Shuting down server ..."
server=$( ps -ef | grep run | grep -v grep | awk '{print $2}' )

if [ "$server" == "" ]; then
	echo "Server is not started!"
else
	kill -9 "$server"
	echo "Server stoped!"
fi

# 关闭训练进程
# Shutdown trianing processes
echo "Shuting down ps server"
ps=$( lsof -i :2222 | grep "(LISTEN)" | awk '{printf $2}' )

if [ "$ps" == "" ]; then
	echo "ps server has been stopped or not started!"
else
	kill -9 "$ps"
	echo "ps server stoped!"
fi

for(( i = 0; i < 10; i++ ));
do
	port = 2225 + i
	worker=$( lsof -i :$port | grep "(LISTEN)" | awk '{printf $2}' )

	if [ "$worker" == "" ]; then
		echo "worker server on port $port has been stopped or not started!"
	else
		kill -9 "$worker"
		echo "worker server on port $port stopped!"
	fi
done
