package gocd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func testUnmarshallMaterialAttributesDependency(t *testing.T) {
	m := MaterialAttributesDependency{}
	unmarshallMaterialAttributesDependency(&m, map[string]interface{}{
		"name":        "test-name",
		"pipeline":    "test-pipeline",
		"stage":       "test-stage",
		"foo":         nil,
		"auto_update": true,
	})

	assert.Equal(t, MaterialAttributesDependency{
		Name:       "test-name",
		Pipeline:   "test-pipeline",
		Stage:      "test-stage",
		AutoUpdate: true,
	}, m)
}

func testGenerateGenericMaterialDependency(t *testing.T) {
	for _, tt := range []struct {
		name           string
		dependency     *MaterialAttributesDependency
		dependencyWant map[string]interface{}
	}{
		{
			name: "basic",
			dependency: &MaterialAttributesDependency{
				Name:       "mock-name",
				Pipeline:   "mock-pipeline",
				Stage:      "mock-stage",
				AutoUpdate: true,
			},
			dependencyWant: map[string]interface{}{
				"name":        "mock-name",
				"pipeline":    "mock-pipeline",
				"stage":       "mock-stage",
				"auto_update": true,
			},
		},
		{
			name:           "null",
			dependency:     &MaterialAttributesDependency{},
			dependencyWant: map[string]interface{}{},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.dependency.GenerateGeneric()
			assert.Equal(t, tt.dependencyWant, got)
		})
	}
}
