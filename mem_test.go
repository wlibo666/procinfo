package procinfo

import (
	"testing"
	"time"
)

func TestGetMemInfo(t *testing.T) {
	for i := 0; i < 10; i++ {
		mem := GetMemInfo()
		t.Logf("index:%d,meminfo:%v", i, mem)
		time.Sleep(time.Second)
	}
}
