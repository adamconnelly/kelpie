package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/alecthomas/kong"
	"github.com/pkg/errors"

	"github.com/adamconnelly/kelpie/parser"
)

//go:embed "mock.go.tmpl"
var mockTemplate string

type GenerateCmd struct {
	SourceFile string   `short:"s" required:"" env:"GOFILE" help:"The Go source file containing the interface to mock."`
	Interfaces []string `short:"i" required:"" help:"The names of the interfaces to mock."`
	OutputDir  string   `short:"o" required:"" default:"mocks" help:"The directory to write the mock out to."`
}

func (g *GenerateCmd) Run() error {
	file, err := os.Open(g.SourceFile)
	if err != nil {
		return errors.Wrap(err, "could not open file for parsing")
	}

	filter := parser.IncludingInterfaceFilter{
		InterfacesToInclude: g.Interfaces,
	}

	mockedInterfaces, err := parser.Parse(file, &filter)
	if err != nil {
		return errors.Wrap(err, "could not parse file")
	}

	template := template.Must(template.New("mock").Parse(mockTemplate))

	for _, i := range mockedInterfaces {
		err := func() error {
			outputDirectoryName := filepath.Join(g.OutputDir, i.PackageName)
			if _, err := os.Stat(outputDirectoryName); os.IsNotExist(err) {
				os.MkdirAll(outputDirectoryName, 0700)
			}
			file, err := os.Create(filepath.Join(outputDirectoryName, fmt.Sprintf("%s.go", i.PackageName)))
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
	Generate GenerateCmd `cmd:"" help:"Generate a mock."`
}

func main() {
	ctx := kong.Parse(&cli, kong.Name("kelpie"), kong.Description("A magical tool for generating Go mocks!"))
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
