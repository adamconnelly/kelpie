// Package main contains the Kelpie code generator.
package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/alecthomas/kong"
	"github.com/pkg/errors"

	"github.com/adamconnelly/kelpie/parser"
	"github.com/adamconnelly/kelpie/slices"
)

//go:embed "mock.go.tmpl"
var mockTemplate string

type generateCmd struct {
	Package    string   `short:"p" required:"" help:"The Go package containing the interface to mock."`
	Interfaces []string `short:"i" required:"" help:"The names of the interfaces to mock."`
	OutputDir  string   `short:"o" required:"" default:"mocks" help:"The directory to write the mock out to."`
}

func (g *generateCmd) Run() error {
	filter := parser.IncludingInterfaceFilter{
		InterfacesToInclude: g.Interfaces,
	}

	cwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "could not get current working directory")
	}

	mockedInterfaces, err := parser.Parse(g.Package, cwd, &filter)
	if err != nil {
		return errors.Wrap(err, "could not parse file")
	}

	template := template.Must(template.New("mock").
		Funcs(template.FuncMap{
			"CommentBlock": func(comment string) string {
				lines := strings.Split(comment, "\n")
				return strings.Join(slices.Map(lines, func(line string) string {
					return "// " + line
				}), "\n")
			},
		}).
		Parse(mockTemplate))

	for _, i := range mockedInterfaces {
		err := func() error {
			outputDirectoryName := filepath.Join(g.OutputDir, i.PackageName)
			if _, err := os.Stat(outputDirectoryName); os.IsNotExist(err) {
				if err := os.MkdirAll(outputDirectoryName, 0700); err != nil {
					return errors.Wrap(err, "could not create directory for mock")
				}
			}
			file, err := os.Create(filepath.Clean(filepath.Join(outputDirectoryName, fmt.Sprintf("%s.go", i.PackageName))))
			if err != nil {
				return errors.Wrap(err, "could not open output file")
			}
			defer file.Close()

			if err := template.Execute(file, i); err != nil {
				return errors.Wrap(err, "could not generate mock")
			}

			return nil
		}()

		if err != nil {
			return err
		}
	}

	return nil
}

var cli struct {
	Generate generateCmd `cmd:"" help:"Generate a mock."`
}

func main() {
	ctx := kong.Parse(&cli, kong.Name("kelpie"), kong.Description("A magical tool for generating Go mocks!"))
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
