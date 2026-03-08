package clojurev

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// assertContains ensures that a substring exists within the source string.
func assertContains(t *testing.T, source, substr, msg string) {
	t.Helper()
	if !strings.Contains(source, substr) {
		t.Errorf("[%s] expected to find %q in output\nGot: %s", msg, substr, source)
	}
}

// assertNotContains ensures that a substring does NOT exist within the source string.
func assertNotContains(t *testing.T, source, substr, msg string) {
	t.Helper()
	if strings.Contains(source, substr) {
		t.Errorf("[%s] expected NOT to find %q in output\nGot: %s", msg, substr, source)
	}
}

// runBackendCommand executes a command and returns its trimmed output.
func runBackendCommand(t *testing.T, name string, args ...string) string {
	t.Helper()
	output, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		t.Fatalf("Execution of %s failed: %v\n%s", name, err, string(output))
	}
	return strings.TrimSpace(string(output))
}

// validateTranspilation performs a generic transpilation and checks for multiple expected fragments.
func validateTranspilation(t *testing.T, cljv string, target Target, expected []string, msg string) string {
	t.Helper()
	output, err := Transpile(cljv, target, "")
	if err != nil {
		t.Fatalf("%s: Transpilation to %v failed: %v", msg, target, err)
	}
	for _, exp := range expected {
		assertContains(t, output, exp, fmt.Sprintf("%s (%v)", msg, target))
	}
	return output
}

// transpileAndVerifyTable performs table-driven transpilation verification.
func transpileAndVerifyTable(t *testing.T, target Target, tests []struct {
	name     string
	cljv     string
	expected []string
}) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := Transpile(tt.cljv, target, "")
			if err != nil {
				t.Fatalf("[%s] Transpilation failed: %v", tt.name, err)
			}
			for _, exp := range tt.expected {
				assertContains(t, output, exp, tt.name)
			}
		})
	}
}

// setupTempFile writes content to a temporary file and returns its path.
func setupTempFile(t *testing.T, prefix, suffix, content string) string {
	t.Helper()
	tmpFile := filepath.Join(t.TempDir(), prefix+suffix)
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create temp file %s: %v", tmpFile, err)
	}
	return tmpFile
}
