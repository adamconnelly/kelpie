package main

import (
	_ "embed"
	"fmt"
	"os"
	"text/template"

	"github.com/adamconnelly/kelpie/parser"
)

//go:embed "mock.go.tmpl"
var mockTemplate string

func main() {
	file, err := os.Open("../parser/parser.go")
	if err != nil {
		fmt.Printf("Could not open file for parsing: %v", err)
		return
	}

	filter := parser.IncludingInterfaceFilter{
		InterfacesToInclude: []string{"github.com/adamconnelly/kelpie/parser.InterfaceFilter"},
	}

	mockedInterfaces, err := parser.Parse(file, "github.com/adamconnelly/kelpie/parser", &filter)
	if err != nil {
		fmt.Printf("Could not parse file: %v", err)
		return
	}

	template := template.Must(template.New("mock").Parse(mockTemplate))

	for _, i := range mockedInterfaces {
		if _, err := os.Stat(fmt.Sprintf("out/%s", i.PackageName)); os.IsNotExist(err) {
			os.Mkdir(fmt.Sprintf("out/%s", i.PackageName), 0700)
		}
		file, err := os.Create(fmt.Sprintf("out/%s/%s.go", i.PackageName, i.PackageName))
		if err != nil {
			fmt.Printf("Could not generate file: %v", err)
			return
		}
		defer file.Close()

		if err := template.Execute(file, i); err != nil {
			fmt.Printf("OH NO!!! %v\n", err)
			return
		}
	}
}
