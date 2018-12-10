package fargatedeploy

import (
	"fmt"
	"os/exec"
	"strings"
)

// GitGetCommitHash ... get commit hash of specified branch
func GitGetCommitHash() (string, error) {
	out, err := exec.Command(
		"git",
		"rev-parse",
		"HEAD",
	).Output()

	if err != nil {
		return "", fmt.Errorf("GitGetCommitHash: %s", err)
	}

	return strings.TrimRight(string(out), "\n"), nil
}
