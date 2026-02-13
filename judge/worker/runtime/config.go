package runtime

type Runtime_conf struct{
	Runtime string
	CompileComand []string
	RunComand []string
	FileExtention string
	CodeFileName string
}

const BinaryFileName = "out.bin"
const StdErrorFileName = "std.err"
const StdOutFileName = "std.out"

var allRuntimes = map[string]Runtime_conf{

	"c++17":{
		Runtime: "c++17",
		CompileComand: []string{"g++","-o",BinaryFileName,"Code.cpp"},
		RunComand: []string{"./out.bin"},
		CodeFileName: "Code.cpp",
	},

	"go 1.25":{
		Runtime: "go 1.25",
		CompileComand: []string{}, // can be skiped if no compilation step
		RunComand: []string{"go","run","Code.go"},
		CodeFileName: "Code.go",
	},
}


func GetRuntime(runtime string)(Runtime_conf,bool){
	r,ok := allRuntimes[runtime]
	return r,ok
}

func Exist(runtime string) bool {
	_,ok := allRuntimes[runtime]
	return ok
}

