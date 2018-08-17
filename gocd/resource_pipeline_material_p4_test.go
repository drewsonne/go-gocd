package gocd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func testUnmarshallMaterialAttributesP4(t *testing.T) {
	m := MaterialAttributesP4{}
	expected := MaterialAttributesP4{
		Name:       "test-name",
		Port:       "test-port",
		UseTickets: true,
		View:       "test-view",

		Username:          "test-user",
		Password:          "test-pass",
		EncryptedPassword: "test-encryptedpass",

		Destination: "test-dest",
		Filter: &MaterialFilter{
			Ignore: []string{"one", "two"},
		},
		InvertFilter: true,
		AutoUpdate:   true,
	}
	unmarshallMaterialAttributesP4(&m, map[string]interface{}{
		"name":               expected.Name,
		"port":               expected.Port,
		"use_tickets":        expected.UseTickets,
		"view":               expected.View,
		"username":           expected.Username,
		"password":           expected.Password,
		"encrypted_password": expected.EncryptedPassword,
		"destination":        expected.Destination,
		"filter": map[string]interface{}{
			"ignore": expected.Filter.Ignore,
		},
		"auto_update":   expected.AutoUpdate,
		"invert_filter": expected.InvertFilter,
		"foo":           nil,
	})

	assert.Equal(t, expected, m)
}

func testGenerateGenericP4Dependency(t *testing.T) {
	for _, tt := range []struct {
		name           string
		dependency     *MaterialAttributesP4
		dependencyWant map[string]interface{}
	}{
		{
			name: "basic",
			dependency: &MaterialAttributesP4{
				Name:              "mock-name",
				Port:              "mock-port",
				UseTickets:        true,
				View:              "mock-view",
				Username:          "mock-username",
				Password:          "mock-password",
				EncryptedPassword: "mock-enc-password",
				Destination:       "mock-destination",
				InvertFilter:      true,
				AutoUpdate:        true,
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
			dependency:     &MaterialAttributesP4{},
			dependencyWant: map[string]interface{}{},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.dependency.GenerateGeneric()
			assert.Equal(t, tt.dependencyWant, got)
		})
	}
}
