package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestBuild(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "ktx")

	cmd := exec.Command("go", "build", "-o", dest, ".")
	cmd.Dir = "."

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, output)
	}

	if _, err := os.Stat(dest); os.IsNotExist(err) {
		t.Fatal("binary not created")
	}
}
