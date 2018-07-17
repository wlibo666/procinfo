package procinfo

import (
	"testing"
	"time"
)

func TestGetHostName(t *testing.T) {
	time.Sleep(time.Second)
	t.Logf("host:%v", Host)
}
