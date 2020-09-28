# Application
![Build & Deploy](https://github.com/gsakun/application/workflows/Build%20&%20Deploy/badge.svg?branch=master)

为解决业务基于k8s和istio进行服务发布的复杂性问题（主要包括各类不同资源的创建及搭配使用，用户的使用门槛。），自定义服务发布规范application，并开发服务发布控制器（应用方归一化配置待发布应用及所需服务治理功能整合为自定义资源，controller自动解析生成为各类k8s&istio资源配置并创建相应资源），实现服务的极简化发布。

# 主要功能
* 