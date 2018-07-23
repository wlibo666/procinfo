// get info from /proc/meminfo
package procinfo

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	PROC_MEMINFO = "%s/meminfo"
)

var (
	SysMemInfo    *MemInfo      = &MemInfo{}
	MemInterval                 = 3
	memRwLock     *sync.RWMutex = &sync.RWMutex{}
	disableSysMem               = false
)

type MemInfo struct {
	Total     uint64
	Free      uint64
	Available uint64
}

func DisableSysMem() {
	disableSysMem = true
}

func getMemInfo() error {
	fp, err := os.OpenFile(fmt.Sprintf(PROC_MEMINFO, PROC_BASE_DIR), os.O_RDONLY, 0444)
	if err != nil {
		return err
	}
	defer fp.Close()

	reader := bufio.NewReader(fp)
	memRwLock.Lock()
	defer memRwLock.Unlock()
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if strings.HasPrefix(line, "MemTotal:") {
			fmt.Sscanf(line, "MemTotal: %d kB", &SysMemInfo.Total)
		} else if strings.HasPrefix(line, "MemFree:") {
			fmt.Sscanf(line, "MemFree: %d kB", &SysMemInfo.Free)
		} else if strings.HasPrefix(line, "MemAvailable:") {
			fmt.Sscanf(line, "MemAvailable: %d kB", &SysMemInfo.Available)
			break
		}
	}
	return nil
}

func GetMemInfo() MemInfo {
	if SysMemInfo.Total == 0 {
		getMemInfo()
	}
	memRwLock.RLock()
	defer memRwLock.RUnlock()
	return MemInfo{
		Total:     SysMemInfo.Total,
		Free:      SysMemInfo.Free,
		Available: SysMemInfo.Available,
	}
}

func init() {
	go func() {
		for {
			if disableSysMem {
				break
			}
			getMemInfo()
			time.Sleep(time.Duration(MemInterval) * time.Second)
		}
	}()
}
