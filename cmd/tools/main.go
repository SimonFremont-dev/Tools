// Command tools scaffolds projects from the templates in this repository.
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"tools"
	"tools/command"
	"tools/env"
)

const (
	minArgs   = 2
	exitUsage = 2
)

var (
	errNewProjectArg = errors.New("new: expected exactly one project name")
	errApplyStageArg = errors.New("apply: expected exactly one stage")
)

func main() {
	log.SetFlags(0)

	if len(os.Args) < minArgs {
		usage()
		os.Exit(exitUsage)
	}

	var err error

	switch os.Args[1] {
	case "list":
		err = cmdList()
	case "new":
		err = cmdNew(os.Args[minArgs:])
	case "apply":
		err = cmdApply(os.Args[minArgs:])
	default:
		usage()
		os.Exit(exitUsage)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Print(`tools - scaffold projects from the templates in this repo

Usage:
  tools list
  tools new <project> --lang go --db postgres [--module PATH] [--dir DIR] [--stages dev,quality] [--git] [--force]
  tools apply <stage> --dir DIR [--lang go] [--db postgres] [--project NAME] [--force]

Environment:
  TOOLS_LANG   default language for --lang (default "go")
  TOOLS_DB     default database for --db (default "postgres")
`)
}

func cmdList() error {
	stages, err := tools.Stages()
	if err != nil {
		return fmt.Errorf("list stages: %w", err)
	}

	langs, err := tools.Langs()
	if err != nil {
		return fmt.Errorf("list languages: %w", err)
	}

	dbs, err := tools.DBs()
	if err != nil {
		return fmt.Errorf("list databases: %w", err)
	}

	fmt.Println("stages:              ", strings.Join(stages, ", "))
	fmt.Println("bootstrap languages: ", strings.Join(langs, ", "))
	fmt.Println("databases:           ", strings.Join(dbs, ", "))

	return nil
}

func cmdNew(args []string) error {
	flagSet := flag.NewFlagSet("new", flag.ExitOnError)
	lang := flagSet.String("lang", env.Get("TOOLS_LANG", "go"), "language template")
	database := flagSet.String("db", env.Get("TOOLS_DB", "postgres"), "database overlay (postgres, mysql, sqlite)")
	module := flagSet.String("module", "", "module path (default: project name)")
	dir := flagSet.String("dir", "", "target directory (default: ./<project>)")
	stages := flagSet.String("stages", "", "comma-separated extra stages (dev,quality,...)")
	gitInit := flagSet.Bool("git", false, "run git init in the new project")
	force := flagSet.Bool("force", false, "overwrite existing files")

	positional, err := parseArgs(flagSet, args)
	if err != nil {
		return fmt.Errorf("new: %w", err)
	}

	if len(positional) != 1 {
		return errNewProjectArg
	}

	project := positional[0]

	target := *dir
	if target == "" {
		target = project
	}

	var extra []string
	if *stages != "" {
		extra = strings.Split(*stages, ",")
	}

	vars := tools.DefaultVars(project, *module, *lang, *database)

	err = tools.New(project, *lang, *database, extra, target, vars, *force)
	if err != nil {
		return fmt.Errorf("new: %w", err)
	}

	fmt.Println("created", target)

	if *gitInit {
		err = command.RunIn(target, "git", "init")
		if err != nil {
			return fmt.Errorf("new: git init: %w", err)
		}
	}

	return nil
}

func cmdApply(args []string) error {
	flagSet := flag.NewFlagSet("apply", flag.ExitOnError)
	dir := flagSet.String("dir", ".", "project directory")
	lang := flagSet.String("lang", env.Get("TOOLS_LANG", "go"), "language template")
	database := flagSet.String("db", env.Get("TOOLS_DB", "postgres"), "database overlay (postgres, mysql, sqlite)")
	project := flagSet.String("project", "", "project name for placeholders (default: base name of --dir)")
	force := flagSet.Bool("force", false, "overwrite existing files")

	positional, err := parseArgs(flagSet, args)
	if err != nil {
		return fmt.Errorf("apply: %w", err)
	}

	if len(positional) != 1 {
		return errApplyStageArg
	}

	if *project == "" {
		abs, err := filepath.Abs(*dir)
		if err != nil {
			return fmt.Errorf("apply: resolve dir: %w", err)
		}

		*project = filepath.Base(abs)
	}

	vars := tools.DefaultVars(*project, *project, *lang, *database)

	err = tools.Apply(positional[0], *lang, *database, *dir, vars, *force)
	if err != nil {
		return fmt.Errorf("apply %s: %w", positional[0], err)
	}

	return nil
}

// parseArgs parses flags interspersed with positional arguments and returns
// the positional arguments in order.
func parseArgs(flagSet *flag.FlagSet, args []string) ([]string, error) {
	var positional []string

	for len(args) > 0 {
		err := flagSet.Parse(args)
		if err != nil {
			return nil, fmt.Errorf("parse flags: %w", err)
		}

		rest := flagSet.Args()
		if len(rest) == 0 {
			return positional, nil
		}

		positional = append(positional, rest[0])
		args = rest[1:]
	}

	return positional, nil
}
