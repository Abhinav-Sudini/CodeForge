package utils

import (
	"master/config"
	"os"
	"path/filepath"
	"strconv"
)

func GetTestCaseStdinStdout(question_id int,wa_test_case int) ([]byte,[]byte) {
		tc_in_path := filepath.Join(config.Question_test_case_dir,strconv.Itoa(question_id),
		strconv.Itoa(wa_test_case)+".in")
		stdin,err := os.ReadFile(tc_in_path)
		if err != nil {
			stdin = []byte{}
		}
		tc_out_path := filepath.Join(config.Question_test_case_dir,strconv.Itoa(question_id),
		strconv.Itoa(wa_test_case)+".out")
		stdout,err := os.ReadFile(tc_out_path)
		if err != nil {
			stdout = []byte{}
		}
		return stdin,stdout
}
