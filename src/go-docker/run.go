package main

import (
	"os"
	"strconv"
	"strings" //有数据结构用到string了，需要引入这个包

	"github.com/sirupsen/logrus"

	"go-docker/cgroups"
	"go-docker/cgroups/subsystem"
	"go-docker/container"
	"go-docker/network"
)

// 首字母大写的函数/变量可以被其他package引用（公有），小写表示本包私有
// 启动一个容器，然后进行相应的资源限制
func Run(cmdArray []string, tty bool, res *subsystem.ResourceConfig, containerName, imageName, volume, net string, envs, ports []string) { //这个函数没有返回值
	containerID := container.GenContainerID(10)
	if containerName == "" {
		containerName = containerID
	}

	parent, writePipe := container.NewParentProcess(tty, volume, containerName, imageName, envs)
	//返回一个被namespace隔离的进程，函数在文件process.go里面
	if parent == nil {
		logrus.Errorf("failed to new parent process")
		return
	}
	if err := parent.Start(); err != nil {
		logrus.Errorf("parent start failed, err: %v", err)
		return
	}

	// 记录容器信息
	err := container.RecordContainerInfo(parent.Process.Pid, cmdArray, containerName, containerID)
	if err != nil {
		logrus.Errorf("record container info, err: %v", err)
	}

	//添加资源限制
	cgroupManager := cgroups.NewCGroupManager("go-docker")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res) //设置资源限制
	//容器进程加入到subsystem挂载对应的cgroup中
	cgroupManager.Apply(parent.Process.Pid)

	// 设置网络
	if net != "" {
		// 初始化容器网络
		err = network.Init()
		if err != nil {
			logrus.Errorf("network init failed, err: %v", err)
			return
		}
		containerInfo := &container.ContainerInfo{
			Id:          containerID,
			Pid:         strconv.Itoa(parent.Process.Pid),
			Name:        containerName,
			PortMapping: ports,
		}
		if err := network.Connect(net, containerInfo); err != nil {
			logrus.Errorf("connect network, err: %v", err)
			return
		}
	}

	// 设置初始化命令
	sendInitCommand(cmdArray, writePipe)
	if tty {
		// 等待父进程结束
		err := parent.Wait()
		if err != nil {
			logrus.Errorf("parent wait, err: %v", err)
		}
		// 删除容器工作空间
		err = container.DeleteWorkSpace(containerName, volume)
		if err != nil {
			logrus.Errorf("delete work space, err: %v", err)
		}
		// 删除容器信息
		container.DeleteContainerInfo(containerName)
	}
}

func sendInitCommand(cmdArray []string, writePipe *os.File) {
	command := strings.Join(cmdArray, " ")
	logrus.Infof("command all is %s", command)
	_, _ = writePipe.WriteString(command) //_表示忽略返回值
	_ = writePipe.Close()
}
