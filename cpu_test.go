package procinfo

import (
	"bufio"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestGetCpuUsageRate(t *testing.T) {
	cnt, err := getCpuCnt()
	if err != nil {
		t.Fatalf("get cpu cnt failed,err:%s", err.Error())
	}
	fmt.Fprintf(os.Stdout, "cpu cnt:%d\n", cnt)
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
		cur := GetCpuUsageRate()
		fmt.Fprintf(os.Stdout, "%s\n", cur)
	}
}

func TestGetCpuInfo(t *testing.T) {
	fp, err := os.OpenFile(fmt.Sprintf(PROC_CPU, PROC_BASE_DIR), os.O_RDONLY, 0444)
	if err != nil {
		t.Fatalf("open failed,err:%s\n", err.Error())
	}
	reader := bufio.NewReader(fp)
	line, err := reader.ReadString('\n')
	fp.Close()
	if err != nil {
		t.Fatalf("read failed,err:%s\n", err.Error())
	}
	fmt.Fprintf(os.Stdout, "line:%s\n", string(line))
	cpuinfo := &CpuUsage{}
	fmt.Sscanf(string(line), "cpu %d %d %d %d %d %d %d %*s",
		&cpuinfo.User, &cpuinfo.Nice, &cpuinfo.System, &cpuinfo.Idle, &cpuinfo.Iowait,
		&cpuinfo.Irq, &cpuinfo.SoftIrq, &cpuinfo.St)
	fmt.Fprintf(os.Stdout, "cpuinfo:%v\n", cpuinfo)
}
