// Package main contains the Kelpie code generator.
package main

import (
	_ "embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"

	"github.com/alecthomas/kong"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/adamconnelly/kelpie/parser"
	"github.com/adamconnelly/kelpie/slices"
)

//go:embed "mock.go.tmpl"
var mockTemplate string

type generateCmd struct {
	ConfigFile string   `name:"config-file" short:"c" help:"The path to Kelpie's configuration file."`
	Package    string   `name:"package" short:"p" help:"The Go package containing the interface to mock."`
	Interfaces []string `name:"interfaces" short:"i" help:"The names of the interfaces to mock."`
	OutputDir  string   `name:"output-dir" short:"o" default:"mocks" help:"The directory to write the mock out to."`
}

func (g *generateCmd) Run() (err error) {
	if g.ConfigFile != "" && (g.Package != "" || len(g.Interfaces) > 0 || g.OutputDir != "") {
		return errors.New("please either specify a Kelpie config file, or specify the -package, -interfaces and -output-dir options, but not both")
	}

	var config Config

	if g.Package == "" {
		configFile, err := g.tryOpenConfigFile(g.ConfigFile)
		if err != nil {
			return err
		}
		defer configFile.Close()

		if err = yaml.NewDecoder(configFile).Decode(&config); err != nil {
			return errors.Wrap(err, fmt.Sprintf("could not parse Kelpie's config file at '%s'", configFile.Name()))
		}

		if config.Version != ConfigVersion1 {
			return fmt.Errorf("the only supported config version is '1', but '%s' was specified in the config file", config.Version)
		}
	} else {
		config.Packages = []PackageConfig{
			{
				PackageName:     g.Package,
				OutputDirectory: g.OutputDir,
				Mocks: slices.Map(g.Interfaces, func(interfaceName string) MockConfig {
					return MockConfig{
						InterfaceName: interfaceName,
					}
				}),
			},
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "could not get current working directory")
	}

	fmt.Printf("Kelpie mock generation starting - preparing to add some magic to your code-base!\n\n")

	for _, pkg := range config.Packages {
		if err = g.generatePackageMocks(cwd, pkg); err != nil {
			return err
		}
	}

	fmt.Printf("Mock generation complete!\n")

	return nil
}

var defaultConfigFiles = []string{"kelpie.yaml", "kelpie.yml"}

func (g *generateCmd) tryOpenConfigFile(customFilename string) (*os.File, error) {
	filenames := defaultConfigFiles
	if customFilename != "" {
		filenames = []string{customFilename}
	}

	for _, filename := range filenames {
		// #nosec G304 -- We're opening a potentially user-supplied filename, so we have to pass the filename via a variable.
		file, err := os.Open(filename)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}

			return nil, errors.Wrap(err, fmt.Sprintf("could not open Kelpie configuration file at '%s'", filename))
		}

		return file, nil
	}

	if customFilename == "" {
		return nil, fmt.Errorf("could not find Kelpie config file in any of the following locations: [%s]", strings.Join(filenames, ", "))
	}

	return nil, fmt.Errorf("could not find Kelpie config file in '%s'", customFilename)
}

func (g *generateCmd) generatePackageMocks(cwd string, pkg PackageConfig) error {
	filter := parser.IncludingInterfaceFilter{
		InterfacesToInclude: slices.Map(pkg.Mocks, func(m MockConfig) string { return m.InterfaceName }),
	}

	fmt.Printf("Parsing package '%s' for interfaces to mock.\n", pkg.PackageName)

	parsedPackage, err := parser.Parse(pkg.PackageName, cwd, &filter)
	if err != nil {
		return errors.Wrap(err, "could not parse file")
	}

	template := template.Must(template.New("mock").
		Funcs(template.FuncMap{
			"CommentBlock": func(comment string) string {
				lines := strings.Split(comment, "\n")
				return strings.Join(slices.Map(lines, func(line string) string {
					if line == "" {
						return "//"
					}

					return "// " + line
				}), "\n")
			},
			"Unexport": func(name string) string {
				firstRune, size := utf8.DecodeRuneInString(name)
				if firstRune == utf8.RuneError && size <= 1 {
					return name
				}

				lower := unicode.ToLower(firstRune)
				if firstRune == lower {
					return name
				}

				return string(lower) + name[size:]
			},
		}).
		Parse(mockTemplate))

	baseOutputDirectory := pkg.OutputDirectory
	if baseOutputDirectory == "" {
		baseOutputDirectory = filepath.Join(parsedPackage.PackageDirectory, "mocks")
	}

	for _, i := range parsedPackage.Mocks {
		fmt.Printf("  - Generating a mock for '%s'.\n", i.Name)

		mockConfig := slices.FirstOrPanic(pkg.Mocks, func(m MockConfig) bool { return m.InterfaceName == i.FullName })
		if mockConfig.GenerationOptions.PackageName != "" {
			i.PackageName = mockConfig.GenerationOptions.PackageName
		}

		err := func() error {
			outputDirectoryName := filepath.Join(baseOutputDirectory, i.PackageName)
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

	fmt.Println()

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
