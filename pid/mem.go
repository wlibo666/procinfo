// get info from /proc/$pid/status
package pid

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	PROC_PID_MEM = "/proc/%s/status"
)

type ProcMemInfo struct {
	VmSize   uint64
	VmRss    uint64
	LastLive int64
}

var (
	MemInterval = 1
	// map[pid]*ProcMemInfo
	ProcMemInfos  *sync.Map = &sync.Map{}
	disablePidMem           = false
)

func DisablePidMemMonitor() {
	disablePidMem = true
}

func getPidMemInfo(pid string) (*ProcMemInfo, error) {
	path := fmt.Sprintf(PROC_PID_MEM, pid)
	fp, err := os.OpenFile(path, os.O_RDONLY, 0444)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	info, ok := ProcMemInfos.Load(pid)
	if !ok {
		info = &ProcMemInfo{LastLive: time.Now().Unix()}
		ProcMemInfos.Store(pid, info)
	}

	reader := bufio.NewReader(fp)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if strings.HasPrefix(line, "VmSize:") {
			fmt.Sscanf(line, "VmSize: %d kB", &info.(*ProcMemInfo).VmSize)
		} else if strings.HasPrefix(line, "VmRSS:") {
			fmt.Sscanf(line, "VmRSS: %d kB", &info.(*ProcMemInfo).VmRss)
			info.(*ProcMemInfo).LastLive = time.Now().Unix()
			break
		}
	}
	return info.(*ProcMemInfo), nil
}

func GetProcMemInfo(pid string) (*ProcMemInfo, error) {
	info, ok := ProcMemInfos.Load(pid)
	if ok {
		return info.(*ProcMemInfo), nil
	}
	return getPidMemInfo(pid)
}

func init() {
	go func() {
		for {
			if disablePidMem {
				break
			}
			ProcMemInfos.Range(func(key, value interface{}) bool {
				getPidMemInfo(key.(string))
				return true
			})
			time.Sleep(time.Duration(MemInterval) * time.Second)
			ProcMemInfos.Range(func(key, value interface{}) bool {
				dir := fmt.Sprintf(PROC_DIR, key.(string))
				_, err := os.Stat(dir)
				if err == os.ErrNotExist {
					ProcMemInfos.Delete(key)
				}
				return true
			})
		}
	}()
}
