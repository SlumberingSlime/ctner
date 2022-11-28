package main

import (
	"os" //提供了操作系统函数的不依赖平台的接口，在所有操作系统中都是一致的。非公用的属性可以从操作系统特定的syscall包获取

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const usage = `go-docker`

func main() {
	app := cli.NewApp()
	app.Name = "go-docker"
	app.Usage = usage

	app.Commands = []*cli.Command{
		&runCommand,
		&initCommand,
		&commitCommand,
		&listCommand,
		&logCommand,
		&execCommand,
		&stopCommand,
		&removeCommand,
		&networkCommand,
	}
	//Command数组定义了两个运行命令runCommand和InitCommand，在command.go里面
	app.Before = func(contest *cli.Context) error {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(os.Stdout) //os.Stdout是指向标准输出的文件描述符
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		//os.Args是保管命令行参数的切片（动态数组
		logrus.Fatal(err)
	}
}
