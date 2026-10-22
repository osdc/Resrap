//go:build integration

package resrap

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/osdc/resrap/templates"
)

func loadCpp(t *testing.T, r *Resrap) {
	t.Helper()
	if err := templates.Load(r, templates.Cpp); err != nil {
		t.Fatalf("load cpp template: %v", err)
	}
}

func requireTool(t *testing.T, name string) string {
	t.Helper()
	path, err := exec.LookPath(name)
	if err != nil {
		t.Skipf("%s not installed; skipping language integration test", name)
	}
	return path
}

func compileCpp(t *testing.T, gxx, src string) {
	t.Helper()
	cmd := exec.Command(gxx, "-std=c++17", "-Wall", "-Wextra", "-fsyntax-only", src)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("generated C++ did not compile: %v\n%s\n--- source ---\n%s", err, out, mustRead(t, src))
	}
}

func mustRead(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestCppGeneratedCodeCompiles(t *testing.T) {
	gxx := requireTool(t, "g++")
	r := NewResrap()
	loadCpp(t, r)

	dir := t.TempDir()
	const seeds = 25
	for seed := uint64(1); seed <= seeds; seed++ {
		code := r.GenerateCompleteWithSeeded("cpp", "program", seed, 25)
		if code == "" {
			t.Fatalf("seed %d produced empty output", seed)
		}
		if !strings.HasSuffix(strings.TrimRight(code, "\n"), "}") {
			t.Fatalf("seed %d did not finish the top-level rule:\n%s", seed, code)
		}
		path := filepath.Join(dir, "main.cpp")
		if err := os.WriteFile(path, []byte(code), 0o644); err != nil {
			t.Fatal(err)
		}
		compileCpp(t, gxx, path)
	}
}

func TestCppGeneratedProgramRuns(t *testing.T) {
	gxx := requireTool(t, "g++")
	r := NewResrap()
	loadCpp(t, r)
	code := r.GenerateCompleteWithSeeded("cpp", "program", 42, 25)
	dir := t.TempDir()
	src := filepath.Join(dir, "main.cpp")
	bin := filepath.Join(dir, "main")
	if err := os.WriteFile(src, []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}
	if out, err := exec.Command(gxx, "-std=c++17", "-o", bin, src).CombinedOutput(); err != nil {
		t.Fatalf("compile failed: %v\n%s\n%s", err, out, code)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, bin).CombinedOutput()
	if err != nil {
		t.Fatalf("generated program failed: %v\n%s\n%s", err, out, code)
	}
	if len(strings.TrimSpace(string(out))) == 0 {
		t.Fatalf("generated program printed nothing:\n%s", code)
	}
}

func TestCppWordsDirectiveIsUsed(t *testing.T) {
	r := NewResrap()
	loadCpp(t, r)
	words := []string{
		"hello", "world", "resrap", "generated", "code", "cpp", "compiler", "grammar",
		"alpha", "beta", "gamma", "delta", "omega",
	}
	for seed := uint64(1); seed <= 25; seed++ {
		code := r.GenerateCompleteWithSeeded("cpp", "program", seed, 25)
		found := false
		for _, w := range words {
			if strings.Contains(code, `"`+w+`"`) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("seed %d did not use a @words value:\n%s", seed, code)
		}
	}
}

func TestCppDeterministicCompleteGeneration(t *testing.T) {
	r := NewResrap()
	loadCpp(t, r)
	a := r.GenerateCompleteWithSeeded("cpp", "program", 99, 25)
	b := r.GenerateCompleteWithSeeded("cpp", "program", 99, 25)
	if a != b {
		t.Fatalf("seeded complete generation is not deterministic:\n%s\n---\n%s", a, b)
	}
}

func TestCppGenerateCodebaseCompiles(t *testing.T) {
	gxx := requireTool(t, "g++")
	r := NewResrap()
	loadCpp(t, r)
	dir := t.TempDir()
	config := filepath.Join(dir, "project.dsl")
	if err := os.WriteFile(config, []byte("src/\n  cpp[3x25 code_*.cpp]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(dir, "out")
	if err := r.GenerateCodebase(config, target); err != nil {
		t.Fatal(err)
	}
	for i := 1; i <= 3; i++ {
		path := filepath.Join(target, "src", "code_"+string(rune('0'+i))+".cpp")
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("missing generated file %s: %v", path, err)
		}
		compileCpp(t, gxx, path)
	}
}
