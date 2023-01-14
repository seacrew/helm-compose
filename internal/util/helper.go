package util

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func IsDebug() bool {
	return os.Getenv("HELM_DEBUG") == "true"
}

func DebugPrint(format string, a ...interface{}) {
	if IsDebug() {
		fmt.Printf(format+"\n", a...)
	}
}

func Execute(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()

	if output == nil {
		return "", err
	}

	text := string(output)

	if len(text) > 5 && text[:6] == "Error:" {
		text = strings.TrimSpace(text[6:])
	}

	return text, err
}
