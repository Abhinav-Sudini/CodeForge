package config

import "os"

const Server_Port = 7000

const Max_request_wait_time = 10 // seconds
var InternalApiKey string = os.Getenv("INTERNAL_API_KEY")

