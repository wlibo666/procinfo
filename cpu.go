// get info from /proc/stat
// http://man7.org/linux/man-pages/man5/proc.5.html
package procinfo

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	// get cpu usage from PROC_CPU
	PROC_CPU = "/proc/stat"
	// get cpu number from PROC_CPUINFO
	PROC_CPUINFO = "/proc/cpuinfo"
)

type CpuUsage struct {
	User    uint64
	Nice    uint64
	System  uint64
	Idle    uint64
	Iowait  uint64
	Irq     uint64
	SoftIrq uint64
	St      uint64
}

type CpuUsageRate struct {
	Us float64
	Sy float64
	Ni float64
	Id float64
	Wa float64
	Hi float64
	Si float64
	St float64
}

var (
	CpuCnt = 1
	// sampling frequency
	CpuInterval                 = 1
	curSysCpuInfo *CpuUsage     = &CpuUsage{}
	preSysCpuInfo *CpuUsage     = &CpuUsage{}
	cpuRwLock     *sync.RWMutex = &sync.RWMutex{}
	disableSysCpu               = false
)

func DisableSysCpuMonitor() {
	disableSysCpu = true
}

func getCpuCnt() (int, error) {
	content, err := ioutil.ReadFile(PROC_CPUINFO)
	if err != nil {
		return 0, err
	}
	return strings.Count(string(content), "processor"), nil
}

func getCpuUsage() (*CpuUsage, error) {
	fp, err := os.OpenFile(PROC_CPU, os.O_RDONLY, 0444)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(fp)
	line, err := reader.ReadString('\n')
	fp.Close()
	if err != nil {
		return nil, err
	}
	cpuinfo := &CpuUsage{}
	fmt.Sscanf(string(line), "cpu %d %d %d %d %d %d %d %*s",
		&cpuinfo.User, &cpuinfo.Nice, &cpuinfo.System, &cpuinfo.Idle, &cpuinfo.Iowait,
		&cpuinfo.Irq, &cpuinfo.SoftIrq, &cpuinfo.St)
	cpuRwLock.Lock()
	preSysCpuInfo = curSysCpuInfo
	curSysCpuInfo = cpuinfo
	cpuRwLock.Unlock()
	return cpuinfo, nil
}

func GetCpuUsage() *CpuUsage {
	cpuUsage := &CpuUsage{}
	cpuRwLock.RLock()
	cpuUsage.User = (curSysCpuInfo.User - preSysCpuInfo.User)
	cpuUsage.System = (curSysCpuInfo.System - preSysCpuInfo.System)
	cpuUsage.Nice = (curSysCpuInfo.Nice - preSysCpuInfo.Nice)
	cpuUsage.Idle = (curSysCpuInfo.Idle - preSysCpuInfo.Idle)
	cpuUsage.Iowait = (curSysCpuInfo.Iowait - preSysCpuInfo.Iowait)
	cpuUsage.Irq = (curSysCpuInfo.Irq - preSysCpuInfo.Irq)
	cpuUsage.SoftIrq = (curSysCpuInfo.SoftIrq - preSysCpuInfo.SoftIrq)
	cpuUsage.St = (curSysCpuInfo.St - preSysCpuInfo.St)
	cpuRwLock.RUnlock()
	return cpuUsage
}

func GetCpuUsageRate() *CpuUsageRate {
	cpuUsage := &CpuUsage{}
	cpuRwLock.RLock()
	cpuUsage.User = (curSysCpuInfo.User - preSysCpuInfo.User)
	cpuUsage.System = (curSysCpuInfo.System - preSysCpuInfo.System)
	cpuUsage.Nice = (curSysCpuInfo.Nice - preSysCpuInfo.Nice)
	cpuUsage.Idle = (curSysCpuInfo.Idle - preSysCpuInfo.Idle)
	cpuUsage.Iowait = (curSysCpuInfo.Iowait - preSysCpuInfo.Iowait)
	cpuUsage.Irq = (curSysCpuInfo.Irq - preSysCpuInfo.Irq)
	cpuUsage.SoftIrq = (curSysCpuInfo.SoftIrq - preSysCpuInfo.SoftIrq)
	cpuUsage.St = (curSysCpuInfo.St - preSysCpuInfo.St)
	cpuRwLock.RUnlock()

	cpuTotal := cpuUsage.User + cpuUsage.System + cpuUsage.Nice + cpuUsage.Idle +
		cpuUsage.Iowait + cpuUsage.Irq + cpuUsage.SoftIrq + cpuUsage.St
	cur := &CpuUsageRate{}
	if cpuTotal == 0 {
		return cur
	}
	//fmt.Printf("%d %d %d %d %d %d %d %d\n", cpuUsage.User, cpuUsage.Nice, cpuUsage.System, cpuUsage.Idle,
	//	cpuUsage.Iowait, cpuUsage.Irq, cpuUsage.SoftIrq, cpuUsage.St)

	cur.Us = float64(cpuUsage.User) / float64(cpuTotal)
	cur.Sy = float64(cpuUsage.System) / float64(cpuTotal)
	cur.Ni = float64(cpuUsage.Nice) / float64(cpuTotal)
	cur.Id = float64(cpuUsage.Idle) / float64(cpuTotal)
	cur.Wa = float64(cpuUsage.Iowait) / float64(cpuTotal)
	cur.Hi = float64(cpuUsage.Irq) / float64(cpuTotal)
	cur.Si = float64(cpuUsage.SoftIrq) / float64(cpuTotal)
	cur.St = float64(cpuUsage.St) / float64(cpuTotal)
	return cur
}

func (cur *CpuUsageRate) String() string {
	return fmt.Sprintf("%%Cpu(s): %.1f us,  %.1f sy,  %.1f ni, %.1f id,  %.1f wa,  %.1f hi,  %.1f si,  %.1f st",
		cur.Us*100, cur.Sy*100, cur.Ni*100, cur.Id*100, cur.Wa*100, cur.Hi*100, cur.Si*100, cur.St*100)
}

func init() {
	cnt, err := getCpuCnt()
	if err != nil {
		panic("getCpuCnt failed,err:" + err.Error())
	}
	CpuCnt = cnt
	go func() {
		for {
			if disableSysCpu {
				break
			}
			getCpuUsage()
			time.Sleep(time.Duration(CpuInterval) * time.Second)
		}
	}()
}
