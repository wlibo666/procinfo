package procinfo

import (
	"os"
)

type HostInfo struct {
	HostName string `json:"host_name"`
	CpuCnt   int    `json:"host_cpu"`
	Memory   uint64 `json:"host_memory"`
}

var (
	Host HostInfo
)

func GetHostName() (string, error) {
	return os.Hostname()
}

func InitHostInfo() {
	cnt, err := getCpuCnt()
	if err == nil {
		Host.CpuCnt = cnt
	}
	name, err := GetHostName()
	if err == nil {
		Host.HostName = name
	}
	mem := GetMemInfo()
	Host.Memory = mem.Total
}

func init() {
	InitHostInfo()
}
