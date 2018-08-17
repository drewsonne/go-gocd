package gocd

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func testResourceMaterial(t *testing.T) {
	t.Run("Generic", testMaterialAttributeGeneric)
	t.Run("Equality", testMaterialEquality)
	t.Run("AttributeEquality", testMaterialAttributeEquality)
	t.Run("AttributeInequality", testMaterialAttributeInequality)
	t.Run("HasFilter", testMaterialAttributeFilterable)
	t.Run("Unmarshall", testMaterialUnmarshall)
	t.Run("UnmarshallAttributes", testMaterialUnmarshallAttributes)
	t.Run("DecodeConfigStringList", testDecodeConfigStringList)
}

func testMaterialAttributeGeneric(t *testing.T) {
	for i, test := range []struct {
		a MaterialAttribute
		m map[string]interface{}
	}{
		{
			a: MaterialAttributesGit{
				Name:            "mock-name",
				URL:             "mock-url",
				AutoUpdate:      true,
				Branch:          "mock-branch",
				SubmoduleFolder: "mock-folder",
				Destination:     "mock-destination",
				ShallowClone:    true,
				InvertFilter:    true,
			},
			m: map[string]interface{}{
				"name":             "mock-name",
				"url":              "mock-url",
				"auto_update":      true,
				"branch":           "mock-branch",
				"submodule_folder": "mock-folder",
				"destination":      "mock-destination",
				"shallow_clone":    true,
				"invert_filter":    true,
			},
		},
		{
			a: MaterialAttributesSvn{
				Name:              "mock-name",
				URL:               "mock-url",
				AutoUpdate:        true,
				Username:          "mock-username",
				Password:          "mock-password",
				EncryptedPassword: "mock-encrypted-password",
				CheckExternals:    true,
				Destination:       "mock-destination",
				Filter: &MaterialFilter{
					Ignore: []string{"mock-ignore"},
				},
				InvertFilter: true,
			},
			m: map[string]interface{}{
				"name":               "mock-name",
				"url":                "mock-url",
				"auto_update":        true,
				"username":           "mock-username",
				"password":           "mock-password",
				"encrypted_password": "mock-encrypted-password",
				"check_externals":    true,
				"destination":        "mock-destination",
				"filter": map[string]interface{}{
					"ignore": []interface{}{"mock-ignore"},
				},
				"invert_filter": true,
			},
		},
		{
			a: MaterialAttributesHg{
				Name:        "mock-name",
				URL:         "mock-url",
				Destination: "mock-destination",
				Filter: &MaterialFilter{
					Ignore: []string{"mock-ignore"},
				},
				InvertFilter: true,
				AutoUpdate:   true,
			},
			m: map[string]interface{}{
				"name":        "mock-name",
				"url":         "mock-url",
				"auto_update": true,
				"destination": "mock-destination",
				"filter": map[string]interface{}{
					"ignore": []interface{}{"mock-ignore"},
				},
				"invert_filter": true,
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, test.m, test.a.GenerateGeneric())
		})
	}
}

func testMaterialAttributeFilterable(t *testing.T) {
	for i, test := range []struct {
		a      MaterialAttribute
		result bool
	}{
		{a: MaterialAttributesGit{}, result: true},
		{a: MaterialAttributesSvn{}, result: true},
		{a: MaterialAttributesHg{}, result: true},
		{a: MaterialAttributesP4{}, result: true},
		{a: MaterialAttributesTfs{}, result: true},
		{a: MaterialAttributesDependency{}, result: false},
		{a: MaterialAttributesPackage{}, result: false},
		{a: MaterialAttributesPlugin{}, result: true},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, test.result, test.a.HasFilter())
			if test.result {
				assert.IsType(t, &MaterialFilter{}, test.a.GetFilter())
			} else {
				assert.Nil(t, test.a.GetFilter())
			}
		})
	}
}

func testMaterialEquality(t *testing.T) {
	s1 := Material{
		Type: "git",
		Attributes: MaterialAttributesGit{
			URL: "https://github.com/gocd/gocd",
		},
	}

	s2 := Material{
		Type: "git",
		Attributes: MaterialAttributesGit{
			Name: "gocd-src",
			URL:  "https://github.com/gocd/gocd",
		},
	}
	ok, err := s1.Equal(&s2)
	assert.Nil(t, err)
	assert.True(t, ok)
}

func testMaterialAttributeEquality(t *testing.T) {
	for i, test := range []struct {
		a MaterialAttribute
		b MaterialAttribute
	}{
		{a: MaterialAttributesGit{}, b: MaterialAttributesGit{}},
		{a: MaterialAttributesGit{Branch: ""}, b: MaterialAttributesGit{Branch: "master"}},
		{a: MaterialAttributesGit{Branch: "master"}, b: MaterialAttributesGit{Branch: ""}},
		{a: MaterialAttributesGit{Branch: ""}, b: MaterialAttributesGit{Branch: ""}},
		{a: MaterialAttributesGit{Branch: "master"}, b: MaterialAttributesGit{Branch: "master"}},
		{a: MaterialAttributesSvn{}, b: MaterialAttributesSvn{}},
		{a: MaterialAttributesHg{}, b: MaterialAttributesHg{}},
		{a: MaterialAttributesP4{}, b: MaterialAttributesP4{}},
		{a: MaterialAttributesTfs{}, b: MaterialAttributesTfs{}},
		{a: MaterialAttributesDependency{}, b: MaterialAttributesDependency{}},
		{a: MaterialAttributesPackage{}, b: MaterialAttributesPackage{}},
		{a: MaterialAttributesPlugin{}, b: MaterialAttributesPlugin{}},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ok, err := test.a.equal(test.b)
			assert.True(t, ok)
			assert.Nil(t, err)
		})
	}
}

func testMaterialAttributeInequality(t *testing.T) {
	for i, test := range []struct {
		a         MaterialAttribute
		b         MaterialAttribute
		errString string
	}{
		{a: MaterialAttributesGit{}, b: MaterialAttributesP4{}, errString: "can only compare with same material type"},
		{a: MaterialAttributesSvn{}, b: MaterialAttributesGit{}, errString: "can only compare with same material type"},
		{a: MaterialAttributesHg{}, b: MaterialAttributesGit{}, errString: "can only compare with same material type"},
		{a: MaterialAttributesP4{}, b: MaterialAttributesGit{}, errString: "can only compare with same material type"},
		{a: MaterialAttributesTfs{}, b: MaterialAttributesGit{}, errString: "can only compare with same material type"},
		{a: MaterialAttributesDependency{}, b: MaterialAttributesGit{}, errString: "can only compare with same material type"},
		{a: MaterialAttributesPackage{}, b: MaterialAttributesGit{}, errString: "can only compare with same material type"},
		{a: MaterialAttributesPlugin{}, b: MaterialAttributesGit{}, errString: "can only compare with same material type"},
		{a: MaterialAttributesGit{}, b: MaterialAttributesGit{URL: "https://github.com/gocd/gocd"}},
		{
			a: MaterialAttributesGit{URL: "https://github.com/gocd/gocd"},
			b: MaterialAttributesGit{URL: "https://github.com/gocd/gocd", Branch: "feature/branch"},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ok, err := test.a.equal(test.b)
			assert.False(t, ok)
			if test.errString == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.errString)
			}
		})

	}
}

func testMaterialUnmarshall(t *testing.T) {
	m := Material{}
	for i, test := range []struct {
		source   string
		expected MaterialAttribute
	}{
		{source: `{"type": "git"}`, expected: &MaterialAttributesGit{}},
		{source: `{"type": "svn"}`, expected: &MaterialAttributesSvn{}},
		{source: `{"type": "hg"}`, expected: &MaterialAttributesHg{}},
		{source: `{"type": "p4"}`, expected: &MaterialAttributesP4{}},
		{source: `{"type": "tfs"}`, expected: &MaterialAttributesTfs{}},
		{source: `{"type": "dependency"}`, expected: &MaterialAttributesDependency{}},
		{source: `{"type": "package"}`, expected: &MaterialAttributesPackage{}},
		{source: `{"type": "plugin"}`, expected: &MaterialAttributesPlugin{}},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			m.UnmarshalJSON([]byte(test.source))
			assert.IsType(t, test.expected, m.Attributes)
		})
	}
}

func testMaterialUnmarshallAttributes(t *testing.T) {
	t.Run("Dependency", testUnmarshallMaterialAttributesDependency)
	t.Run("DependencyGeneric", testGenerateGenericMaterialDependency)
	t.Run("Git", testUnmarshallMaterialAttributesGit)
	t.Run("GitGeneric", testGenerateGenericGitDependency)
	t.Run("Hg", testUnmarshallMaterialAttributesHg)
	t.Run("P4", testUnmarshallMaterialAttributesP4)
	t.Run("P4Generic", testGenerateGenericP4Dependency)
	t.Run("Package", testUnmarshallMaterialAttributesPkg)
	t.Run("Plugin", testUnmarshallMaterialAttributesPlugin)
	t.Run("SVN", testUnmarshallMaterialAttributesSvn)
	t.Run("TFS", testUnmarshallMaterialAttributesTfs)
}

func testDecodeConfigStringList(t *testing.T) {
	type args struct {
		lI []interface{}
	}
	for _, tt := range []struct {
		name    string
		args    args
		wantRet []string
	}{
		{
			name: "basic",
			args: args{
				lI: []interface{}{
					"one",
				},
			},
			wantRet: []string{"one"},
		},
		{
			name: "multi",
			args: args{
				lI: []interface{}{
					"one",
					"two",
				},
			},
			wantRet: []string{"one", "two"},
		},
	} {
		got := decodeConfigStringList(tt.args.lI)
		assert.Equal(t, tt.wantRet, got)
	}
}
