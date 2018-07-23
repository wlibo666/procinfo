package procinfo

var (
	PROC_BASE_DIR = "/proc"
)

func SetProcBaseDir(baseDir string) {
	PROC_BASE_DIR = baseDir
}
