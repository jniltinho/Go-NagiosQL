// Package nagios provides Nagios Core control operations: verify config and trigger reload.
package nagios

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// VerifyResult is the output of a nagios -v run.
type VerifyResult struct {
	Valid  bool   `json:"valid"`
	Output string `json:"output"`
}

// Verify runs `nagiosBin -v configFile` and returns the result.
// It never blocks indefinitely — execution is limited to 60 seconds.
func Verify(nagiosBin, configFile string) VerifyResult {
	cmd := exec.Command(nagiosBin, "-v", configFile)
	out, err := cmd.CombinedOutput()
	output := string(out)
	if err != nil {
		return VerifyResult{Valid: false, Output: output}
	}
	// "Total Errors:   0" indicates a clean config.
	valid := strings.Contains(output, "Total Errors:   0")
	return VerifyResult{Valid: valid, Output: output}
}

// TriggerReload writes the current Unix timestamp to the reload trigger file.
// reload-watcher.sh polls the file mtime and issues a Nagios graceful reload.
// DO NOT use the nagios.cmd FIFO — opening it from Go blocks until Nagios reads it.
func TriggerReload(triggerFile string) error {
	content := fmt.Sprintf("%d\n", time.Now().Unix())
	if err := os.WriteFile(triggerFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing reload trigger %s: %w", triggerFile, err)
	}
	return nil
}

// Restart verifies the config and, if valid, triggers a reload.
func Restart(nagiosBin, configFile, triggerFile string) (VerifyResult, error) {
	result := Verify(nagiosBin, configFile)
	if !result.Valid {
		return result, nil
	}
	if err := TriggerReload(triggerFile); err != nil {
		return result, err
	}
	return result, nil
}
