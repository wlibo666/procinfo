package pid

import (
	"testing"
)

func TestGetPidCmdline(t *testing.T) {
	cmdline, err := GetPidCmdline("1")
	if err != nil {
		t.Fatalf("get failed,err:%s", err.Error())
	}
	t.Logf("pid 1,cmdline:%s", cmdline)
}
