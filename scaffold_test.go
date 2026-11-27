package tools

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStagesAndLangs(t *testing.T) {
	stages, err := Stages()
	if err != nil {
		t.Fatal(err)
	}
	if !contains(stages, "bootstrap") || !contains(stages, "dev") {
		t.Errorf("Stages() = %v, want bootstrap and dev", stages)
	}
	langs, err := Langs()
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"go", "node", "python"} {
		if !contains(langs, want) {
			t.Errorf("Langs() = %v, missing %s", langs, want)
		}
	}
}

func TestNew(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "demo")
	vars := DefaultVars("demo", "github.com/me/demo", "go")
	if err := New("demo", "go", []string{"dev"}, dir, vars, false); err != nil {
		t.Fatalf("New() error = %v", err)
	}

	gomod, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(gomod), "module github.com/me/demo") {
		t.Errorf("go.mod = %q, want module github.com/me/demo", gomod)
	}

	data, err := os.ReadFile(filepath.Join(dir, "cmd", "demo", "main.go"))
	if err != nil {
		t.Fatalf("rendered cmd/demo/main.go missing: %v", err)
	}
	if strings.Contains(string(data), "{{") {
		t.Errorf("main.go still contains placeholders: %q", data)
	}

	for _, f := range []string{"Makefile", ".env.example", "docker-compose.yml", ".gitignore", ".editorconfig", "README.md"} {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Errorf("expected %s in scaffolded project: %v", f, err)
		}
	}

	for _, other := range []string{"node", "python"} {
		if _, err := os.Stat(filepath.Join(dir, other)); err == nil {
			t.Errorf("scaffolded project must not contain the %q language directory", other)
		}
	}

	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if strings.Contains(string(data), "{{") {
			t.Errorf("%s contains an unresolved placeholder", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewNonEmptyDir(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "keep.txt"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	vars := DefaultVars("demo", "", "go")
	if err := New("demo", "go", nil, dir, vars, false); err == nil {
		t.Error("New() into a non-empty dir: want an error, got nil")
	}
	if err := New("demo", "go", nil, dir, vars, true); err != nil {
		t.Errorf("New() with force: error = %v", err)
	}
}

func TestApply(t *testing.T) {
	dir := t.TempDir()
	vars := DefaultVars("demo", "", "go")
	if err := New("demo", "go", nil, dir, vars, false); err != nil {
		t.Fatal(err)
	}
	if err := Apply("dev", "go", dir, vars, false); err != nil {
		t.Fatalf("Apply(dev) error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "Makefile")); err != nil {
		t.Errorf("Makefile missing after Apply: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "Makefile"), []byte("# custom\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := Apply("dev", "go", dir, vars, false); err == nil {
		t.Error("Apply() over an existing file without force: want an error, got nil")
	}
	if err := Apply("dev", "go", dir, vars, true); err != nil {
		t.Errorf("Apply() with force: error = %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "Makefile"))
	if err != nil || string(data) == "# custom\n" {
		t.Errorf("Apply() with force did not overwrite Makefile: %q, %v", data, err)
	}
	if err := Apply("bootstrap", "go", dir, vars, false); err == nil {
		t.Error("Apply(bootstrap): want an error, got nil")
	}
}

func TestRender(t *testing.T) {
	got, err := Render("hello {{PROJECT}} ({{YEAR}})", Vars{"PROJECT": "demo", "YEAR": "2025"})
	if err != nil {
		t.Fatal(err)
	}
	if got != "hello demo (2025)" {
		t.Errorf("Render() = %q, want %q", got, "hello demo (2025)")
	}
	if _, err := Render("{{UNKNOWN}}", Vars{}); err == nil {
		t.Error("Render() with an unknown placeholder: want an error, got nil")
	}
}

func TestDefaultVars(t *testing.T) {
	v := DefaultVars("demo", "", "go")
	if v["MODULE"] != "demo" {
		t.Errorf("MODULE = %q, want %q", v["MODULE"], "demo")
	}
	if v["PROJECT"] != "demo" || v["LANG"] != "go" || v["YEAR"] == "" {
		t.Errorf("DefaultVars() = %v", v)
	}
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}
