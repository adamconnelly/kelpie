package main

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type MockedInterface struct {
	Name        string
	PackageName string
	Methods     []MethodDefinition
}

type MethodDefinition struct {
	Name       string
	Parameters []ParameterDefinition
	Returns    []ReturnDefinition
}

type ParameterDefinition struct {
	Name string
	Type string
}

type ReturnDefinition struct {
	Name string
	Type string
}

//go:embed "mock.go.tmpl"
var mockTemplate string

func main() {
	mockedInterface := MockedInterface{
		Name:        "Maths",
		PackageName: strings.ToLower("Maths"),
		Methods: []MethodDefinition{
			{
				Name: "Add",
				Parameters: []ParameterDefinition{
					{
						Name: "a",
						Type: "int",
					},
					{
						Name: "b",
						Type: "int",
					},
				},
				Returns: []ReturnDefinition{
					{
						Name: "",
						Type: "int",
					},
				},
			},
			{
				Name: "ParseInt",
				Parameters: []ParameterDefinition{
					{
						Name: "input",
						Type: "string",
					},
				},
				Returns: []ReturnDefinition{
					{
						Name: "",
						Type: "int",
					},
					{
						Name: "",
						Type: "error",
					},
				},
			},
		},
	}

	template := template.Must(template.New("mock").Parse(mockTemplate))

	file, err := os.Create("out/generated.go")
	if err != nil {
		fmt.Printf("Could not generate file: %v", err)
		return
	}
	defer file.Close()

	if err := template.Execute(file, mockedInterface); err != nil {
		fmt.Printf("OH NO!!! %v\n", err)
		return
	}
}
