package pid

import (
	"fmt"
	"testing"
)

func TestGetProcCpuUsage(t *testing.T) {
	var tmpS string
	content := "1 (systemd) S 0 1 1 0 -1 4202752 124272 908123 58 1836 778 1855 30372 109457 20 0 1 0 19 129167360 1150 18446744073709551615 94754739941376 94754741379833 140724869776864 140724869771840 139855373312387 0 671173123 4096 1260 18446744072365702430 0 0 17 0 0 0 28 0 0 94754743477336 94754743621176 94754761469952 140724869779356 140724869779423 140724869779423 140724869779423 0"
	tmpInfo := &ProcCpuInfo{}
	n, err := fmt.Sscanf(content, "%s %s %s %s %s %s %s %s %s %s %s %s %s %d %d %d %d", &tmpS, &tmpS, &tmpS, &tmpS, &tmpS, &tmpS, &tmpS, &tmpS, &tmpS, &tmpS, &tmpS, &tmpS, &tmpS,
		&tmpInfo.Utime, &tmpInfo.Stime, &tmpInfo.Cutime, &tmpInfo.Cstime)
	if err != nil {
		t.Fatalf("scan failed,n:%d,err:%s", n, err.Error())
	}
	t.Logf("tmpInfo:%v", tmpInfo)
	if (tmpInfo.Utime != uint64(778)) || (tmpInfo.Stime != uint64(1855)) || (tmpInfo.Cutime != uint64(30372)) || (tmpInfo.Cstime != uint64(109457)) {
		t.Fatalf("scan failed")
	}
}
