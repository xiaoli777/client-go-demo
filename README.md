# client-go-demo

## Introduction

这个仓库记录一些 client-go 的用法以及 demo:

common: 连接 Kubernetes 集群

epCounter: 动态监控 Kubernetes 集群中endpoints的数量，并输出 Adresses 和 NotReadyAddresses

dpCreator: 读取 yaml 文件，转换成 deployment 结构体，修改 replicas 的数量后部署到 k8s 集群

imageUpdater: 读取 deployment 的 yaml 文件，替换容器的 image ，并更新到 k8s 集群

podRoller: 滚动升级 deployment 中的 Pods



## Reference

https://github.com/kubernetes/client-go

https://github.com/owenliang/k8s-client-go

