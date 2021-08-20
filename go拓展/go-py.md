# 1 配置环境

## 1.1 Docker and Docker-Compose

```bash
# 安装
curl -fsSL https://get.docker.com | bash -s docker --mirror Aliyun

# 开机启动docker
systemctl enable docker

# 镜像加速
tee /etc/docker/daemon.json <<-'EOF'
{
	"registry-mirrors": ["https://8sqb5nwq.mirror.aliyuncs.com"] 
}
EOF

sudo systemctl daemon-reload 
sudo systemctl restart docker 
yum makecache fast


# Docker Compose
sudo curl -L "https://get.daocloud.io/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

sudo chmod +x /usr/local/bin/docker-compose
docker-compose --version
```





## 1.2 python

```bash
# 安装依赖包
# Centos
yum install openssl-devel bzip2-devel expat-devel gdbm-devel readline-devel sqlite-devel gcc gcc-c++ openssl-devel libffi-devel mariadb-devel

# 获取python
wget https://www.python.org/ftp/python/3.8.6/Python-3.8.6.tgz
tar -xvf Python-3.8.6.tgz -C /tmp
cd /tmp/Python-3.8.6

# 安装
./configure --prefix=/usr/local
make
make altinstall

# 更改/usr/bin/python链接
rm -f /usr/bin/python3 /usr/bin/pip3
ln -s /usr/local/bin/python3.8 /usr/bin/python3
ln -s /usr/local/bin/pip3.8 /usr/bin/pip3 

```

### 1.2.1 python虚拟环境

```bash
# windows 
pip install virtualenvwrapper-win

# linux
yum -y install python-setuptools python-devel
pip install virtualenvwrapper

# 导入环境 , virtualenvwrapper.sh根据系统的不同可能需要查询
vim ~/.bashrc

export VIRTUALENVWRAPPER_PYTHON=/usr/bin/python3
export WORKON_HOME=$HOME/.virtualenvs
source /usr/local/bin/virtualenvwrapper.sh
```

测试

```bash
mkvirtualenv python_test_1
```



## 1.3 配置nodejs

```bash
wget https://nodejs.org/dist/v14.17.3/node-v14.17.3-linux-x64.tar.xz
tar -xvf node-v14.17.3-linux-x64.tar.xz  -C /opt/

ln -s /opt/node-v14.17.3-linux-x64/bin/node  /usr/bin/node
ln -s /opt/node-v14.17.3-linux-x64/bin/npm  /usr/bin/npm
```









































































