package cgroups

import (
	"go-docker/cgroups/subsystem"

	"github.com/sirupsen/logrus"
)

//资源限制管理器

type CGroupManager struct {
	Path string
}

func NewCGroupManager(path string) *CGroupManager {
	return &CGroupManager{Path: path}
}

func (c *CGroupManager) Set(res *subsystem.ResourceConfig) {
	for _, subsystem := range subsystem.Subsystems {
		err := subsystem.Set(c.Path, res)
		if err != nil {
			logrus.Errorf("set %s err: %v", subsystem.Name(), err)
		}
	}
}

func (c *CGroupManager) Apply(pid int) {
	for _, subsystem := range subsystem.Subsystems {
		err := subsystem.Apply(c.Path, pid)
		if err != nil {
			logrus.Errorf("apply task,  err: %v", err)
		}
	}
}

func (c *CGroupManager) Destroy() {
	for _, subsystem := range subsystem.Subsystems {
		err := subsystem.Remove(c.Path)
		if err != nil {
			logrus.Errorf("remove %s task,  err: %v", subsystem.Name(), err)
		}
	}
}

//案例：memory.go
