package parser

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ImportHelperTests struct {
	suite.Suite
}

func (t *ImportHelperTests) Test_RequiredImports_SortsImportsCorrectly() {
	type testCase struct {
		imports         []string
		expectedImports []string
	}

	testCases := map[string]testCase{
		"Returns imports in same order if alphabetical": {
			imports: []string{
				"github.com/adamconnelly/kelpie/examples/users",
				"github.com/adamconnelly/kelpie/mocking",
			},
			expectedImports: []string{
				"github.com/adamconnelly/kelpie/examples/users",
				"github.com/adamconnelly/kelpie/mocking",
			},
		},
		"Re-orders imports if not in alphabetical order": {
			imports: []string{
				"github.com/adamconnelly/kelpie/mocking",
				"github.com/adamconnelly/kelpie/examples/users",
			},
			expectedImports: []string{
				"github.com/adamconnelly/kelpie/examples/users",
				"github.com/adamconnelly/kelpie/mocking",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func() {
			helper := importHelper{
				requiredImports: tc.imports,
			}

			imports := helper.RequiredImports()

			t.Equal(tc.expectedImports, imports)
		})
	}
}

func TestImportHelper(t *testing.T) {
	suite.Run(t, new(ImportHelperTests))
}
