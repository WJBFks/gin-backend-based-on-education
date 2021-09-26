# gin-backend-based-on-education

基于[education](https://github.com/sxguan/education)的Fabric 2.2 and Gin纯后端

# 运行环境

Ubuntu 20.04

Go 1.17.1

docker 20.10.7

docker-compose 1.25.0

[参考Fabric中文文档](https://hyperledger-fabric.readthedocs.io/zh_CN/release-2.2/prereqs.html)

# 运行方式

在`/etc/hosts`中添加：

```
127.0.0.1  orderer.example.com

127.0.0.1  peer0.org1.example.com

127.0.0.1  peer1.org1.example.com
```

添加依赖：

```
cd education && go mod tidy
```
运行项目：

```
./clean_docker.sh
```

在`127.0.0.1:9000`进行访问

# 注意

不要修改`education`文件夹文件名

不要修改`fixtures`和`sdkInit`下的任何东西

以`education`为工作路径打开`vscode`

无视`education\chaincode\edu.go`的报错

链码：`education\chaincode\edu.go`

后端服务：`education\service\eduService.go`

Gin路由：`education\main.go`中的`ginInit()`函数

