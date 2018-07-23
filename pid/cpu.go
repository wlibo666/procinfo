// get info from /proc/$pid/stat
package pid

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/wlibo666/procinfo"
)

const (
	PID_CLEAN_INTERVAL = 60

	PID_PROC_PATH = "%s/%s/stat"
)

type ProcCpuInfo struct {
	Utime  uint64
	Stime  uint64
	Cutime uint64
	Cstime uint64
}

type ProcCpuInfos struct {
	Cur      *ProcCpuInfo
	Pre      *ProcCpuInfo
	LastLive int64
}

type ProcCpuUsageRate struct {
	Rate float64
}

var (
	// map[pid]*ProcCpuInfos
	PidCpuInfo    *sync.Map = &sync.Map{}
	disablePidCpu           = false
)

func DisablePidCpuMonitor() {
	disablePidCpu = true
	procinfo.DisableSysCpuMonitor()
}

func MoniPidCpu(pid string) {
	_, ok := PidCpuInfo.Load(pid)
	if !ok {
		infos := &ProcCpuInfos{
			Cur:      &ProcCpuInfo{},
			Pre:      &ProcCpuInfo{},
			LastLive: time.Now().Unix(),
		}
		PidCpuInfo.Store(pid, infos)
	}
}

func getProcCpuUsage(pid string) error {
	path := fmt.Sprintf(PID_PROC_PATH, procinfo.PROC_BASE_DIR, pid)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	tmpInfo := &ProcCpuInfo{}
	var tmpS string
	fmt.Sscanf(string(content), "%s %s %s %s %s %s %s %s %s %s %s %s %s %d %d %d %d",
		&tmpS, &tmpS, &tmpS, &tmpS, &tmpS, &tmpS, &tmpS, &tmpS, &tmpS, &tmpS, &tmpS, &tmpS, &tmpS,
		&tmpInfo.Utime, &tmpInfo.Stime, &tmpInfo.Cutime, &tmpInfo.Cstime)

	value, ok := PidCpuInfo.Load(pid)
	if !ok {
		value = &ProcCpuInfos{
			Cur: &ProcCpuInfo{
				Utime:  tmpInfo.Utime,
				Stime:  tmpInfo.Stime,
				Cutime: tmpInfo.Cutime,
				Cstime: tmpInfo.Cstime,
			},
			Pre:      &ProcCpuInfo{},
			LastLive: time.Now().Unix(),
		}
		PidCpuInfo.Store(pid, value)
	} else {
		value.(*ProcCpuInfos).Pre.Utime = value.(*ProcCpuInfos).Cur.Utime
		value.(*ProcCpuInfos).Pre.Stime = value.(*ProcCpuInfos).Cur.Stime
		value.(*ProcCpuInfos).Pre.Cutime = value.(*ProcCpuInfos).Cur.Cutime
		value.(*ProcCpuInfos).Pre.Cstime = value.(*ProcCpuInfos).Cur.Cstime

		value.(*ProcCpuInfos).Cur.Utime = tmpInfo.Utime
		value.(*ProcCpuInfos).Cur.Stime = tmpInfo.Stime
		value.(*ProcCpuInfos).Cur.Cutime = tmpInfo.Cutime
		value.(*ProcCpuInfos).Cur.Cstime = tmpInfo.Cstime

		value.(*ProcCpuInfos).LastLive = time.Now().Unix()
	}

	return nil
}

func GetProcCpuUsage(pid string) *ProcCpuInfo {
	value, ok := PidCpuInfo.Load(pid)
	if !ok {
		return nil
	}
	info := &ProcCpuInfo{}
	info.Utime = (value.(*ProcCpuInfos).Cur.Utime - value.(*ProcCpuInfos).Pre.Utime)
	info.Stime = (value.(*ProcCpuInfos).Cur.Stime - value.(*ProcCpuInfos).Pre.Stime)
	info.Cutime = (value.(*ProcCpuInfos).Cur.Cutime - value.(*ProcCpuInfos).Pre.Cutime)
	info.Cstime = (value.(*ProcCpuInfos).Cur.Cstime - value.(*ProcCpuInfos).Pre.Cstime)
	return info
}

func GetProcCpuUsageRate(pid string) (*ProcCpuUsageRate, error) {
	procInfo := GetProcCpuUsage(pid)
	if procInfo == nil {
		return nil, fmt.Errorf("not found info by pid:%s", pid)
	}
	allInfo := procinfo.GetCpuUsage()
	cpuTotal := allInfo.User + allInfo.System + allInfo.Nice + allInfo.Idle +
		allInfo.Iowait + allInfo.Irq + allInfo.SoftIrq + allInfo.St
	rate := (float64(procInfo.Utime+procInfo.Stime+procInfo.Cutime+procInfo.Cstime) / float64(cpuTotal)) * float64(procinfo.CpuCnt)
	return &ProcCpuUsageRate{
		Rate: rate,
	}, nil
}

func GetAllProcCpuUsageRate() map[string]*ProcCpuUsageRate {
	cpuRate := make(map[string]*ProcCpuUsageRate)

	PidCpuInfo.Range(func(key, value interface{}) bool {
		pid := key.(string)
		rate, err := GetProcCpuUsageRate(pid)
		if err != nil {
			return false
		}
		cpuRate[pid] = &ProcCpuUsageRate{Rate: rate.Rate}
		return true
	})
	return cpuRate
}

func init() {
	go func() {
		for {
			if disablePidCpu {
				break
			}
			PidCpuInfo.Range(func(key, value interface{}) bool {
				getProcCpuUsage(key.(string))
				return true
			})
			time.Sleep(time.Duration(procinfo.CpuInterval) * time.Second)
			PidCpuInfo.Range(func(key, value interface{}) bool {
				dir := fmt.Sprintf(PROC_DIR, key.(string))
				_, err := os.Stat(dir)
				if err == os.ErrNotExist {
					PidCpuInfo.Delete(key)
				}
				return true
			})
		}
	}()
}
