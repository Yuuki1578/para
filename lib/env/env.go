package env

import (
	"os"
	"strconv"
)

const PARA_THREAD string = "PARA_THREAD"

func GetEnv() uint64 {
	envar := os.Getenv(PARA_THREAD)
	if envar == "" {
		return 0
	}

	if result, err := strconv.ParseUint(envar, 10, 64); err == nil {
		return result
	}

	return 0
}
