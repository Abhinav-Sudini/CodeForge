package runner

import (
	"strconv"
	"worker/config"
)

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

func GenerateResultMSG(res SubmitionResult) string {
	str_test := strconv.Itoa(res.WA_Test_case)
	switch Verdict(res.Result) {
	case VerdictAccepted:
		return config.MsgOnAccepted
	case VerdictWrongAns:
		return config.MsgOnWrongAns + str_test
	case VerdictCompilationError:
		return config.MsgOnCompilationErr
	case VerdictTLE:
		return config.MsgOnTLE + str_test
	case VerdictMLE:
		return config.MsgOnMLE + str_test
	case VerdictCodeToBig:
		return config.MsgOnCodeToBig
	case VerdictInternalError:
		return config.MsgOnInternalErr
	case VerdictBadRequest:
		return config.MsgOnBadRequest
	}
	return "idk ig someting is wrong if i hit this code"
}
