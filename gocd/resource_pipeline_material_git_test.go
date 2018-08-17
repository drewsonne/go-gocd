package gocd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func testUnmarshallMaterialAttributesGit(t *testing.T) {
	expected := MaterialAttributesGit{
		Name:   "test-name",
		URL:    "test-url",
		Branch: "test-branch",

		SubmoduleFolder: "test-submodule_folder",
		ShallowClone:    true,

		Destination: "test-destination",
		Filter: &MaterialFilter{
			Ignore: []string{"one", "two"},
		},
		InvertFilter: true,
		AutoUpdate:   true,
	}

	m := MaterialAttributesGit{}
	unmarshallMaterialAttributesGit(&m, map[string]interface{}{
		"name":             "test-name",
		"url":              expected.URL,
		"auto_update":      expected.AutoUpdate,
		"branch":           expected.Branch,
		"submodule_folder": expected.SubmoduleFolder,
		"destination":      expected.Destination,
		"shallow_clone":    expected.ShallowClone,
		"invert_filter":    expected.InvertFilter,
		"filter": map[string]interface{}{
			"ignore": expected.Filter.Ignore,
		},
		"foo": nil,
	})

	assert.Equal(t, expected, m)
}

func testGenerateGenericGitDependency(t *testing.T) {
	for _, tt := range []struct {
		name           string
		dependency     *MaterialAttributesGit
		dependencyWant map[string]interface{}
	}{
		{
			name: "basic",
			dependency: &MaterialAttributesGit{
				Name:            "mock-name",
				URL:             "mock-url",
				Branch:          "mock-branch",
				SubmoduleFolder: "mock-folder",
				ShallowClone:    true,
				Destination:     "mock-destination",
				InvertFilter:    true,
				AutoUpdate:      true,
				Filter: &MaterialFilter{
					Ignore: []string{"one", "two"},
				},
			},
			dependencyWant: map[string]interface{}{
				"name":             "mock-name",
				"url":              "mock-url",
				"branch":           "mock-branch",
				"submodule_folder": "mock-folder",
				"shallow_clone":    true,
				"destination":      "mock-destination",
				"invert_filter":    true,
				"auto_update":      true,
				"filter": map[string]interface{}{
					"ignore": []interface{}{"one", "two"},
				},
			},
		},
		{
			name:           "null",
			dependency:     &MaterialAttributesGit{},
			dependencyWant: map[string]interface{}{},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.dependency.GenerateGeneric()
			assert.Equal(t, tt.dependencyWant, got)
		})
	}
}
