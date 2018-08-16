package gocd

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
)

func testServerVersionResource(t *testing.T) {
	t.Run("LessThan", testServerVersionLessThan)
	t.Run("Equal", testServerVersionEqual)
	t.Run("GetAPIVersion", testServerVersionGetAPIVersion)
	t.Run("GetAPIVersionFail", testServerVersionGetAPIVersionFail)
}

func testServerVersionEqual(t *testing.T) {
	for _, test := range []struct {
		v1   *ServerVersion
		v2   *ServerVersion
		want bool
	}{
		{v1: &ServerVersion{Version: "1.2.3"}, v2: &ServerVersion{Version: "1.2.3"}, want: true},
		{v1: &ServerVersion{Version: "1.2.3"}, v2: &ServerVersion{Version: "2.2.3"}, want: false},
	} {
		assert.Equal(t, test.want, test.v1.Equal(test.v2))
		assert.Equal(t, test.want, test.v2.Equal(test.v1))
	}
}

func testServerVersionLessThan(t *testing.T) {
	for _, test := range []struct {
		v1   *ServerVersion
		v2   *ServerVersion
		want bool
	}{
		{v1: &ServerVersion{Version: "1.0.0"}, v2: &ServerVersion{Version: "2.0.0"}, want: true},
		{v1: &ServerVersion{Version: "2.0.1"}, v2: &ServerVersion{Version: "2.0.0"}, want: false},
		{v1: &ServerVersion{Version: "2.0.0"}, v2: &ServerVersion{Version: "2.0.1"}, want: true},
		{v1: &ServerVersion{Version: "2.0.0"}, v2: &ServerVersion{Version: "1.0.0"}, want: false},
	} {
		name := fmt.Sprintf("%s < %s = %t", test.v1.Version, test.v2.Version, test.want)
		t.Run(name, func(t *testing.T) {

			test.v1.parseVersion()
			test.v2.parseVersion()

			assert.Equal(t, test.want, test.v1.LessThan(test.v2))
			assert.Equal(t, !test.want, test.v2.LessThan(test.v1))
		})
	}
}

func testServerVersionGetAPIVersion(t *testing.T) {
	for _, test := range []struct {
		v        *ServerVersion
		endpoint string
		want     string
	}{
		{
			endpoint: "/api/version",
			want:     apiV1,
			v:        &ServerVersion{Version: "16.7.0"},
		},
		{
			endpoint: "/api/admin/pipelines/:pipeline_name",
			want:     apiV5,
			v:        &ServerVersion{Version: "17.13.0"},
		},
		{
			endpoint: "/api/admin/pipelines/:pipeline_name",
			want:     apiV5,
			v:        &ServerVersion{Version: "18.6.0"},
		},
		{
			endpoint: "/api/admin/pipelines/:pipeline_name",
			want:     apiV6,
			v:        &ServerVersion{Version: "18.8.0"},
		},
	} {
		test.v.parseVersion()
		apiV, err := test.v.GetAPIVersion(test.endpoint)

		assert.NoError(t, err)
		assert.Equal(t, test.want, apiV)
	}
}

func testServerVersionGetAPIVersionFail(t *testing.T) {
	for _, test := range []struct {
		v        *ServerVersion
		endpoint string
		want     string
	}{
		{
			endpoint: "/api/version",
			want:     "could not find api version for server version '1.0.0'",
			v:        &ServerVersion{Version: "1.0.0"},
		},
		{
			endpoint: "/api/foobar",
			want:     "could not find API version tag for '/api/foobar'",
			v:        &ServerVersion{Version: "1.0.0"},
		},
		{
			endpoint: "/api/admin/pipelines/:pipeline_name",
			want:     "could not find api version for server version '0.1.0'",
			v:        &ServerVersion{Version: "0.1.0"},
		},
	} {
		test.v.parseVersion()
		apiV, err := test.v.GetAPIVersion(test.endpoint)

		assert.EqualError(t, err, test.want)
		assert.Empty(t, apiV)
	}
}

func TestNewserverAPIVersionMapping(t *testing.T) {

	mockVersion, err := version.NewVersion("1.0.0")
	assert.NoError(t, err)
	type args struct {
		serverVersion ServerVersionString
		apiVersion    string
	}
	tests := []struct {
		name        string
		args        args
		wantMapping *serverVersionToAcceptMapping
	}{
		{
			name: "base",
			args: args{serverVersion: "1.0.0", apiVersion: apiV1},
			wantMapping: &serverVersionToAcceptMapping{
				API:    apiV1,
				Server: mockVersion,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t,
				tt.wantMapping,
				newServerAPI(tt.args.serverVersion, tt.args.apiVersion),
			)
		})
	}
}

func TestServerAPIVersionMappingCollection_Sort(t *testing.T) {
	tests := []struct {
		name string
		have *serverAPIVersionMappingCollection
		want *serverAPIVersionMappingCollection
	}{
		{
			name: "base",
			have: &serverAPIVersionMappingCollection{
				mappings: []*serverVersionToAcceptMapping{
					newServerAPI("2.0.0", apiV2),
					newServerAPI("1.0.0", apiV1),
					newServerAPI("4.0.0", apiV4),
					newServerAPI("3.0.0", apiV3),
				},
			},
			want: &serverAPIVersionMappingCollection{
				mappings: []*serverVersionToAcceptMapping{
					newServerAPI("1.0.0", apiV1),
					newServerAPI("2.0.0", apiV2),
					newServerAPI("3.0.0", apiV3),
					newServerAPI("4.0.0", apiV4),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.have.Sort()
			assert.Equal(t, tt.want, tt.have)
		})
	}
}
