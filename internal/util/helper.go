package util

import (
	"fmt"
	"os"
)

func IsDebug() bool {
	return os.Getenv("HELM_DEBUG") == "true"
}

func DebugPrint(format string, a ...interface{}) {
	if IsDebug() {
		fmt.Printf(format+"\n", a...)
	}
}
