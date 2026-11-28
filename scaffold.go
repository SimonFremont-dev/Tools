// Package tools renders the embedded templates/ tree into new projects.
package tools

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

//go:embed all:templates
var templatesFS embed.FS

const (
	templatesRoot  = "templates"
	stageBootstrap = "bootstrap"
	defaultPerm    = 0o644
	execPerm       = 0o755
	dirPerm        = 0o750
)

var (
	placeholderRe = regexp.MustCompile(`\{\{[A-Z0-9_]+\}\}`)

	// ErrUnknownStage is returned when a stage does not exist under templates/.
	ErrUnknownStage = errors.New("tools: unknown stage")
	// ErrUnknownLanguage is returned when a language has no bootstrap template.
	ErrUnknownLanguage = errors.New("tools: unknown language")
	// ErrBootstrapViaNew is returned when Apply is called with the bootstrap stage.
	ErrBootstrapViaNew = errors.New("tools: the bootstrap stage is only available via New")
	// ErrUnresolvedPlaceholder is returned when a template keeps a {{PLACEHOLDER}}.
	ErrUnresolvedPlaceholder = errors.New("tools: unresolved placeholder")
	// ErrTargetExists is returned when a file would be overwritten without force.
	ErrTargetExists = errors.New("tools: file already exists (use force to overwrite)")
	// ErrTargetNotEmpty is returned when the target directory is not empty.
	ErrTargetNotEmpty = errors.New("tools: directory is not empty (use force to overwrite)")
)

// Vars holds the values substituted into {{PLACEHOLDER}}s.
type Vars map[string]string

// DefaultVars returns the standard variables for a project. module defaults
// to the project name when empty. IMAGE is the container image reference:
// ghcr.io/<owner>/<repo> when module is a github.com path, else the project
// name (for local builds).
func DefaultVars(project, module, lang, database string) Vars {
	if module == "" {
		module = project
	}

	if database == "" {
		database = "postgres"
	}

	image := project

	if rest, ok := strings.CutPrefix(module, "github.com/"); ok {
		image = "ghcr.io/" + rest
	}

	return Vars{
		"PROJECT": project,
		"MODULE":  module,
		"LANG":    lang,
		"DB":      database,
		"YEAR":    strconv.Itoa(time.Now().Year()),
		"IMAGE":   image,
	}
}

// Stages lists the template stages available under templates/.
func Stages() ([]string, error) {
	return subdirs(templatesRoot)
}

// Langs lists the languages available for the bootstrap stage. The "common"
// directory holds language-independent files and is not a language.
func Langs() ([]string, error) {
	dirs, err := subdirs(templatesRoot + "/" + stageBootstrap)
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

// DBs lists the databases available as overlays of the database stage: the
// top-level subdirectories of templates/database that contain an atlas.hcl.
// It returns nil when the stage does not exist.
func DBs() ([]string, error) {
	root := templatesRoot + "/database"

	dirs, err := subdirs(root)
	if err != nil {
		return nil, nil //nolint:nilerr // no database stage means no overlays
	}

	var dbs []string

	for _, name := range dirs {
		entries, err := templatesFS.ReadDir(root + "/" + name)
		if err != nil {
			continue
		}

		for _, e := range entries {
			if e.Name() == "atlas.hcl" {
				dbs = append(dbs, name)

				break
			}
		}
	}

	return dbs, nil
}

// New scaffolds a new project: the bootstrap template for lang plus any
// extra stages, rendered into dir (default: ./<project>). It refuses to
// write into a non-empty directory unless force is true.
func New(project, lang, database string, stages []string, dir string, vars Vars, force bool) error {
	if dir == "" {
		dir = project
	}

	err := checkTarget(dir, force)
	if err != nil {
		return err
	}

	all := []string{stageBootstrap}

	for _, s := range stages {
		if s != stageBootstrap {
			all = append(all, s)
		}
	}

	for _, stage := range all {
		err = copyStage(stage, lang, database, dir, vars, force)
		if err != nil {
			return err
		}
	}

	return nil
}

// Apply renders a single stage into an existing project directory. The
// bootstrap stage is only available through New.
func Apply(stage, lang, database, dir string, vars Vars, force bool) error {
	if stage == stageBootstrap {
		return ErrBootstrapViaNew
	}

	return copyStage(stage, lang, database, dir, vars, force)
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
		return "", fmt.Errorf("%w: %s", ErrUnresolvedPlaceholder, left)
	}

	return out, nil
}

func copyStage(stage, lang, database, dir string, vars Vars, force bool) error {
	if stage == stageBootstrap {
		return copyBootstrap(lang, dir, vars, force)
	}

	return copyOtherStage(stage, lang, database, dir, vars, force)
}

func copyBootstrap(lang, dir string, vars Vars, force bool) error {
	common := templatesRoot + "/" + stageBootstrap + "/common"

	_, commonErr := templatesFS.ReadDir(common)
	if commonErr == nil {
		err := copyTree(common, dir, vars, force, nil)
		if err != nil {
			return err
		}
	}

	base := templatesRoot + "/" + stageBootstrap + "/" + lang

	_, err := templatesFS.ReadDir(base)
	if err != nil {
		return fmt.Errorf("%w %q for stage %q", ErrUnknownLanguage, lang, stageBootstrap)
	}

	return copyTree(base, dir, vars, force, nil)
}

func copyOtherStage(stage, lang, database, dir string, vars Vars, force bool) error {
	base := templatesRoot + "/" + stage

	_, err := templatesFS.ReadDir(base)
	if err != nil {
		return fmt.Errorf("%w %q", ErrUnknownStage, stage)
	}

	skip, err := overlayDirs()
	if err != nil {
		return err
	}

	err = copyTree(base, dir, vars, force, skip)
	if err != nil {
		return err
	}

	for _, ov := range []string{lang, database} {
		if ov == "" {
			continue
		}

		overlay := base + "/" + ov

		_, overlayErr := templatesFS.ReadDir(overlay)
		if overlayErr == nil {
			err = copyTree(overlay, dir, vars, force, nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// overlayDirs returns the top-level stage subdirectories that are overlays
// (languages and databases) rather than content.
func overlayDirs() (map[string]bool, error) {
	langs, err := Langs()
	if err != nil {
		return nil, err
	}

	dbs, err := DBs()
	if err != nil {
		return nil, err
	}

	skip := make(map[string]bool, len(langs)+len(dbs))

	for _, l := range langs {
		skip[l] = true
	}

	for _, d := range dbs {
		skip[d] = true
	}

	return skip, nil
}

func copyTree(base, dir string, vars Vars, force bool, skipLangs map[string]bool) error {
	walk := func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() && skipLangs[entry.Name()] && path == base+"/"+entry.Name() {
			return fs.SkipDir
		}

		return copyEntry(base, dir, path, entry, vars, force)
	}

	err := fs.WalkDir(templatesFS, base, walk)
	if err != nil {
		return fmt.Errorf("tools: walk %s: %w", base, err)
	}

	return nil
}

func copyEntry(base, dir, path string, entry fs.DirEntry, vars Vars, force bool) error {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return fmt.Errorf("tools: resolve %s: %w", path, err)
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

	if entry.IsDir() {
		return makeDir(target)
	}

	err = targetExists(target, force)
	if err != nil {
		return err
	}

	content, err := renderTemplate(path, vars)
	if err != nil {
		return err
	}

	err = makeDir(filepath.Dir(target))
	if err != nil {
		return err
	}

	return writeFile(target, content, filePerm(name))
}

func makeDir(target string) error {
	err := os.MkdirAll(target, dirPerm)
	if err != nil {
		return fmt.Errorf("tools: mkdir %s: %w", target, err)
	}

	return nil
}

func filePerm(name string) fs.FileMode {
	if strings.HasSuffix(name, ".sh") || strings.Contains(name, "/githooks/") || strings.HasPrefix(name, "githooks/") {
		return execPerm
	}

	return defaultPerm
}

func writeFile(target, content string, perm fs.FileMode) error {
	err := os.WriteFile(target, []byte(content), perm) // #nosec G304 -- target derives from embedded template names
	if err != nil {
		return fmt.Errorf("tools: write %s: %w", target, err)
	}

	return nil
}

func renderTemplate(path string, vars Vars) (string, error) {
	data, err := templatesFS.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("tools: read template %s: %w", path, err)
	}

	content, err := Render(string(data), vars)
	if err != nil {
		return "", fmt.Errorf("tools: %s: %w", path, err)
	}

	return content, nil
}

func targetExists(target string, force bool) error {
	if force {
		return nil
	}

	_, err := os.Stat(target) // #nosec G304 -- target derives from embedded template names
	if err == nil {
		return fmt.Errorf("%w: %s", ErrTargetExists, target)
	}

	return nil
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
		return fmt.Errorf("tools: read dir %s: %w", dir, err)
	}

	if len(entries) > 0 {
		return fmt.Errorf("%w: %s", ErrTargetNotEmpty, dir)
	}

	return nil
}

func subdirs(root string) ([]string, error) {
	entries, err := templatesFS.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("tools: read templates %s: %w", root, err)
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
