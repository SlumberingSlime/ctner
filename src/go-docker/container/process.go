package container

import (
	"go-docker/common"
	"os"
	"os/exec"
	"path"
	"syscall"

	"github.com/sirupsen/logrus"
)

// 创建一个会隔离namespace进程的Command
func NewParentProcess(tty bool, volume, containerName, imageName string, envs []string) (*exec.Cmd, *os.File) {
	//*exec.Cmd, *os.File是这个函数的返回值。Cmd代表一个正在准备或者在执行中的外部命令。
	//File代表一个打开的文件对象

	readPipe, writePipe, _ := os.Pipe()

	//调用自身，传入init参数，执行initCommannd命令， 设置相应隔离信息
	//initCommand内容在init.go里面
	cmd := exec.Command("/proc/self/exe", "init")

	//cmd.SysProcAttr保管可选的、各操作系统特定的sys执行属性
	cmd.SysProcAttr = &syscall.SysProcAttr{
		//Cloneflags: Flags for clone calls (Linux only)
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC | //= 0x8000000
			syscall.CLONE_NEWPID | //= 0x20000000
			syscall.CLONE_NEWNS | //= 0x20000
			syscall.CLONE_NEWUSER | //= 0x10000000
			syscall.CLONE_NEWNET, //= 0x40000000
	}
	if tty {
		//指定进程的标准输入、输出、错误
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		//日志输出到文件里
		logDir := path.Join(common.DefaultContainerInfoPath, containerName)
		if _, err := os.Stat(logDir); err != nil && os.IsNotExist(err) {
			err := os.Mkdir(logDir, os.ModePerm)
			if err != nil {
				logrus.Errorf("mkdir container log, err: %v", err)
			}
		}
		logFileName := path.Join(logDir, common.ContainerLogFileName)
		file, err := os.Create(logFileName)
		if err != nil {
			logrus.Errorf("create log file, error: %v", err)
		}

		cmd.Stdout = file //cmd的输出流改到流文件
	}

	//设置额外文件句柄
	cmd.ExtraFiles = []*os.File{ //ExtraFiles指定额外被新进程继承的已打开文件流，不包括标准输入、标准输出、标准错误输出
		readPipe,
	}

	// 设置环境变量
	cmd.Env = append(os.Environ(), envs...)
	err := NewWorkSpace(volume, containerName, imageName)
	if err != nil {
		logrus.Errorf("new work space, err: %v", err)
	}

	// 指定容器初始化后的工作目录
	cmd.Dir = common.MntPath

	return cmd, writePipe
}
