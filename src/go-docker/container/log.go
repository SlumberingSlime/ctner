package container

import (
	"fmt"
	"go-docker/common"
	"io/ioutil"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

func LookContainerLog(containerName string) {
	logFileName := path.Join(common.DefaultContainerInfoPath, containerName, common.ContainerLogFileName)
	file, err := os.Open(logFileName)
	if err != nil {
		logrus.Errorf("open log file, path: %s, error: %v", logFileName, err)
	}
	bs, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Errorf("read log file, error: %v", err)
	}
	_, _ = fmt.Fprint(os.Stdout, string(bs))
}
