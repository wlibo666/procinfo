package pid

import (
	"testing"
)

func TestGetProcMemInfo(t *testing.T) {
	info, err := GetProcMemInfo("1")
	if err != nil {
		t.Fatalf("get failed,err:%s", err.Error())
	}
	t.Logf("mem info,vmsize:%d KB,vmrss:%d KB", info.VmSize, info.VmRss)
}
