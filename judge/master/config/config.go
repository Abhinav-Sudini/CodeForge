package config

import "os"

const Server_Port = 7000

const Max_request_wait_time = 10 // seconds
var InternalApiKey string = os.Getenv("INTERNAL_API_KEY")
var Question_test_case_dir = os.Getenv("QUESTIONS_DIR")

var AllRuntimes = map[string]bool{
	"c++17":   true,
	"c++20":   true,
	"c++23":   true,
	"gcc-c17": true,
	"python3": true,
	"node-25": true,
}
