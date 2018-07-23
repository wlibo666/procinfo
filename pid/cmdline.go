// get info from /proc/$pid/cmdline
package pid

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/wlibo666/procinfo"
)

const (
	PROC_DIR     = "%s/%s"
	PROC_CMDLINE = "%s/%s/cmdline"
)

var (
	// map[pid]string
	ProcCmdlines   *sync.Map = &sync.Map{}
	disableCmdline           = false

	oldBytes []byte
	newBytes []byte
)

func DisableCmdLineMonitor() {
	disableCmdline = true
}

func GetPidCmdline(pid string) (string, error) {
	value, ok := ProcCmdlines.Load(pid)
	if ok {
		return value.(string), nil
	}
	path := fmt.Sprintf(PROC_CMDLINE, procinfo.PROC_BASE_DIR, pid)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	tmpByte := bytes.Replace(content, oldBytes, newBytes, -1)
	tmpStr := strings.Replace(string(tmpByte), "\"", "", -1)
	ProcCmdlines.Store(pid, tmpStr)
	return string(content), nil
}

// clean proccess's cmdline name every some seconds
func init() {
	oldBytes = append(oldBytes, byte(0x0))
	newBytes = append(newBytes, byte(0x20))
	go func() {
		for {
			if disableCmdline {
				break
			}
			time.Sleep(PID_CLEAN_INTERVAL * time.Second)
			ProcCmdlines.Range(func(key, value interface{}) bool {
				dir := fmt.Sprintf(PROC_DIR, procinfo.PROC_BASE_DIR, key.(string))
				_, err := os.Stat(dir)
				if err == os.ErrNotExist {
					ProcCmdlines.Delete(key)
				}
				return true
			})
		}
	}()
}
