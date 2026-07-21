package logging

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewWritesLevelFiles(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("LOG_DIR", dir)

	logger, closeLogs, err := New("development")
	if err != nil {
		t.Fatal(err)
	}
	defer closeLogs()

	logger.Debug("debug line", "k", 1)
	logger.Info("info line", "k", 2)
	logger.Warn("warn line", "k", 3)
	logger.Error("error line", "k", 4)
	closeLogs()

	debugBody := readFile(t, filepath.Join(dir, "debug.log"))
	infoBody := readFile(t, filepath.Join(dir, "info.log"))
	errorBody := readFile(t, filepath.Join(dir, "error.log"))

	if !strings.Contains(debugBody, `"msg":"debug line"`) {
		t.Fatalf("debug.log missing debug line: %s", debugBody)
	}
	if strings.Contains(debugBody, `"msg":"info line"`) {
		t.Fatal("debug.log should not contain info line")
	}
	if !strings.Contains(infoBody, `"msg":"info line"`) || !strings.Contains(infoBody, `"msg":"warn line"`) {
		t.Fatalf("info.log missing info/warn: %s", infoBody)
	}
	if strings.Contains(infoBody, `"msg":"error line"`) {
		t.Fatal("info.log should not contain error line")
	}
	if !strings.Contains(errorBody, `"msg":"error line"`) {
		t.Fatalf("error.log missing error line: %s", errorBody)
	}
	if strings.Contains(errorBody, `"msg":"info line"`) {
		t.Fatal("error.log should not contain info line")
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}
