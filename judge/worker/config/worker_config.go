package config

import "os"

// const (
// 	VerdictAccepted Verdict = iota
// 	VerdictWrongAns
// 	VerdictCompilationError
// 	VerdictTLE
// 	VerdictMLE
// 	VerdictCodeToBig
// 	VerdictInternalError
// 	VerdictBadRequest
// )

const Max_req_timeout = 5 //s
const MaxOutputBufferSize = 50<<20

var marster_ip = os.Getenv("MASTER_IP_ADDR")
var port = os.Getenv("MASTER_PORT_CF")
var Master_url = marster_ip + ":" + port

const VerdictApiLocation = "/api/worker/verdict/"
const WorkerRegisterApiLocation = "/api/worker/register/"

var WorkerRuntimeENVVariablesToInclude = []string{"PATH", "HOME"}

const (
	MsgOnAccepted       = "Accepted"
	MsgOnWrongAns       = "Wrong Answer on Test - "
	MsgOnCompilationErr = "Compilation Error"
	MsgOnTLE            = "Time Limit Exceeded on Test - "
	MsgOnMLE            = "Memory Limit Exceeded on Test - "
	MsgOnCodeToBig      = "Code Size To Big"
	MsgOnInternalErr    = "Internal Server Error"
	MsgOnBadRequest     = "Bad Request"
)
