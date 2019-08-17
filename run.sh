# 运行go服务器
# run go server
:<<!
如果需要使用GOPATH，且未配置，请在home或者root文件夹下编辑.bashrc，并在末尾添加export GOPATH=/home/code(你想要添加的文件夹的绝对路径)
If you want to use GOPATH and it's not configured, then you can edit .bashrc in your current user's folder(like /home or /root), and add export GOPATH=/home/code(the absolute location of the folder you want to use)
!

# 检测使用到的端口是否被占用
# Checking for ports whether they are hosted
echo "Checking for ports ..."
bash stop.sh
echo "Done!"

echo "Now installing run.go ..."
go install run.go
echo "Now building run.go ..."
go build -o startTrainServer run.go
echo "Starting server ..."
./startTrainServer
