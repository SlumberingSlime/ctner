package subsystem

type ResourceConfig struct { //资源配置限制
	MemoryLimit string
	CpuShare    string //cpu时间片权重
	CpuSet      string //CPU核数
}

// 将cgroup抽象成path, 因为在hierarchy中，cgroup便是虚拟的路径
type Subsystem interface {
	Name() string //返回subsystem名字，如CPU， Memory
	Set(cgroupPath string, res *ResourceConfig) error
	//设置crgoup在这个subsystem的资源限制

	Remove(cgroupPath string) error         //撤销资源限制
	Apply(cgroupPath string, pid int) error //把一个进程加入cgroup
}

var (
	Subsystems = []Subsystem{
		&MemorySubSystem{},
		&CpuSubSystem{},
		&CpuSetSubSystem{},
	}
)
