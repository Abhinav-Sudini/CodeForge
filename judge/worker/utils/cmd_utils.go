package utils

import (
	"fmt"
	"os"
)

func CopyEnvVariablesOfParent(variables []string) []string {
	newenv := []string{}
	for _,var_name := range(variables) {
		newenv = append(newenv, fmt.Sprintf("%s=%s",var_name,os.Getenv(var_name)))
	}
	return newenv
}

