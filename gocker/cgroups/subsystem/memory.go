package subsystem

import (
	"os"
	"path"    //实现了对斜杠分隔的路径的实用操作函数
	"strconv" //实现了基本数据类型和其字符串表示的相互转换

	"github.com/sirupsen/logrus"
)

type MemorySubSystem struct {
	apply bool
}

func (*MemorySubSystem) Name() string {
	return "memory"
}

func (m *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	subsystemCgroupPath, err := GetCgroupPath(m.Name(), cgroupPath, true)
	if err != nil {
		logrus.Errorf("get %s path, err: %v", cgroupPath, err)
		return err
	}

	//os.WriteFile: 向第一个参数指定的文件中写入[]byte数据。如果文件不存在将按给出的权限创建文件，否则在写入数据之前清空文件
	//path.Join: 将任意数量的路径元素放入一个单一路径里，会根据需要添加斜杠。结果是经过简化的，所有的空字符串元素会被忽略
	if res.MemoryLimit != "" {
		m.apply = true
		// 设置cgroup内存限制，
		// 将这个限制写入到cgroup对应目录的 memory.limit_in_bytes文件中即可
		err := os.WriteFile(path.Join(subsystemCgroupPath, "memory.limit_in_bytes"), []byte(res.MemoryLimit), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MemorySubSystem) Remove(cgroupPath string) error {
	subsystemCgroupPath, err := GetCgroupPath(m.Name(), cgroupPath, false)
	if err != nil {
		return err
	}
	return os.RemoveAll(subsystemCgroupPath)
}

func (m *MemorySubSystem) Apply(cgroupPath string, pid int) error {
	if m.apply {
		subsystemCgroupPath, err := GetCgroupPath(m.Name(), cgroupPath, false)
		if err != nil {
			return err
		}
		tasksPath := path.Join(subsystemCgroupPath, "tasks")
		err = os.WriteFile(tasksPath, []byte(strconv.Itoa(pid)), os.ModePerm)
		if err != nil {
			logrus.Errorf("write pid to tasks, path: %s, pid: %d, err: %v", tasksPath, pid, err)
			return err
		}
	}
	return nil
}
