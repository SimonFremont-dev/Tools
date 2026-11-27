package main

import (
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

func main() {
	log.SetFlags(0)
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "list":
		err = cmdList()
	case "new":
		err = cmdNew(os.Args[2:])
	case "apply":
		err = cmdApply(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Print(`tools - scaffold projects from the templates in this repo

Usage:
  tools list
  tools new <project> --lang go [--module PATH] [--dir DIR] [--stages dev,quality] [--git] [--force]
  tools apply <stage> --dir DIR [--lang go] [--project NAME] [--force]

Environment:
  TOOLS_LANG   default language for --lang (default "go")
`)
}

func cmdList() error {
	stages, err := tools.Stages()
	if err != nil {
		return err
	}
	langs, err := tools.Langs()
	if err != nil {
		return err
	}
	fmt.Println("stages:              ", strings.Join(stages, ", "))
	fmt.Println("bootstrap languages: ", strings.Join(langs, ", "))
	return nil
}

func cmdNew(args []string) error {
	fs := flag.NewFlagSet("new", flag.ExitOnError)
	lang := fs.String("lang", env.Get("TOOLS_LANG", "go"), "language template")
	module := fs.String("module", "", "module path (default: project name)")
	dir := fs.String("dir", "", "target directory (default: ./<project>)")
	stages := fs.String("stages", "", "comma-separated extra stages (dev,quality,...)")
	gitInit := fs.Bool("git", false, "run git init in the new project")
	force := fs.Bool("force", false, "overwrite existing files")
	positional, err := parseArgs(fs, args)
	if err != nil {
		return fmt.Errorf("new: %w", err)
	}
	if len(positional) != 1 {
		return fmt.Errorf("new: expected exactly one project name")
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
	vars := tools.DefaultVars(project, *module, *lang)
	if err := tools.New(project, *lang, extra, target, vars, *force); err != nil {
		return err
	}
	fmt.Println("created", target)
	if *gitInit {
		return command.RunIn(target, "git", "init")
	}
	return nil
}

func cmdApply(args []string) error {
	fs := flag.NewFlagSet("apply", flag.ExitOnError)
	dir := fs.String("dir", ".", "project directory")
	lang := fs.String("lang", env.Get("TOOLS_LANG", "go"), "language template")
	project := fs.String("project", "", "project name for placeholders (default: base name of --dir)")
	force := fs.Bool("force", false, "overwrite existing files")
	positional, err := parseArgs(fs, args)
	if err != nil {
		return fmt.Errorf("apply: %w", err)
	}
	if len(positional) != 1 {
		return fmt.Errorf("apply: expected exactly one stage")
	}
	if *project == "" {
		abs, err := filepath.Abs(*dir)
		if err != nil {
			return err
		}
		*project = filepath.Base(abs)
	}
	vars := tools.DefaultVars(*project, *project, *lang)
	return tools.Apply(positional[0], *lang, *dir, vars, *force)
}

// parseArgs parses flags interspersed with positional arguments and returns
// the positional arguments in order.
func parseArgs(fs *flag.FlagSet, args []string) ([]string, error) {
	var positional []string
	for len(args) > 0 {
		if err := fs.Parse(args); err != nil {
			return nil, err
		}
		rest := fs.Args()
		if len(rest) == 0 {
			return positional, nil
		}
		positional = append(positional, rest[0])
		args = rest[1:]
	}
	return positional, nil
}
