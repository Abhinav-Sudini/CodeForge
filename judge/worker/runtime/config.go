package runtime

type Runtime_conf struct {
	Runtime       string
	CompileComand []string
	RunComand     []string
	FileExtention string
	CodeFileName  string
}

const Is_production = true

const BinaryFileName = "out.bin"
const StdErrorFileName = "std.err"
const StdOutFileName = "std.out"
const VerdictFileName = "verdict.json"

var ExecWorkerStartComand = "/codeForge/exec_worker.bin"

var allRuntimes = map[string]Runtime_conf{

	"c++17": {
		Runtime:       "c++17",
		CompileComand: []string{"g++", "-o", BinaryFileName, "Code.cpp"},
		RunComand:     []string{"./out.bin"},
		CodeFileName:  "Code.cpp",
	},

	"go 1.25": {
		Runtime:       "go 1.25",
		CompileComand: []string{}, // can be skiped if no compilation step
		RunComand:     []string{"go", "run", "Code.go"},
		CodeFileName:  "Code.go",
	},
}

func GetCodeFileName(runtime string) string {
	r, ok := allRuntimes[runtime]
	if ok == false {
		return "Code"
	}
	return r.CodeFileName
}

func GetRuntime(runtime string) (Runtime_conf, bool) {
	r, ok := allRuntimes[runtime]
	return r, ok
}

func Exist(runtime string) bool {
	_, ok := allRuntimes[runtime]
	return ok
}
