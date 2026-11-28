package tools

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestStagesAndLangs(t *testing.T) {
	stages, err := Stages()
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{"bootstrap", "dev", "quality", "build", "cicd", "kubernetes", "ai", "docs", "database"} {
		if !contains(stages, want) {
			t.Errorf("Stages() = %v, missing %s", stages, want)
		}
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

	if contains(langs, "common") {
		t.Errorf("Langs() = %v, must not list the common directory", langs)
	}
}

func TestNew(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "demo")

	vars := DefaultVars("demo", "github.com/me/demo", "go", "postgres")

	if err := New("demo", "go", "postgres", []string{"dev"}, dir, vars, false); err != nil {
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

	if placeholderRe.Match(data) {
		t.Errorf("main.go still contains placeholders: %q", data)
	}

	gitignore, err := os.ReadFile(filepath.Join(dir, ".gitignore"))
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(gitignore), ".env") {
		t.Error(".gitignore must ignore .env (credentials stay local)")
	}

	for _, f := range []string{"Makefile", ".env.example", "docker-compose.yml", ".gitignore", ".editorconfig", ".gitattributes", "LICENSE", "README.md"} {
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

		if placeholderRe.Match(data) {
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
	err := os.WriteFile(filepath.Join(dir, "keep.txt"), []byte("x"), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	vars := DefaultVars("demo", "", "go", "postgres")
	err = New("demo", "go", "postgres", nil, dir, vars, false)
	if err == nil {
		t.Error("New() into a non-empty dir: want an error, got nil")
	}

	err = New("demo", "go", "postgres", nil, dir, vars, true)

	if err != nil {
		t.Errorf("New() with force: error = %v", err)
	}
}

func TestApply(t *testing.T) {
	dir := t.TempDir()

	vars := DefaultVars("demo", "", "go", "postgres")

	if err := New("demo", "go", "postgres", nil, dir, vars, false); err != nil {
		t.Fatal(err)
	}

	if err := Apply("dev", "go", "postgres", dir, vars, false); err != nil {
		t.Fatalf("Apply(dev) error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "Makefile")); err != nil {
		t.Errorf("Makefile missing after Apply: %v", err)
	}

	if err := os.WriteFile(filepath.Join(dir, "Makefile"), []byte("# custom\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := Apply("dev", "go", "postgres", dir, vars, false); err == nil {
		t.Error("Apply() over an existing file without force: want an error, got nil")
	}

	if err := Apply("dev", "go", "postgres", dir, vars, true); err != nil {
		t.Errorf("Apply() with force: error = %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "Makefile"))
	if err != nil || string(data) == "# custom\n" {
		t.Errorf("Apply() with force did not overwrite Makefile: %q, %v", data, err)
	}

	if err := Apply("bootstrap", "go", "postgres", dir, vars, false); err == nil {
		t.Error("Apply(bootstrap): want an error, got nil")
	}
}

func TestNewWithDevAndQualityStages(t *testing.T) {
	langFiles := map[string][]string{
		"go":     {"Makefile.lang", ".tool-versions", ".golangci.yml", "scripts/lint-report.sh", ".goreleaser.yaml", "Dockerfile", ".dockerignore", ".github/workflows/ci.yml", ".github/workflows/cd.yml", ".github/workflows/docker.yml", ".github/workflows/release.yml", ".github/dependabot.yml", ".github/pull_request_template.md", ".github/ISSUE_TEMPLATE/bug_report.yml", "PIPELINE.md", "deployments/kustomization.yaml", "deployments/deployment.yaml", "deployments/service.yaml", "deployments/ingress.yaml", "deployments/hpa.yaml", "deployments/overlays/dev/kustomization.yaml", "deployments/overlays/staging/kustomization.yaml", "deployments/overlays/prod/kustomization.yaml", "skaffold.yaml", "AGENTS.md", ".skills/project-spec/SKILL.md", ".skills/go/SKILL.md", ".skills/tdd/SKILL.md", ".skills/clean-code/SKILL.md", ".skills/clean-architecture/SKILL.md", "docs/README.md", "docs/spec.md", "docs/adr/0000-template.md", "docs/runbook.md", "docs/git-workflow.md", "CONTRIBUTING.md", "atlas.postgres.hcl", ".env.postgres.example", "migrations/README.md", "scripts/check-branch-name.sh", "githooks/pre-push", ".skills/sql-migrations/SKILL.md", ".skills/git-conventions/SKILL.md"},
		"node":   {"Makefile.lang", ".tool-versions", "eslint.config.js", ".prettierrc", "commitlint.config.js", "Dockerfile", ".dockerignore", ".github/workflows/ci.yml", ".github/workflows/cd.yml", ".github/workflows/docker.yml", ".github/dependabot.yml", ".github/pull_request_template.md", "PIPELINE.md", "AGENTS.md", ".skills/project-spec/SKILL.md", ".skills/tdd/SKILL.md", ".skills/clean-code/SKILL.md", ".skills/clean-architecture/SKILL.md", "docs/README.md", "docs/spec.md", "docs/adr/0000-template.md", "docs/runbook.md", "docs/git-workflow.md", "CONTRIBUTING.md", "atlas.postgres.hcl", ".env.postgres.example", "migrations/README.md", "scripts/check-branch-name.sh", "githooks/pre-push", ".skills/sql-migrations/SKILL.md", ".skills/git-conventions/SKILL.md"},
		"python": {"Makefile.lang", ".tool-versions", "ruff.toml", "mypy.ini", "Dockerfile", ".dockerignore", ".github/workflows/ci.yml", ".github/workflows/cd.yml", ".github/workflows/docker.yml", ".github/dependabot.yml", ".github/pull_request_template.md", "PIPELINE.md", "AGENTS.md", ".skills/project-spec/SKILL.md", ".skills/tdd/SKILL.md", ".skills/clean-code/SKILL.md", ".skills/clean-architecture/SKILL.md", "docs/README.md", "docs/spec.md", "docs/adr/0000-template.md", "docs/runbook.md", "docs/git-workflow.md", "CONTRIBUTING.md", "atlas.postgres.hcl", ".env.postgres.example", "migrations/README.md", "scripts/check-branch-name.sh", "githooks/pre-push", ".skills/sql-migrations/SKILL.md", ".skills/git-conventions/SKILL.md"},
	}
	for lang, files := range langFiles {
		t.Run(lang, func(t *testing.T) {
			dir := filepath.Join(t.TempDir(), "demo")

			vars := DefaultVars("demo", "", lang, "postgres")

			if err := New("demo", lang, "postgres", []string{"dev", "quality", "build", "cicd", "kubernetes", "ai", "docs", "database"}, dir, vars, false); err != nil {
				t.Fatalf("New() error = %v", err)
			}

			for _, f := range append(files, ".pre-commit-config.yaml", ".markdownlint.json", "Makefile") {
				if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
					t.Errorf("lang %s: expected %s: %v", lang, f, err)
				}
			}

			data, err := os.ReadFile(filepath.Join(dir, "Makefile.lang"))
			if err != nil {
				t.Fatal(err)
			}

			if placeholderRe.Match(data) {
				t.Errorf("lang %s: Makefile.lang has unresolved placeholders: %q", lang, data)
			}

			for _, cmd := range []string{"FMT_CMD", "FMT_CHECK_CMD", "SECURITY_CMD"} {
				if !strings.Contains(string(data), cmd+" =") {
					t.Errorf("lang %s: Makefile.lang missing %s", lang, cmd)
				}
			}

			mk, err := os.ReadFile(filepath.Join(dir, "Makefile"))
			if err != nil {
				t.Fatal(err)
			}

			for _, target := range []string{"fmt-check:", "security:", "lint:", "test:", "build:"} {
				if !strings.Contains(string(mk), target) {
					t.Errorf("lang %s: Makefile missing target %q", lang, target)
				}
			}

			ci, err := os.ReadFile(filepath.Join(dir, ".github", "workflows", "ci.yml"))
			if err != nil {
				t.Fatal(err)
			}

			for _, job := range []string{"pre-checks:", "format:", "lint:", "security:", "test:", "build:"} {
				if !strings.Contains(string(ci), job) {
					t.Errorf("lang %s: ci.yml missing gate job %q", lang, job)
				}
			}

			if lang == "go" && !strings.Contains(string(data), "demo") {
				t.Errorf("lang go: Makefile.lang was not rendered with the project name: %q", data)
			}

			dockerfile, err := os.ReadFile(filepath.Join(dir, "Dockerfile"))
			if err != nil {
				t.Fatal(err)
			}

			if lang == "go" && !strings.Contains(string(dockerfile), "demo") {
				t.Errorf("lang go: Dockerfile was not rendered with the project name")
			}

			if placeholderRe.Match(dockerfile) {
				t.Errorf("lang %s: Dockerfile has unresolved placeholders", lang)
			}

			if lang == "go" {
				info, err := os.Stat(filepath.Join(dir, "scripts", "lint-report.sh"))
				if err != nil {
					t.Fatal(err)
				}

				if info.Mode().Perm()&0o111 == 0 {
					t.Errorf("scripts/lint-report.sh lost its executable bit: %v", info.Mode())
				}
			}

			hook, err := os.Stat(filepath.Join(dir, "githooks", "pre-push"))
			if err != nil {
				t.Fatal(err)
			}
			if hook.Mode().Perm()&0o111 == 0 {
				t.Errorf("githooks/pre-push lost its executable bit: %v", hook.Mode())
			}

			for _, other := range []string{"go", "node", "python"} {
				if _, err := os.Stat(filepath.Join(dir, other)); err == nil {
					t.Errorf("lang %s: stage language directories must not be copied verbatim (found %s/)", lang, other)
				}
			}

			goSkill := filepath.Join(dir, ".skills", "go", "SKILL.md")
			if lang == "go" {
				if _, err := os.Stat(goSkill); err != nil {
					t.Errorf("lang go: expected the go skill at %s: %v", goSkill, err)
				}
			} else if _, err := os.Stat(goSkill); err == nil {
				t.Errorf("lang %s: .skills/go must only be scaffolded for go projects", lang)
			}
		})
	}
}

func TestCicdWorkflowsRender(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "demo")

	vars := DefaultVars("demo", "github.com/me/demo", "go", "postgres")

	if err := New("demo", "go", "postgres", []string{"cicd"}, dir, vars, false); err != nil {
		t.Fatal(err)
	}

	wf, err := os.ReadFile(filepath.Join(dir, ".github", "workflows", "docker.yml"))
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(wf), "${{ github.repository }}") {
		t.Errorf("docker.yml lost its GitHub Actions expression:\n%s", wf)
	}

	if placeholderRe.Match(wf) {
		t.Errorf("docker.yml has unresolved placeholders:\n%s", wf)
	}
}

func TestDBs(t *testing.T) {
	dbs, err := DBs()
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{"postgres", "mysql", "sqlite"} {
		if !contains(dbs, want) {
			t.Errorf("DBs() = %v, missing %s", dbs, want)
		}
	}
}

func TestDatabaseStageOverlay(t *testing.T) {
	for _, db := range []string{"postgres", "mysql", "sqlite"} {
		t.Run(db, func(t *testing.T) {
			dir := filepath.Join(t.TempDir(), "demo")
			vars := DefaultVars("demo", "", "go", db)

			err := New("demo", "go", db, []string{"database"}, dir, vars, false)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			atlas, err := os.ReadFile(filepath.Join(dir, "atlas."+db+".hcl"))
			if err != nil {
				t.Fatal(err)
			}

			dbEnv := strings.ToUpper(db)
			if db != "sqlite" && !strings.Contains(string(atlas), dbEnv+"_DATABASE_URL") {
				t.Errorf("db %s: atlas.%s.hcl does not read %s_Database_URL from env:\n%s", db, db, dbEnv, atlas)
			}

			if strings.Contains(string(atlas), "dev:dev@") {
				t.Errorf("db %s: atlas config must not embed credentials:\n%s", db, atlas)
			}

			envExample, err := os.ReadFile(filepath.Join(dir, ".env."+db+".example"))
			if err != nil {
				t.Fatal(err)
			}

			if !strings.Contains(string(envExample), dbEnv+"_DATABASE_URL=") &&
				!strings.Contains(string(envExample), dbEnv+"_STAGING_DATABASE_URL=") {
				t.Errorf("db %s: .env.%s.example missing a %s_*_DATABASE_URL entry", db, db, dbEnv)
			}

			if _, err := os.Stat(filepath.Join(dir, "migrations", "README.md")); err != nil {
				t.Errorf("db %s: migrations/README.md missing: %v", db, err)
			}

			for _, other := range []string{"postgres", "mysql", "sqlite"} {
				if other == db {
					continue
				}

				if _, err := os.Stat(filepath.Join(dir, other)); err == nil {
					t.Errorf("db %s: overlay directory %q must not be copied verbatim", db, other)
				}
			}
		})
	}
}

func TestMultipleDatabases(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "demo")
	vars := DefaultVars("demo", "", "go", "postgres,mysql")

	err := New("demo", "go", "postgres,mysql", []string{"database"}, dir, vars, false)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	for _, f := range []string{"atlas.postgres.hcl", "atlas.mysql.hcl", ".env.postgres.example", ".env.mysql.example", "migrations/README.md"} {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Errorf("multi-db: expected %s: %v", f, err)
		}
	}

	atlas, err := os.ReadFile(filepath.Join(dir, "atlas.mysql.hcl"))
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(atlas), "MYSQL_DATABASE_URL") {
		t.Errorf("atlas.mysql.hcl was not rendered with DB=mysql:\n%s", atlas)
	}
}

func TestGolangciConfigSeverity(t *testing.T) {
	data, err := templatesFS.ReadFile("templates/quality/go/.golangci.yml")
	if err != nil {
		t.Fatal(err)
	}

	cfg := string(data)
	for _, want := range []string{"default: all", "default: warning", "severity: error", "severity: warning", "severity: info"} {
		if !strings.Contains(cfg, want) {
			t.Errorf(".golangci.yml missing %q", want)
		}
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
	v := DefaultVars("demo", "", "go", "postgres")
	if v["MODULE"] != "demo" {
		t.Errorf("MODULE = %q, want %q", v["MODULE"], "demo")
	}

	if v["IMAGE"] != "demo" {
		t.Errorf("IMAGE = %q, want %q", v["IMAGE"], "demo")
	}

	if v["PROJECT"] != "demo" || v["LANG"] != "go" || v["YEAR"] == "" {
		t.Errorf("DefaultVars() = %v", v)
	}

	v = DefaultVars("demo", "github.com/me/demo", "go", "postgres")
	if v["IMAGE"] != "ghcr.io/me/demo" {
		t.Errorf("IMAGE = %q, want %q", v["IMAGE"], "ghcr.io/me/demo")
	}
}

func contains(list []string, s string) bool {
	return slices.Contains(list, s)
}
