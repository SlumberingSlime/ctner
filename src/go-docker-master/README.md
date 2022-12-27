# go-docker
> 用go写一个docker

## 环境配置
- 开发环境： windows
- 运行环境： CentOS

### windows中goland配置
> windows下要修改goland的OS环境为 linux,不然只会引用`exec_windows.go`而不会引用`exec_linxu_go`
> 在Setting->Go->Build Tags & Vendoring -> OS=linux

### 设置CentOS支持aufs
查看是否支持
```bash
cat /proc/filesystems
```
安装aufs
```bash
cd /etc/yum.repo.d
# 下载文件
wget https://yum.spaceduck.org/kernel-ml-aufs/kernel-ml-aufs.repo
# 安装
yum install kernel-ml-aufs
# 修改内核启动
vim /etc/default/grub
## 修改参数
GRUB_DEFAULT=0

# 重新生成grub.cfg
grub2-mkconfig -o /boot/grub2/grub.cfg

# 重启计算机
reboot
```
### 配置busybox
```bash
# 下载 busybox
docker pull busybox
# 运行
docker run -d busybox top -b
# 导出
docker export -o busybox.tar (容器ID)
# 解压到 /root文件夹下
cd /root
mkdir busybox
tar -xvf busybox.tar -C busybox/
```

## 使用指南
```bash
# 编译
go build .

# 启动一个容器, busybox为镜像名，存放在 /root/busybox.tar
./go-docker run -ti --name test busybox sh

# 后台启动
./go-docker run -d --name test busybox sh

# 挂载文件
./go-docker run -d -v /root/test:/test --name test busybox sh

# 进入容器
./go-docker exec test sh

# 查看容器日志
./go-docker logs test

# 查看容器列表
./go-docker ps

# 停止容器
./go-docker stop test

# 删除容器
./go-docker rm test
```

## docker核心技术
### namespace
- uts : 隔离主机名
- pid : 隔离进程pid
- user : 隔离用户
- network : 隔离网络
- mount : 隔离挂载点
- ipc : 隔离System VIPC和POSIX message queues

### cgroup
> 主要是使用三个组件相互协作实现的，分别是：subsystem, hierarchy, cgroup,

- cgroup: 是对进程分组管理的一种机制
- subsystem: 是一组资源控制的模块
- hierarchy: 把一组cgroup串成一个树状结构(可让其实现继承)

#### 实现方式
> 主要实现方式是在`/sys/fs/cgroup/` 文件夹下，根据限制的不同，创建一个新的文件夹即可，kernel会将这个文件夹
> 标记为它的`子cgroup`, 比如要限制内存使用，则在`/sys/fs/cgroup/memory/` 下创建`test-limit-memory`文件夹即可，将
> 内存限制数写到该文件夹里面的 `memory.limit_in_bytes`即可

## 指令小记

- 查看Linux程序父进程
```bash
pstree -pl | grep main
```
- 查看进程id
```bash
echo $$
```
- 查看进程的uts
```bash
readling /proc/进程id/ns/uts
```
- 修改hostname
```bash
hostname -b 新名称
```
- 常看当前用户和用户组
```bash
id
```
- 创建并挂载一个hierarchy
> 在这个文件夹下面创建新的文件夹，会被kernel标记为该`cgroup`的子`cgroup`
```bash
mkdir cgroup-test
mount -t cgroup -o none,name=cgroup-test cgroup-test ./cgroup-test
```
- 将其他进程移动到其他的`cgroup`中
> 只要将该进程的ID放到其`cgroup`的`tasks`里面即可
```bash
echo "进程ID" >> cgroup/tasks 
```

- 导出容器
```bash
docker export -o busybox.tar 45c98e055883(容器ID)
```
- 移除mount
```bash
unshare -m
```