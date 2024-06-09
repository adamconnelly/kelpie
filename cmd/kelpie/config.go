package main

// ConfigVersion defines the version of Kelpie's config file.
type ConfigVersion string

const (
	// ConfigVersion1 is v1 of the config file.
	ConfigVersion1 ConfigVersion = "1"
)

// Config represents Kelpie's generation config.
type Config struct {
	// Version is the version of the config file.
	Version ConfigVersion

	// Packages contains the configuration of the packages to generate mocks from.
	Packages []PackageConfig
}

// PackageConfig defines the configuration for a single package to mock.
type PackageConfig struct {
	// PackageName is the full path of the package, for example github.com/adamconnelly/kelpie/examples.
	PackageName string `yaml:"package"`

	// Mocks contains the list of interfaces to mock.
	Mocks []MockConfig

	// Output directory is the directory to output generated mocks for this package. Defaults
	// to a folder called "mocks" in the package directory if not specified.
	OutputDirectory string `yaml:"directory"`
}

// MockConfig is configuration of an individual mock.
type MockConfig struct {
	// InterfaceName is the name of the interface to mock, for example "Maths" or "SomeService.SomeRepository".
	InterfaceName string `yaml:"interface"`

	// GenerationOptions allows generation of the mock to be customized.
	GenerationOptions MockGenerationOptions `yaml:"generation"`
}

// MockGenerationOptions allows generation of the mock to be customized.
type MockGenerationOptions struct {
	// PackageName is the name of the generated package for the mock. Defaults to the lowercased
	// interface name. For example an interface called EmailSender would generate a package
	// called emailsender.
	PackageName string `yaml:"package"`
}
