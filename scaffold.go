package tools

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

//go:embed all:templates
var templatesFS embed.FS

const templatesRoot = "templates"

var placeholderRe = regexp.MustCompile(`\{\{[A-Z0-9_]+\}\}`)

// Vars holds the values substituted into {{PLACEHOLDER}}s.
type Vars map[string]string

// DefaultVars returns the standard variables for a project. module defaults
// to the project name when empty.
func DefaultVars(project, module, lang string) Vars {
	if module == "" {
		module = project
	}
	return Vars{
		"PROJECT": project,
		"MODULE":  module,
		"LANG":    lang,
		"YEAR":    fmt.Sprintf("%d", time.Now().Year()),
	}
}

// Stages lists the template stages available under templates/.
func Stages() ([]string, error) {
	return subdirs(templatesRoot)
}

// Langs lists the languages available for the bootstrap stage. The "common"
// directory holds language-independent files and is not a language.
func Langs() ([]string, error) {
	dirs, err := subdirs(templatesRoot + "/bootstrap")
	if err != nil {
		return nil, err
	}
	langs := make([]string, 0, len(dirs))
	for _, d := range dirs {
		if d != "common" {
			langs = append(langs, d)
		}
	}
	return langs, nil
}

// New scaffolds a new project: the bootstrap template for lang plus any
// extra stages, rendered into dir (default: ./<project>). It refuses to
// write into a non-empty directory unless force is true.
func New(project, lang string, stages []string, dir string, vars Vars, force bool) error {
	if dir == "" {
		dir = project
	}
	if err := checkTarget(dir, force); err != nil {
		return err
	}
	all := []string{"bootstrap"}
	for _, s := range stages {
		if s != "bootstrap" {
			all = append(all, s)
		}
	}
	for _, stage := range all {
		if err := copyStage(stage, lang, dir, vars, force); err != nil {
			return err
		}
	}
	return nil
}

// Apply renders a single stage into an existing project directory. The
// bootstrap stage is only available through New.
func Apply(stage, lang, dir string, vars Vars, force bool) error {
	if stage == "bootstrap" {
		return fmt.Errorf("tools: the bootstrap stage is only available via New")
	}
	return copyStage(stage, lang, dir, vars, force)
}

// Render replaces every {{KEY}} in s with vars[KEY]. It returns an error if
// a placeholder is left unresolved.
func Render(s string, vars Vars) (string, error) {
	out := placeholderRe.ReplaceAllStringFunc(s, func(p string) string {
		if v, ok := vars[p[2:len(p)-2]]; ok {
			return v
		}
		return p
	})
	if left := placeholderRe.FindString(out); left != "" {
		return "", fmt.Errorf("tools: unresolved placeholder %s", left)
	}
	return out, nil
}

func copyStage(stage, lang, dir string, vars Vars, force bool) error {
	if stage == "bootstrap" {
		common := templatesRoot + "/bootstrap/common"
		if _, err := templatesFS.ReadDir(common); err == nil {
			if err := copyTree(common, dir, vars, force, nil); err != nil {
				return err
			}
		}
		base := templatesRoot + "/bootstrap/" + lang
		if _, err := templatesFS.ReadDir(base); err != nil {
			return fmt.Errorf("tools: unknown language %q for stage %q", lang, stage)
		}
		return copyTree(base, dir, vars, force, nil)
	}
	base := templatesRoot + "/" + stage
	if _, err := templatesFS.ReadDir(base); err != nil {
		return fmt.Errorf("tools: unknown stage %q", stage)
	}
	langs, err := Langs()
	if err != nil {
		return err
	}
	skip := map[string]bool{}
	for _, l := range langs {
		skip[l] = true
	}
	if err := copyTree(base, dir, vars, force, skip); err != nil {
		return err
	}
	if lang != "" {
		overlay := base + "/" + lang
		if _, err := templatesFS.ReadDir(overlay); err == nil {
			if err := copyTree(overlay, dir, vars, force, nil); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyTree(base, dir string, vars Vars, force bool, skipLangs map[string]bool) error {
	return fs.WalkDir(templatesFS, base, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && path != base && skipLangs[d.Name()] {
			return fs.SkipDir
		}
		rel, err := filepath.Rel(base, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		name, err := Render(rel, vars)
		if err != nil {
			return fmt.Errorf("tools: %s: %w", path, err)
		}
		name = strings.TrimSuffix(name, ".tmpl")
		target := filepath.Join(dir, filepath.FromSlash(name))
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if !force {
			if _, err := os.Stat(target); err == nil {
				return fmt.Errorf("tools: %s already exists (use force to overwrite)", target)
			}
		}
		data, err := templatesFS.ReadFile(path)
		if err != nil {
			return err
		}
		content, err := Render(string(data), vars)
		if err != nil {
			return fmt.Errorf("tools: %s: %w", path, err)
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		perm := fs.FileMode(0o644)
		if info, err := d.Info(); err == nil && info.Mode().Perm()&0o111 != 0 {
			perm = 0o755
		}
		return os.WriteFile(target, []byte(content), perm)
	})
}

func checkTarget(dir string, force bool) error {
	if force {
		return nil
	}
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if len(entries) > 0 {
		return fmt.Errorf("tools: %s is not empty (use force to overwrite)", dir)
	}
	return nil
}

func subdirs(root string) ([]string, error) {
	entries, err := templatesFS.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("tools: %w", err)
	}
	var dirs []string
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e.Name())
		}
	}
	sort.Strings(dirs)
	return dirs, nil
}
