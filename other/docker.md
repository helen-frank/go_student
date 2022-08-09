# 1.需求

- `docker`发布一个项目带上环境打包过去 

打包项目带上环境（镜像）--->(docker仓库：商店) --->下载发布的镜像--->直接运行即可 

- 多个应用（端口冲突） 

隔离：`docker`核心思想！打包装箱！每个箱子是相互隔离的 

`docker`通过隔离机制，可以将服务器利用到极致  

- `docker`重在轻量化 

虚拟机：虚拟出一个真实机的物理环境，大，全 

`docker`是容器技术，也是一种虚拟化技术，隔离，镜像（最核心的环境 4m+其他应用环境），运行镜像就行 



## 1.1 普遍的虚拟化技术

![image-20210531134630897](https://raw.githubusercontent.com/helen-frank/img/master/img/202208082159759.png)

**虚拟化的技术缺点** 

1. 资源占用十分多 

2. 冗余步骤多 

3. 启动很慢 

## 1.2 容器化技术

容器化技术不是模拟一个完整的操作系统 

![image-20210531134721889](https://raw.githubusercontent.com/helen-frank/img/master/img/202208082159723.png)

## 1.3 比较docker和虚拟机技术的不同

- 传统虚拟机，虚拟出一条硬件，运行一个完整的操作系统，然后再这个系统上安装和运行软件 
- 容器内的应用直接运行在宿主机的内核中，容器是没有自己的内核的，也没有虚拟硬件 
- 每个容器间是相互隔离的，每个容器都有一个属于自己的文件系统，互不影响 

# 2.基础

- 镜像`image` 

`docker`镜像就好比是一个模板，可以通过这个模板来创建容器服务，tomcat镜像--->`run`--->`tomcat01`容器（提供服务器），通过这个镜像可以创建多个容器，最终服务运行或者项目运行就是在容器中的 

- 容器`container` 

`docker`利用容器技术，独立运行一个或者一个组应用，通过镜像来创建的， 

启动，停止，删除，基本命令 

目前就可以把这个容器理解为就是一个简易的linux系统 

- 仓库`repository` 

仓库就是存放镜像的地方 

仓库分为共有仓库和私有仓库，Docker Hub（默认是国外的，国内阿里云等都有容量服务器

![image-20210531135025783](https://raw.githubusercontent.com/helen-frank/img/master/img/202208082159184.png)

## 2.1 docker的安装与卸载

### 2.1.1 for centos

1. 卸载旧的docker 

```shell
yum remove docker docker-client docker-client-latest docker-common docker-latest docker-latest-logrotate docker-logrotate docker-engine 
```

2. 需要的安装包

```shell
yum install -y yum-utils device-mapper-persistent-data lvm2  
```

3. 设置镜像仓库

```shell
yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo #官方源 

yum-config-manager  --add-repo  http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo                          # 阿里源 

#设置加速仓库(安装好后需要下载镜像之类的进行加速) 
 tee /etc/docker/daemon.json <<-'EOF'
{
	"registry-mirrors": ["https://8sqb5nwq.mirror.aliyuncs.com"] 
}
EOF

sudo systemctl daemon-reload 
sudo systemctl restart docker 
yum makecache fast  									#更新索引 
```

4. 安装docker相关

```shell
yum install -y docker-ce docker-ce-cli containerd.io 	#docker-ce社区版 docker-ee企业版 
```

5. 启动docker

```shell
systemctl start docker 
docker  version 			#显示版本，顺便查看是否安装完成 
```

6. 测试docker

```shell
docker run hello-world
```

7. 查看hello-world镜像

```shell
docker images
```

8. 卸载docker

```shell
yum remove docker-ce docker-ce-cli containerd.io 	#删除依赖 
rm -rf /var/lib/docker  							#删除资源
#/var/lib/docker     docker的默认工作路径 
```

### 2.1.2 for ubuntu

### 2.1.3 for arch







# 3.底层原理

## 3.1 docker是怎么工作的？

docker是一个Client-Server (CS) 机构的系统，dokcer的守护进程运行在主机上。通过Socket从客户端访问 

docker-Server接收到docker-Client的指令，就会执行这个命令 

![image-20210531140417721](https://raw.githubusercontent.com/helen-frank/img/master/img/202208082200936.png)

## 3.2 docker为什么比vm快 

1. `docker`有着比虚拟机更少的抽象层 

2. `docker`利用的是宿主机的内核，`vm`需要的是`GuestOS `

![image-20210531140500212](https://raw.githubusercontent.com/helen-frank/img/master/img/202208082200717.png)

so,新建一个容器的时候，docker不需要像虚拟机一样重新加载一个操作系统，避免引导。 

虚拟机是加载`GuestOS`，分钟级别的，慢，而docker是利用宿主机的操作系统，省略了这个复杂的过程，秒级 

![image-20210531140525134](https://raw.githubusercontent.com/helen-frank/img/master/img/202208082200499.png)

# 4.常用命令

## 4.1 帮助命令

| code               | tips                                       |
| ------------------ | ------------------------------------------ |
| docker version     | 显示docker的版本信息                       |
| docker info        | 显示docker的系统信息，包括镜像和容器的数量 |
| docker 命令 --help | 帮助命令                                   |

## 4.2 镜像命令

### docker images

```shell
查看所有本地的主机上的镜像 
[root@helen ~]# docker images 
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE 

hello-world         latest              bf756fb1ae65        10 months ago       13.3kB 

#解释 
REPOSITORY				#镜像的仓库源 
TAG						#镜像的标签 
IMAGE ID				#镜像的id 
CREATED					#镜像的创建时间 
SIZE					#镜像的大小 

#可选项 
-a, -all				#列出所有镜像 
-q, --quiet				#只显示镜像的id 
```

### docker search 

```shell
搜索镜像 
[root@helen ~]# docker search mysql 
NAME                              DESCRIPTION                                     STARS               OFFICIAL            AUTOMATED 

mysql                             MySQL is a widely used, open-source relation…   10142               [OK]                 

mariadb                           MariaDB is a community-developed fork of MyS…   3735                [OK]                 

#可选项，通过收藏量来过滤 
--filter=STARS=3000               #搜索出来的镜像就是STARS大于3000的 
```

### docker pull 

```shell
下载镜像 docker pull 镜像名[:tag] 
[root@helen ~]# docker pull mysql 
Using default tag: latest                 #如果不写tag，默认就是latest 
latest: Pulling from library/mysql 
bb79b6b2107f: Pull complete               #分层下载，docker iamges的核心，联合文件系统 
49e22f6fb9f7: Pull complete  
842b1255668c: Pull complete  
9f48d1f43000: Pull complete  
c693f0615bce: Pull complete  
8a621b9dbed2: Pull complete  
0807d32aef13: Pull complete  
a56aca0feb17: Pull complete  
de9d45fd0f07: Pull complete  
1d68a49161cc: Pull complete  
d16d318b774e: Pull complete  
49e112c55976: Pull complete  
Digest: sha256:8c17271df53ee3b843d6e16d46cff13f22c9c04d6982eb15a9a47bd5c9ac7e2d #签名 
Status: Downloaded newer image for mysql:latest 
docker.io/library/mysql:latest  #真实地址信息 
#docker pull mysql  等同于 docker pull docker.io/library/mysql:latest 

指定版本下载 
[root@helen ~]# docker pull mysql:5.7 
5.7: Pulling from library/mysql 
bb79b6b2107f: Already exists  
49e22f6fb9f7: Already exists  .ini放到和二进制文件一个目录，并修改配置文件 配置文件：
port 为http代理功能保留，暂空
842b1255668c: Already exists  
9f48d1f43000: Already exists  
c693f0615bce: Already exists  
8a621b9dbed2: Already exists  
0807d32aef13: Already exists  
f15d42f48bd9: Pull complete  
098ceecc0c8d: Pull complete  
b6fead9737bc: Pull complete  
351d223d3d76: Pull complete  
Digest: sha256:4d2b34e99c14edb99cdd95ddad4d9aa7ea3f2c4405ff0c3509a29dc40bcb10ef 
Status: Downloaded newer image for mysql:5.7 
docker.io/library/mysql:5.7 
```

### docker rmi 

```shell
删除镜像（可以基于id,正则匹配删除） 
-f 强制删除 
docker rmi -f 镜像id							  #删除指定的镜像 
docker rmi -f 镜像id 镜像id 镜像id         	  #删除多个镜像 
docker rmi -f  $(docker images -aq)				#删除所有镜像 
```

## 4.3 容器命令 

### 4.3.1 容器基于镜像创建

先下载一个centos镜像用于测试 

```shell
docker pull centos 
```

### 4.3.2 新建容器并启动

```shell
docker run [可选参数] image
#参数说明

--name name			容器名字，用户来区分容器 
-d 					后台方式运行 
-it 				使用交互方式运行，进入容器查看内容 
-p 					指定容器的端口 -p 8080:8080 
					-p ip:主机端口:容器端口 
					-p 主机端口:容器端口 
					-p 容器端口 
					容器端口 
					
-P 					随机指定端口 
```

### 4.3.3 退出容器

```shell
exit 		容器停止并退出 
ctrl+p+q 	容器不停止退出 
```

### 4.3.4 列出所有运行的容器

```she
docker ps 	#列出当前正在运行的容器 

-a 			显示所有容器 
-n=? 		显示最近创建的容器，个数 
-q 			静默模式，只显示容器编号 
```

### 4.3.5 删除容器

```shell
docker rm 容器id 					#删除指定容器，不能删除正在运行的容器 
docker rm -f $(docker ps -aq) 	  #删除所有容器 
docker ps -a -q | xargs doker rm  #删除所有容器 
```

### 4.3.6 启动和停止容器

```shell
docker start 容器id 				#启动容器 
docker restart 容器id 			#重启容器 
docker stop 容器id				#停止当前正在运行的容器 
docker kill 容器id 				#强制停止当前容器 
```

### 4.3.7 常用其他命令

#### 后台启动容器

```shell
docker run –d 容器名 			#常见问题，docker容器使用后台运行，就必须有一个前台进程，docker发现没有应用，就会自动停止 
```

#### 查看日志

```shell
docker logs -tf 容器id 			#查看所有日志 
docker logs -tf --tail 10 容器id  #查看最新的10条日志 
```

#### 查看容器中进程信息

```she
docker top 容器id 
```

#### 查看镜像的元数据

```shell
docker inspect 容器id
```

#### 进入当前正在运行的容器

```shell
容器通常是使用后台方式运行的，需要进入容器，修改一些配置 
docker exec -it 容器id /bin/bash 	# 进入容器后开启一个新的终端 
docker attach  	容器id 			# 进入容器当前正在执行的终端，不会开启新的进程 
```

#### 文件拷贝

```shell
docker cp 容器id:容器内路径 主机路径 
```

# 5.镜像

镜像是一种轻量级，可执行的独立软件包，用来打包软件运行环境和基于运行环境开发的软件，它包含运行某个软件所需的所有内容，包括代码，运行时，库，环境变量和配置文件 

所有的应用，直接打包docker镜像，就可以直接运行 

如何得到镜像： 

- 远程仓库下载 
- 拷贝 
- 自行制作一个镜像DockerFile 

## 5.1 Docker镜像加载原理 

### 5.1.1 UnionFS(联合文件系统) 

UbionFS(联合文件系统)：UnionFS文件系统(UnionFS)是一种分层，轻量级并且高性能的文件系统，它支持对文件系统的修改作为一次提交来一层层的叠加，同时可以将不同目录挂载到同一虚拟文件系统下（unite several directories into a single virtual filesystem）.Union文件系统是Docker镜像的基础。镜像可以通过分层来进行继承，基于基础镜像（没有父镜像），可以制作各种具体的应用镜像 

特性：一次同时加载多个文件系统，但从外面看起来，只能看到一个文件系统，联合加载会把各层文件系统叠加起来，这样最终的文件系统会包含所有底层的文件和目录

### 5.1.2 Docker镜像加载原理 

docker的镜像实际上由一层一层的文件系统组成，这种层级的文件系统Unionfs. 

bootfs(boot file system)主要包含bootloader和kernel。bootloader主要引导加载kernel,linux刚启动会加载bootfs文件系统，在Docker镜像的最底层是bootfs。这一层与我们典型的linux/unix系统是一样的，包含boot加载器和内核。当boot加载完成之后整个内核就在内存中了，此时内存的使用权已由bootfs转交给内核，此时系统也会卸载bootfs 

rootfs(root file system) 在bootfs之上，包含的就是典型linux系统中的/dev,/proc,/bin/etc 等标准目录和文件。rootfs就是各种不同的操作系统发行版，比如ubuntu,centos等等 

![image-20210531145315164](https://raw.githubusercontent.com/helen-frank/img/master/img/202208082200399.png)

对于一个精简的OS,rootfs可以很小，只需要包含最基本的命令，工具和程序库就可以了，因为底层直接用Host的kernel，自己只需要提供rootfs就可以了。由此可见对于不同的linux发行版，bootfs基本是一致的，rootfs会有差别，因此不同的发行版可以共用bootfs 

### 5.1.3 总结

镜像层为只读，容器层可读写，层级复用，提高资源利用率，容器之下的都叫镜像层，当容器启动时，新的可写层被加载到镜像的顶部 

 

## 5.2 如何提交一个自己的镜像 

### 5.2.1 commit镜像

```shell
docker commit 		#提交容器成为一个新的副本（保存当前容器的状态）  
docker commit –m="提交的描述" -a="作者" 容器id 目标镜像名:[TAG] #docker commit -m="this tomcat test" -a="helen" 66cb687f9b8f tomcat_test:1.0 
#提交到了本地，本地可以直接使用这个镜像 
```

# 6.可视化

Portainer | Rancher (CI/CD)  

## 6.1 portainer

```shell
docker volume create portainer_data 
docker run -d -p 9000:9000 -p 8000:8000 --name portainer --restart always -v /var/run/docker.sock:/var/run/docker.sock -v portainer_data:/data portainer/portainer 
#运行 Portainer 的服务器上的端口 9000 来使用 Portainer 
#访问：http://127.0.0.1:9000/ 
```

# 7.容器数据卷

docker的理念：将应用和环境打包成一个镜像 

数据都在容器中，容器删除，数据就会丢失，就产生新的需求，数据可以持久化,数据存储在本地 

容器之间可以有一个数据共享的技术，docker容器中产生的数据，同步到本地 

容器的持久化和同步操作，容器间数据共享 

 

## 7.1 使用数据卷 

### 7.1.1 直接使用命令来挂载 

```shell
docker run -it -v 主机目录:容器内目录 centos /bin/bash 
docker inspect 容器id #查询 
```

![image-20210531145933689](https://raw.githubusercontent.com/helen-frank/img/master/img/202208082200950.png)

```shell
#试验：安装mysql 
docker run -d -p 3310:3306 -v /mnt/mysql/conf:/etc/mysql/conf.d -v /mnt/mysql/data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=123456 --name mysql_ceshi_001 mysql:5.7 

mysql -P3310 -h 172.16.109.130 -uroot -p123456 #远程连接 
#匿名和具名挂载（这里没有写容器外地址） 
#匿名挂载 
-v 容器内路径 
docker run -d -P --name nginx001 -v /etc/nginx nginx 
#具名挂载 
docker run -d -P -name nginx002 -v nginx002_nginx:/etc/nginx nginx 
#查看具名挂载的详细信息 
docker volume inspect nginx002_nginx 
#查看所有的volume情况 
docker volume ls 
#所有docker容器的卷，没有指定目录的情况下都是在/var/lib/docker/volumes 
#可以在容器路径后写权限ro(readonly 只读)，rw(readwrite 读写) 【xxx:xxx:ro】,默认 rw ,设为ro则在容器里不可以修改，只能在宿主机里修改 
```

## 7.2 Dockerfile

`Dockerfile`就是用来构建`docker`镜像的构建文件,命令脚本 

- 通过这个脚本可以生成镜像 

- 创建一个dockerfile文件，名字可以随意，建议`Dockerfile `

- 这里面的每一个命令就是镜像的一层 

```shell
#文件中的内容 指令(大写)  参数 
vim dockerfile1 

FROM centos                               	#基于centos 
VOLUME ["volume01","volume02"]  			#匿名挂载 
CMD echo "----end----" 
CMD /bin/bash 
```

通过dockerfile来构建镜像 

```shell
docker build -f dockerfile1 -t centos-private:1.0 . 
```

## 7.3 数据卷容器

例：多个mysql同步数据 、

![image-20210531150352707](https://raw.githubusercontent.com/helen-frank/img/master/img/202208082200975.png)

```shell
docker run -it --name centos_001 -v vl_001:/home centos:7.9.2009			#父容器 
docker run -it --name centos_002 --volumes-from centos_001 centos:7.9.2009	#继承 
```

# 8.DockerFile

dockerfile是用来构建docker镜像的文件 

构建步骤： 

1. 编写一个dockerfile文件 
2. `docker build`构建成一个镜像 

3. `docker run `运行镜像 

4. `docker push`发布镜像（DockerHub,阿里云镜像仓库） 

## 8.1 基础

1. 每个保留关键字（指令）都必须是大写字母 

2. 执行从上到下的顺序执行 

3. #表示注释 

4. 每一个指令都会创建提交一个新的镜像层，并提交

![image-20210531150614291](https://raw.githubusercontent.com/helen-frank/img/master/img/202208082200651.png)

dockerfile是面向开发的，发布项目，做镜像，就要编写dockerfile 

```shell
DockerFile 		#构建文件，定义了一切的步骤，源代码 
DockerImage 	#通过DockerFile构建生成的镜像，最终发布和运行的产品 
Docker容器 	  #容器是镜像运行起来提供服务 
```

## 8.2 DockerFile的指令 

| code       | tips                                                         |
| ---------- | ------------------------------------------------------------ |
| FROM       | 指定基础镜像（centos,ubuntu）                                |
| MAINTAINER | 指定作者信息                                                 |
| RUN        | 在命令前面加上RUN即可                                        |
| ADD        | COPY文件，会自动解压                                         |
| WORKDIR    | 镜像的工作目录                                               |
| VOLUME     | 设置卷，挂载主机目录                                         |
| EXPOSE     | 指定对外的端口                                               |
| CMD        | 指定这个容器启动的时候要运行的命令，只有最后一个会生效，可被替代 |
| ENTRYPOINT | 指定这个容器启动的时候要运行的命令，可以追加命令             |
| ONBUILD    | 但构建一个被继承DockerFile这个时候就会运行ONBUILD的指令。触发指令 |
| COPY       | 类似ADD,将文件copy到镜像中                                   |
| ENV        | 构建的时候设置环境变量                                       |

### 8.2.1 实例一

Docker Hub中99%的镜像都是从这个基础镜像过来的 FROM scratch 

1. 编写自定义centos镜像的DockerFIle 

```shell
vim centos_helen_private

FROM centos:7.9.2009
MAINTAINER helen-frank<helenfrank@protonmail.com>

ENV MYPATH /root 
WORKDIR $MYPATH

RUN yum -y install vim gcc gcc-c++ wget net-tools
RUN yum -y install tmux
RUN yum -y update
RUN yum -y upgrade

EXPOSE 20 21 22 80 

CMD echo $MYPATH 
CMD echo "----end----" 
CMD /bin/bash 
```

2. 构建镜像 

```shell
docker build -f centos_helen_private -t centos_helen_private:0.1 . 
#最后的.是指定DockerFile路径为当前目录 
```

### 8.2.2 实例二

1. 编写tomcat的镜像 

（需要tomcat，jdk的压缩包），dockerfile文件官方命令Dockerfile，build会自动寻找这个文件，不需要 -f 指定 

```shell
vim tomcat_helen_private 

FROM centos:7.9.2009 
MAINTAINER helen-frank<helenfrank@protonmail.com> 

ADD jdk-8u271-linux-x64.tar.gz /usr/local 
ADD apache-tomcat-9.0.41.tar.gz /usr/local 

RUN yum -y install vim 
ENV MYPATH /usr/local 
WORKDIR $MYPATH 

ENV JAVA_HOME /usr/local/jdk1.8.0_271 
ENV CLASSPATH $JAVA_HOME/lib/dt.jar:$JAVA_HOME/lib/tools.jar 
ENV CATALINA_HOME /usr/local/apache-tomcat-9.0.41 
ENV CATALINA_BASH /usr/local/apache-tomcat-9.0.41 
ENV PATH $PATH:$JAVA_HOME/bin:$CATALINA_HOME/lib:$CATALINA_HOME/bin 

EXPOSE 8080 

CMD /usr/local/apache-tomcat-9.0.41/bin/startup.sh && tail -F /usr/local/apache-tomcat-9.0.41/bin/logs/catalina.out 
```

2. 构建镜像 

```shell
docker build -f tomcat_helen_private -t tomcat_helen_private:0.1 . 
```

## 8.3 CMD和ENTRYPOINT的区别 

| code       | tips                                                         |
| ---------- | ------------------------------------------------------------ |
| CMD        | 指定这个容器启动的时候要运行的命令，只有最后一个会生效，可被替代（写完整命令替换） |
| ENTRYPOINT | 指定这个容器启动的时候要运行的命令，可以追加命令(在运行容器的时候，补充参数之类的) |

## 8.4 发布自己的镜像到dockerhub上 

```shell
docker login -u 账号 -p 密码 
```

发布tomcat_helen_private:0.1 

```shell
docker tag tomcat_helen_private:0.1 helenfrank/privatehub:tomcat_helen #先标记 
docker push helenfrank/privatehub:tomcat_helen 			#再push 
```

```shell
docker logout			#退出账号 
```

## 8.5 发布镜像到阿里云镜像上 

创建命名空间 

创建镜像仓库 

按照官方文档走就行 

![image-20210531151730939](https://raw.githubusercontent.com/helen-frank/img/master/img/202208082200028.png)

# 9.docer网络

## 9.1 docker0

docker0是dockerce的虚拟网卡，当启动docker容器时，docker就会给docker容器分配ip,docker0使用的是桥接模式，使用的技术是veth-pair 

veth-pair就是一对的虚拟设备接口，成对出现，一段连着协议，一段彼此相连 

openStac,Docker容器之间的连接，OVS的连接，都是使用veth-pair技术 

## 9.2 --link（不推荐使用）

容器的ip可能会改变，使用--link可以绑定容器，通过容器名来寻址 

实际就是在容器的 /etc/hosts 里写了主机名对应的ip 

docker0问题：不支持容器名访问 

## 9.3 自定义网络

```shell
docker network ls #查看所有的docker 网络 
```

### 9.3.1 网络模式 

| model     | tips                           |
| --------- | ------------------------------ |
| bridge    | 桥接 docker （默认)            |
| none      | 不配置网络                     |
| host      | 和宿主机共享网络               |
| container | 容器网络连接（使用少，局限大） |

直接启动的命令默认使用 --net bridge ，这个就是docker0 

```shell
docker run -d -P --name tomcat_001 --net bridge tomcat 
```

### 9.3.2 自定义网络(先创建网络，创建容器时--net指定) 

```shell
docker network create --driver bridge --subnet 192.168.0.0/16 --gateway 192.168.0.1 mynet 


--driver 	bridge 
--subnet 	192.168.0.0/16 
--gateway 	192.168.0.1 
```

使用该网络启动容器 

```shell
docker run -d -P --name tomcat_mynet_001 --net mynet tomcat 
docker run -d -P --name tomcat_mynet_002 --net mynet tomcat 
```

```shell
docker exec -it tomcat_mynet_001 ping tomcat_mynet_002 
#使用自定义网络能直接通过容器名来访问 
#不同的集群使用不同的网络，保证集群是安全和健康的 
```

### 9.3.3 网络联通 

一个容器能够访问到另一个网络的所有容器 

原理：将这个容器添加到该网络中，构建一块虚拟网卡 

```shell
docker network connect mynet tomcat_001 #从此tomcat_001就具有了mynet给ta分配的地址 
```

## 9.4 实战：部署redis集群

1. 创建网卡 

```shell
docker network create redis_network --subnet 172.38.0.0/16  
```

2. 通过脚本创建六个redis配置 

```shell
for port in $(seq 1 6); \ 
do \ 
mkdir -p /mydata/redis/node-${port}/conf 
touch /mydata/redis/node-${port}/conf/redis.conf 
cat << EOF >/mydata/redis/node-${port}/conf/redis.conf 
port 6379 
bind 0.0.0.0 
cluster-enabled yes 
cluster-config-file nodes.conf 
cluster-node-timeout 5000 
cluster-announce-ip 172.38.0.1${port} 
cluster-announce-port 6379 
cluster-announce-bus-port 16379 
appendonly yes 
EOF 
done 
```

3. 创建redis容器 

```shell
docker run -p 6371:6379 -p 16371:16379 --name redis-1 \ 
-v /mydata/redis/node-1/data:/data \ 
-v /mydata/redis/node-1/conf/redis.conf:/etc/redis/redis.conf \ 
-d --net redis_network --ip 172.38.0.11 redis:5.0.9-alpine3.11 redis-server /etc/redis/redis.conf 

docker run -p 6372:6379 -p 16372:16379 --name redis-2 \ 
-v /mydata/redis/node-2/data:/data \ 
-v /mydata/redis/node-2/conf/redis.conf:/etc/redis/redis.conf \ 
-d --net redis_network --ip 172.38.0.12 redis:5.0.9-alpine3.11 redis-server /etc/redis/redis.conf 

docker run -p 6373:6379 -p 16373:16379 --name redis-3 \ 
-v /mydata/redis/node-3/data:/data \ 
-v /mydata/redis/node-3/conf/redis.conf:/etc/redis/redis.conf \ 
-d --net redis_network --ip 172.38.0.13 redis:5.0.9-alpine3.11 redis-server /etc/redis/redis.conf 

docker run -p 6374:6379 -p 16374:16379 --name redis-4 \ 
-v /mydata/redis/node-4/data:/data \ 
-v /mydata/redis/node-4/conf/redis.conf:/etc/redis/redis.conf \ 
-d --net redis_network --ip 172.38.0.14 redis:5.0.9-alpine3.11 redis-server /etc/redis/redis.conf

docker run -p 6375:6379 -p 16375:16379 --name redis-5 \ 
-v /mydata/redis/node-5/data:/data \ 
-v /mydata/redis/node-5/conf/redis.conf:/etc/redis/redis.conf \ 
-d --net redis_network --ip 172.38.0.15 redis:5.0.9-alpine3.11 redis-server /etc/redis/redis.conf 

docker run -p 6376:6379 -p 16376:16379 --name redis-6 \ 
-v /mydata/redis/node-6/data:/data \ 
-v /mydata/redis/node-6/conf/redis.conf:/etc/redis/redis.conf \ 
-d --net redis_network --ip 172.38.0.16 redis:5.0.9-alpine3.11 redis-server /etc/redis/redis.conf 
```

4. 创建集群(进入redis-1中) 

```shell
docker exec -it redis-1 /bin/sh
```

```shell
redis-cli --cluster create 172.38.0.11:6379 172.38.0.12:6379 172.38.0.13:6379 172.38.0.14:6379 172.38.0.15:6379 172.38.0.16:6379 --cluster-replicas 1 
```

​	real code

```shell
/data # redis-cli --cluster create 172.38.0.11:6379 172.38.0.12:6379 172.38.0.13:6379 172.38.0.14:6379 172.38.0.15: 
6379 172.38.0.16:6379 --cluster-replicas 1 
>>> Performing hash slots allocation on 6 nodes... 
Master[0] -> Slots 0 - 5460 
Master[1] -> Slots 5461 - 10922 
Master[2] -> Slots 10923 - 16383 
Adding replica 172.38.0.15:6379 to 172.38.0.11:6379 
Adding replica 172.38.0.16:6379 to 172.38.0.12:6379 
Adding replica 172.38.0.14:6379 to 172.38.0.13:6379 
M: 94975e1f9ea4d07e2b9ec1220d6267dc813c88cf 172.38.0.11:6379 
   slots:[0-5460] (5461 slots) master 
M: 95c7f5cd617fd0fda3809711093cabced91a624c 172.38.0.12:6379 
   slots:[5461-10922] (5462 slots) master 
M: 5e69f989b66663067ed2cc00446556f442cf03e7 172.38.0.13:6379 
   slots:[10923-16383] (5461 slots) master 
S: 9ed9347806df471d23f05e15ef72f6193fe34345 172.38.0.14:6379 
   replicates 5e69f989b66663067ed2cc00446556f442cf03e7 
S: 0be276cf1f09198b2992a001265924cf7075ef01 172.38.0.15:6379 
   replicates 94975e1f9ea4d07e2b9ec1220d6267dc813c88cf 
S: 3044a00e48a33851fddbfee80af53fe4e8bb9a10 172.38.0.16:6379 
   replicates 95c7f5cd617fd0fda3809711093cabced91a624c 
Can I set the above configuration? (type 'yes' to accept): yes 
>>> Nodes configuration updated 
>>> Assign a different config epoch to each node 
>>> Sending CLUSTER MEET messages to join the cluster 
Waiting for the cluster to join 
.. 
>>> Performing Cluster Check (using node 172.38.0.11:6379) 
M: 94975e1f9ea4d07e2b9ec1220d6267dc813c88cf 172.38.0.11:6379 
   slots:[0-5460] (5461 slots) master 
   1 additional replica(s) 
M: 95c7f5cd617fd0fda3809711093cabced91a624c 172.38.0.12:6379 
   slots:[5461-10922] (5462 slots) master 
   1 additional replica(s) 
S: 9ed9347806df471d23f05e15ef72f6193fe34345 172.38.0.14:6379 
   slots: (0 slots) slave 
   replicates 5e69f989b66663067ed2cc00446556f442cf03e7 
S: 3044a00e48a33851fddbfee80af53fe4e8bb9a10 172.38.0.16:6379 
   slots: (0 slots) slave 
   replicates 95c7f5cd617fd0fda3809711093cabced91a624c 
M: 5e69f989b66663067ed2cc00446556f442cf03e7 172.38.0.13:6379 
   slots:[10923-16383] (5461 slots) master 
   1 additional replica(s) 
S: 0be276cf1f09198b2992a001265924cf7075ef01 172.38.0.15:6379 
   slots: (0 slots) slave 
   replicates 94975e1f9ea4d07e2b9ec1220d6267dc813c88cf 
[OK] All nodes agree about slots configuration. 
>>> Check for open slots... 
>>> Check slots coverage... 
[OK] All 16384 slots covered. 
/data #  
```

test code

```shell
/data # redis-cli -c 
127.0.0.1:6379> cluster info 
cluster_state:ok 
cluster_slots_assigned:16384 
cluster_slots_ok:16384 
cluster_slots_pfail:0 
cluster_slots_fail:0 
cluster_known_nodes:6 
cluster_size:3 
cluster_current_epoch:6 
cluster_my_epoch:1 
cluster_stats_messages_ping_sent:293 
cluster_stats_messages_pong_sent:290 
cluster_stats_messages_sent:583 
cluster_stats_messages_ping_received:285 
cluster_stats_messages_pong_received:293 
cluster_stats_messages_meet_received:5 
cluster_stats_messages_received:583 
127.0.0.1:6379> cluster nodes 
95c7f5cd617fd0fda3809711093cabced91a624c 172.38.0.12:6379@16379 master - 0 1610115893991 2 connected 5461-10922 
94975e1f9ea4d07e2b9ec1220d6267dc813c88cf 172.38.0.11:6379@16379 myself,master - 0 1610115895000 1 connected 0-5460 
9ed9347806df471d23f05e15ef72f6193fe34345 172.38.0.14:6379@16379 slave 5e69f989b66663067ed2cc00446556f442cf03e7 0 1610115894538 4 connected 
3044a00e48a33851fddbfee80af53fe4e8bb9a10 172.38.0.16:6379@16379 slave 95c7f5cd617fd0fda3809711093cabced91a624c 0 1610115894645 6 connected 
5e69f989b66663067ed2cc00446556f442cf03e7 172.38.0.13:6379@16379 master - 0 1610115895000 3 connected 10923-16383 
0be276cf1f09198b2992a001265924cf7075ef01 172.38.0.15:6379@16379 slave 94975e1f9ea4d07e2b9ec1220d6267dc813c88cf 0 1610115895520 5 connected 
127.0.0.1:6379>  
```

# 10.Docker Compose

定义，运行多个容器， 批量容器编排

`Dockerfile`

配置文件 YAML file `docker-compose.yml`

启动 `docker-compose up`

> tips

`Compose`是`Docker`官方的开源项目，需要安装

Dokcerfile让程序在任何地方运行，web服务，redis,mysql,nginx....多个容器

`docker-compose.yml`

```yaml
version: "3.9"  # optional since v1.27.0
services:
  web:
    build: .
    ports:
      - "5000:5000"
    volumes:
      - .:/code
      - logvolume01:/var/log
    links:
      - redis
  redis:
    image: redis
volumes:
  logvolume01: {}
```

- 服务`services` ，容器，应用（web,redis,mysql...）
- 项目`project`，一组关联的容器

## 10.1 安装

1. 安装Docker Compose当前稳定版本

```bash
sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
```

> 镜像加速

```bash
sudo curl -L "https://get.daocloud.io/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
```

2. 对二进制文件应用可执行权限

```bash
sudo chmod +x /usr/local/bin/docker-compose
```

3. 测试

```bash
docker-compose --version
```

## 10.2 升级

如果您从 Compose 1.2 或更早版本升级，请在升级 Compose 后移除或迁移现有容器。这是因为，从 1.3 版开始，Compose 使用 Docker 标签来跟踪容器，并且需要重新创建容器以添加标签。

如果 Compose 检测到创建时没有标签的容器，它会拒绝运行，这样您就不会得到两组。如果您想继续使用现有容器（例如，因为它们有您想要保留的数据卷），您可以使用 Compose 1.5.x 使用以下命令迁移它们

```bash
docker-compose migrate-to-labels
```

或者，如果您不担心保留它们，则可以删除它们。Compose 只是创建新的。

```bash
docker container rm -f -v myapp_web_1 myapp_db_1 ...
```

10.3 卸载

```bash
sudo rm /usr/local/bin/docker-compose
```

使用`pip`以下方式安装，则卸载 Docker Compose

```bash
pip uninstall docker-compose
```

## 10.3 基本使用

[Compose示例](https://docs.docker.com/compose/gettingstarted/)

### 10.3.1 设置

1. 创建项目文件夹

   ```bash
   mkdir composetest
   cd composetest
   ```

2. `app.py`

   ```python
   import time
   
   import redis
   from flask import Flask
   
   app = Flask(__name__)
   cache = redis.Redis(host='redis', port=6379)
   
   def get_hit_count():
       retries = 5
       while True:
           try:
               return cache.incr('hits')
           except redis.exceptions.ConnectionError as exc:
               if retries == 0:
                   raise exc
               retries -= 1
               time.sleep(0.5)
   
   @app.route('/')
   def hello():
       count = get_hit_count()
       return 'Hello World! I have been seen {} times.\n'.format(count)
   ```

在此示例中，`redis`是应用程序网络上的 redis 容器的主机名。我们使用 Redis 的默认端口`6379`

> 处理瞬态错误
>
> 注意`get_hit_count`函数的编写方式。如果 redis 服务不可用，这个基本的重试循环让我们可以多次尝试我们的请求。这在应用程序上线时的启动时很有用，但如果 Redis 服务需要在应用程序的生命周期内随时重新启动，这也会使我们的应用程序更具弹性。在集群中，这也有助于处理节点之间的瞬时连接中断。

3. `requirements.txt`

```txt
flask
redis
```

### 10.3.2 创建Dockerfile

编写一个用于构建 Docker 映像的 Dockerfile。该图像包含 Python 应用程序所需的所有依赖项，包括 Python 本身

`Dockerfile`

```dockerfile
# syntax=docker/dockerfile:1
FROM python:3.7-alpine
WORKDIR /code
ENV FLASK_APP=app.py
ENV FLASK_RUN_HOST=0.0.0.0
RUN apk add --no-cache gcc musl-dev linux-headers
COPY requirements.txt requirements.txt
RUN pip install -r requirements.txt
EXPOSE 5000
COPY . .
CMD ["flask", "run"]
```

- 从 Python 3.7 映像开始构建映像。
- 将工作目录设置为`/code`.
- 设置`flask`命令使用的环境变量。
- 安装 gcc 和其他依赖项
- 复制`requirements.txt`并安装 Python 依赖项。
- 将元数据添加到图像以描述容器正在侦听端口 5000
- 将`.`项目中的当前目录复制到`.`镜像中的workdir 。
- 将容器的默认命令设置为`flask run`.

### 10.3.3 在Compose文件中定义服务

```yaml
version: "3.9"
services:
  web:
    build: .
    ports:
      - "5000:5000"
  redis:
    image: "redis:alpine"
```

这个 Compose 文件定义了两个服务：`web`和`redis`

- 网络服务

该`web`服务使用从`Dockerfile`当前目录中构建的映像。然后它将容器和主机绑定到暴露的端口`5000`. 此示例服务使用 Flask Web 服务器的默认端口`5000`.

- Redis服务

该`redis`服务使用 从 Docker Hub 注册表中提取的公共[Redis](https://registry.hub.docker.com/_/redis/)映像。

### 10.3.4 使用Compose构建并运行项目

```bash
docker-compose up
```

### 10.3.5 编辑Compose文件以添加绑定挂载

`docker-compose.yml`

```yaml
version: "3.9"
services:
  web:
    build: .
    ports:
      - "5000:5000"
    volumes:
      - .:/code
    environment:
      FLASK_ENV: development
  redis:
    image: "redis:alpine"
```

新`volumes`密钥将主机上的项目目录（当前目录）挂载到`/code`容器内部，允许您即时修改代码，而无需重新构建映像。该`environment`键设置 `FLASK_ENV`环境变量，它告诉sudo chmod +x /usr/local/bin/docker-compose`flask run`在开发模式下运行，并重新加载更改代码。这种模式应该只在开发中使用。

共享文件夹、卷和绑定安装

- 如果您的项目在`Users`目录 ( `cd ~`) 之外，那么您需要共享您正在使用的 Dockerfile 和卷的驱动器或位置。如果您收到指示未找到应用程序文件、卷安装被拒绝或服务无法启动的运行时错误，请尝试启用文件或驱动器共享。对于位于`C:\Users`(Windows) 或`/Users`(Mac)之外的项目，卷安装需要共享驱动器 ，并且对于使用[Linux 容器的](https://docs.docker.com/docker-for-windows/#switch-between-windows-and-linux-containers)Docker Desktop for Windows 上的*任何*项目都需要。有关更多信息，请参阅Docker for Mac 上的[文件共享](https://docs.docker.com/docker-for-mac/#file-sharing)，以及有关如何[在容器中管理数据](https://docs.docker.com/storage/volumes/)的一般示例 。
- 如果您在较旧的 Windows 操作系统上使用 Oracle VirtualBox，您可能会遇到此[VB 故障单 中](https://www.virtualbox.org/ticket/14920)所述的共享文件夹[问题](https://www.virtualbox.org/ticket/14920)。较新的 Windows 系统满足[Docker Desktop for Windows](https://docs.docker.com/docker-for-windows/install/)的要求，不需要 VirtualBox。

### 10.3.6 使用Compose重新构建并运行项目

```bash
docker-compose up
```

### 10.3.7 更新应用程序

由于应用程序代码现在使用卷挂载到容器中，因此您可以对其代码进行更改并立即查看更改，而无需重新构建映像。

更改问候语`app.py`并保存。例如，将`Hello World!` 消息更改为`Hello from Docker!`：

```python
return 'Hello from Docker! I have been seen {} times.\n'.format(count)
```

### 10.3.8 其他命令

- 如果您想在后台运行您的服务，您可以将`-d`标志（用于“分离”模式）传递给`docker-compose up`并用于`docker-compose ps`查看当前正在运行的内容：

```bash
docker-compose up -d
docker-compsoe ps
```

- 该`docker-compose run`命令允许您为您的服务运行一次性命令。例如，要查看`web`服务可用的环境变量 ：

```bas
docker-compose run web env
```

- 停止服务

```bash
docker-compose stop
```

- 使用以下`down` 命令关闭所有内容，完全删除容器。传递`--volumes`给 Redis 容器使用的数据卷

```bash
docker-compose down --volumes
```

## 10.4 yaml 规则

docker-compose.yml 核心

[version3_example](https://docs.docker.com/compose/compose-file/compose-file-v3/#compose-file-structure-and-examples)

```yaml
# 3层
#版本
version: ''

#服务
services: 	
	服务1:web
		#服务配置
		images
		build
		network
		...
	服务2:redis
		...
	服务3:redis
	 	...
	 	
#其他配置
volumes:
configs:
```

## 10.5 搭建开源博客 wordpress

```bash
mkdir my_wordpress
cd my_wordpress
```

`docker-compose.yml`启动 `WordPress`博客和一个单独的`MySQL`实例，该实例具有用于数据持久性的卷挂载：

```yaml
version: "3.9"
    
services:
  db:
    image: mysql:5.7
    volumes:
      - db_data:/var/lib/mysql
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: somewordpress
      MYSQL_DATABASE: wordpress
      MYSQL_USER: wordpress
      MYSQL_PASSWORD: wordpress
    
  wordpress:
    depends_on:
      - db
    image: wordpress:latest
    volumes:
      - wordpress_data:/var/www/html
    ports:
      - "8000:80"
    restart: always
    environment:
      WORDPRESS_DB_HOST: db:3306
      WORDPRESS_DB_USER: wordpress
      WORDPRESS_DB_PASSWORD: wordpress
      WORDPRESS_DB_NAME: wordpress
volumes:
  db_data: {}
  wordpress_data: {}
```

```bash
docker-compose up -d
```



















# 11.Docker Swarm

集群方式的部署



## 11.1 基本操作

192.168.194.170作为主节点

1. 初始化主节点

`docker swarm init`

```bash
[root@localhost ~]# docker swarm init --advertise-addr 192.168.194.170
Swarm initialized: current node (vi279zrhj0trkxtmcoss09cu4) is now a manager.

To add a worker to this swarm, run the following command:

    docker swarm join --token SWMTKN-1-525upu0t3lyvsvjy8p5cgbf1tb9fw40js3u24ojqyi6nagkevn-8ov835xi1owc5zayrjct5d7z5 192.168.194.170:2377

To add a manager to this swarm, run 'docker swarm join-token manager' and follow the instructions.

[root@localhost ~]# docker swarm init --help

Usage:  docker swarm init [OPTIONS]

Initialize a swarm

Options:
      --advertise-addr string                  Advertised address (format: <ip|interface>[:port])
      --autolock                               Enable manager autolocking (requiring an unlock
                                               key to start a stopped manager)
      --availability string                    Availability of the node
                                               ("active"|"pause"|"drain") (default "active")
      --cert-expiry duration                   Validity period for node certificates
                                               (ns|us|ms|s|m|h) (default 2160h0m0s)
      --data-path-addr string                  Address or interface to use for data path
                                               traffic (format: <ip|interface>)
      --data-path-port uint32                  Port number to use for data path traffic (1024
                                               - 49151). If no value is set or is set to 0,
                                               the default port (4789) is used.
      --default-addr-pool ipNetSlice           default address pool in CIDR format (default [])
      --default-addr-pool-mask-length uint32   default address pool subnet mask length (default 24)
      --dispatcher-heartbeat duration          Dispatcher heartbeat period (ns|us|ms|s|m|h)
                                               (default 5s)
      --external-ca external-ca                Specifications of one or more certificate
                                               signing endpoints
      --force-new-cluster                      Force create a new cluster from current state
      --listen-addr node-addr                  Listen address (format: <ip|interface>[:port])
                                               (default 0.0.0.0:2377)
      --max-snapshots uint                     Number of additional Raft snapshots to retain
      --snapshot-interval uint                 Number of log entries between Raft snapshots
                                               (default 10000)
      --task-history-limit int                 Task history retention limit (default 5)
```

2. 加入节点`docker swarm join`

```bash
# 获取令牌,生成命令， 在manager节点上
# manager
docker swarm join-token manager
# worker
docker swarm join-token worker

# 加入节点，以worker
docker swarm join --token SWMTKN-1-525upu0t3lyvsvjy8p5cgbf1tb9fw40js3u24ojqyi6nagkevn-8ov835xi1owc5zayrjct5d7z5 192.168.194.170:2377
```

3. 在主节点上查看

```bash
docker node ls
```

4. 更改身份

```bash
docker swarm leave

#加入节点

```



5. 离开节点

```bash
docker swarm leave
```



## 11.2 Raft协议

保证大多数节点存活才可以用，集群至少三台`manager`节点

即三台中一台突然宕机，其他还能正常工作

Raft协议： 保证大多数节点存活， 高可用



## 11.3 service

创建服务，动态扩展服务，动态更新服务

```bash
[root@localhost ~]# docker service --help

Usage:  docker service COMMAND

Manage services

Commands:
  create      Create a new service
  inspect     Display detailed information on one or more services
  logs        Fetch the logs of a service or task
  ls          List services
  ps          List the tasks of one or more services
  rm          Remove one or more services
  rollback    Revert changes to a service's configuration
  scale       Scale one or multiple replicated services
  update      Update a service

Run 'docker service COMMAND --help' for more information on a command.

```

## 11.4 灰度发布（金丝雀发布）

```bash
[root@localhost ~]# docker service create --help

Usage:  docker service create [OPTIONS] IMAGE [COMMAND] [ARG...]

Create a new service

Options:
      --cap-add list                       Add Linux capabilities
      --cap-drop list                      Drop Linux capabilities
      --config config                      Specify configurations to expose to the service
      --constraint list                    Placement constraints
      --container-label list               Container labels
      --credential-spec credential-spec    Credential spec for managed service account
                                           (Windows only)
  -d, --detach                             Exit immediately instead of waiting for the service
                                           to converge
      --dns list                           Set custom DNS servers
      --dns-option list                    Set DNS options
      --dns-search list                    Set custom DNS search domains
      --endpoint-mode string               Endpoint mode (vip or dnsrr) (default "vip")
      --entrypoint command                 Overwrite the default ENTRYPOINT of the image
  -e, --env list                           Set environment variables
      --env-file list                      Read in a file of environment variables
      --generic-resource list              User defined resources
      --group list                         Set one or more supplementary user groups for the
                                           container
      --health-cmd string                  Command to run to check health
      --health-interval duration           Time between running the check (ms|s|m|h)
      --health-retries int                 Consecutive failures needed to report unhealthy
      --health-start-period duration       Start period for the container to initialize before
                                           counting retries towards unstable (ms|s|m|h)
      --health-timeout duration            Maximum time to allow one check to run (ms|s|m|h)
      --host list                          Set one or more custom host-to-IP mappings (host:ip)
      --hostname string                    Container hostname
      --init                               Use an init inside each service container to
                                           forward signals and reap processes
      --isolation string                   Service container isolation mode
  -l, --label list                         Service labels
      --limit-cpu decimal                  Limit CPUs
      --limit-memory bytes                 Limit Memory
      --limit-pids int                     Limit maximum number of processes (default 0 =
                                           unlimited)
      --log-driver string                  Logging driver for service
      --log-opt list                       Logging driver options
      --max-concurrent uint                Number of job tasks to run concurrently (default
                                           equal to --replicas)
      --mode string                        Service mode (replicated, global, replicated-job,
                                           or global-job) (default "replicated")
      --mount mount                        Attach a filesystem mount to the service
      --name string                        Service name
      --network network                    Network attachments
      --no-healthcheck                     Disable any container-specified HEALTHCHECK
      --no-resolve-image                   Do not query the registry to resolve image digest
                                           and supported platforms
      --placement-pref pref                Add a placement preference
  -p, --publish port                       Publish a port as a node port
  -q, --quiet                              Suppress progress output
      --read-only                          Mount the container's root filesystem as read only
      --replicas uint                      Number of tasks
      --replicas-max-per-node uint         Maximum number of tasks per node (default 0 = unlimited)
      --reserve-cpu decimal                Reserve CPUs
      --reserve-memory bytes               Reserve Memory
      --restart-condition string           Restart when condition is met
                                           ("none"|"on-failure"|"any") (default "any")
      --restart-delay duration             Delay between restart attempts (ns|us|ms|s|m|h)
                                           (default 5s)
      --restart-max-attempts uint          Maximum number of restarts before giving up
      --restart-window duration            Window used to evaluate the restart policy
                                           (ns|us|ms|s|m|h)
      --rollback-delay duration            Delay between task rollbacks (ns|us|ms|s|m|h)
                                           (default 0s)
      --rollback-failure-action string     Action on rollback failure ("pause"|"continue")
                                           (default "pause")
      --rollback-max-failure-ratio float   Failure rate to tolerate during a rollback (default 0)
      --rollback-monitor duration          Duration after each task rollback to monitor for
                                           failure (ns|us|ms|s|m|h) (default 5s)
      --rollback-order string              Rollback order ("start-first"|"stop-first")
                                           (default "stop-first")
      --rollback-parallelism uint          Maximum number of tasks rolled back simultaneously
                                           (0 to roll back all at once) (default 1)
      --secret secret                      Specify secrets to expose to the service
      --stop-grace-period duration         Time to wait before force killing a container
                                           (ns|us|ms|s|m|h) (default 10s)
      --stop-signal string                 Signal to stop the container
      --sysctl list                        Sysctl options
  -t, --tty                                Allocate a pseudo-TTY
      --ulimit ulimit                      Ulimit options (default [])
      --update-delay duration              Delay between updates (ns|us|ms|s|m|h) (default 0s)
      --update-failure-action string       Action on update failure
                                           ("pause"|"continue"|"rollback") (default "pause")
      --update-max-failure-ratio float     Failure rate to tolerate during an update (default 0)
      --update-monitor duration            Duration after each task update to monitor for
                                           failure (ns|us|ms|s|m|h) (default 5s)
      --update-order string                Update order ("start-first"|"stop-first") (default
                                           "stop-first")
      --update-parallelism uint            Maximum number of tasks updated simultaneously (0
                                           to update all at once) (default 1)
  -u, --user string                        Username or UID (format: <name|uid>[:<group|gid>])
      --with-registry-auth                 Send registry authentication details to swarm agents
  -w, --workdir string                     Working directory inside the container

```

| code           | tips                         |
| -------------- | ---------------------------- |
| docker run     | 容器启动，不具有扩缩容器     |
| docekr service | 服务，具有扩缩容器，滚动更新 |



1. nginx

### 11.4.1 创建服务

```bash
[root@localhost ~]# docker service create -p 8000:80 --name nginx_001 nginx
5t03xxxnrlg7m8hqp8421z6w5
overall progress: 1 out of 1 tasks 
1/1: running   
verify: Service converged
```

参数 --mode可以决定service 以什么方式运行

```bash
# 默认
docker service create --mode relicated -p 8000:80 --name nginx_001 nginx

# 日志收集
# 每一个节点都有自己的日志收集器，过滤，把所有的日志最终再传给日志中心，服务监控，状态性能
docker service create --mode global -p 8000:80 --name nginx_001 nginx
```





### 11.4.2 查看服务

```bash
[root@localhost ~]# docker service ps nginx_001
ID             NAME          IMAGE          NODE                    DESIRED STATE   CURRENT STATE           ERROR     PORTS
hvtgeznt34da   nginx_001.1   nginx:latest   localhost.localdomain   Running         Running 7 minutes ago             
```

### 11.4.3 创建副本（动态扩缩容）

开启n个副本

```bash
[root@localhost ~]# docker service update --replicas 3 nginx_001 
nginx_001
overall progress: 3 out of 3 tasks 
1/3: running   
2/3: running   
3/3: running   
verify: Service converged 
```

等同于

```bash
[root@localhost ~]# docker service scale nginx_001=10
nginx_001 scaled to 10
overall progress: 10 out of 10 tasks 
1/10: running   
2/10: running   
3/10: running   
4/10: running   
5/10: running   
6/10: running   
7/10: running   
8/10: running   
9/10: running   
10/10: running   
verify: Service converged 
[root@localhost ~]#
```

服务，集群中的任意节点都可以访问。服务可以有多个副本，动态扩缩容

### 11.4.4 移除服务

```bash
[root@localhost ~]# docker service rm nginx_001 
nginx_001
```

## 11.5 总结

1. swarm

集群的管理和编号，docker 可以初始化一个swarm 集群，其他节点可以加入（manager , worker）

2. Node

docker 节点，多个节点组成了一个网络集群

3. service 

任务，可以在管理节点或者工作节点来运行。核心，用户访问

4. Task

容器内的命令，细节任务















# 12.Docker Stack

| code           | tips         |
| -------------- | ------------ |
| docker-compose | 单机部署项目 |
| docker Stack   | 集群部署     |

```bash
# 单机
docker-compose up -d wordpress.ymal
# 集群
docker stack deploy wordpress.ymal
```



```bash
[root@localhost ~]# docker stack 

Usage:  docker stack [OPTIONS] COMMAND

Manage Docker stacks

Options:
      --orchestrator string   Orchestrator to use (swarm|kubernetes|all)

Commands:
  deploy      Deploy a new stack or update an existing stack
  ls          List stacks
  ps          List the tasks in the stack
  rm          Remove one or more stacks
  services    List the services in the stack

Run 'docker stack COMMAND --help' for more information on a command.
```











# 13.Docker Secret

安全，配置密码，证书

```bash
[root@localhost ~]# docker secret 

Usage:  docker secret COMMAND

Manage Docker secrets

Commands:
  create      Create a secret from a file or STDIN as content
  inspect     Display detailed information on one or more secrets
  ls          List secrets
  rm          Remove one or more secrets

Run 'docker secret COMMAND --help' for more information on a command.

```











# 14.Docker Config

给多个容器统一配置

```bash
[root@localhost ~]# docker config

Usage:  docker config COMMAND

Manage Docker configs

Commands:
  create      Create a config from a file or STDIN
  inspect     Display detailed information on one or more configs
  ls          List configs
  rm          Remove one or more configs

Run 'docker config COMMAND --help' for more information on a command.

```

