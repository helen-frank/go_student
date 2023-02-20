image-registry , image-tag

idea来源：当处于离线环境时，使用自建镜像仓库安装 karmada 。`karmadactl init`需要为每一个镜像单独指定镜像（kube image除外）,会导致命令变得十分冗长。于是期望可以通过设置 karmada 的镜像仓库地址来达到简化命令的目的。



实现预期：当设置`image-registry` , `image-tag`且组件的镜像为默认时使用`image-registry`，`image-tag`替换默认镜像仓库

对于离线环境，`docker tag` karmada 相关镜像并统一push到镜像仓库可以实现对于安装 karmada 相关组件自定义版本的管理



Q: 仅设置`image-registry`

A: 组件的镜像为默认时，替换默认镜像仓库，tag 使用默认



Q: 仅设置`image-tag`



v1 k8s 1.18 karmada 1.3 -- install image-tag 



