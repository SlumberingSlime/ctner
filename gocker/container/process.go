package container

import (
	"os"
	"os/exec"
	"path"
	"syscall"

	"github.com/sirupsen/logrus"

	"gocker/common"
)

// 创建一个会隔离namespace进程的Command
func NewParentProcess(tty bool, volume, containerName, imageName string, envs []string) (*exec.Cmd, *os.File) {
	//Cmd代表一个正在准备或者在执行中的外部命令, File代表一个打开的文件对象

	readPipe, writePipe, _ := os.Pipe()
	// 调用自身，传入 init 参数，也就是执行 initCommand
	//initCommand内容在init.go里面
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
		//Cloneflags: Flags for clone calls (Linux only)
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		// 把日志输出到文件里
		logDir := path.Join(common.DefaultContainerInfoPath, containerName)
		if _, err := os.Stat(logDir); err != nil && os.IsNotExist(err) {
			err := os.MkdirAll(logDir, os.ModePerm)
			if err != nil {
				logrus.Errorf("mkdir container log, err: %v", err)
			}
		}
		logFileName := path.Join(logDir, common.ContainerLogFileName)
		file, err := os.Create(logFileName)
		if err != nil {
			logrus.Errorf("create log file, err: %v", err)
		}
		cmd.Stdout = file
	}
	// 设置额外文件句柄
	cmd.ExtraFiles = []*os.File{
		readPipe,
	}
	// 设置环境变量
	cmd.Env = append(os.Environ(), envs...)
	err := NewWorkSpace(volume, containerName, imageName)
	if err != nil {
		logrus.Errorf("new work space, err: %v", err)
	}
	// 指定容器初始化后的工作目录
	cmd.Dir = path.Join(common.MntPath, containerName)
	return cmd, writePipe
}
