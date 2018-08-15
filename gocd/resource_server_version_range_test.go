package gocd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newServerVersion(version string) *ServerVersion {
	v := &ServerVersion{
		Version: version,
	}
	v.parseVersion()
	return v
}

func TestNewServerVersionRange(t *testing.T) {
	type args struct {
		min *ServerVersion
		max *ServerVersion
	}
	for _, tt := range []struct {
		name        string
		args        args
		want        *ServerVersionRange
		wantErr     bool
		wantErrType error
	}{
		{
			name: "basic",
			args: args{
				min: newServerVersion("1.0.0"),
				max: newServerVersion("2.0.0"),
			},
			want: &ServerVersionRange{
				min: newServerVersion("1.0.0"),
				max: newServerVersion("2.0.0"),
			},
			wantErr: false,
		},
		{
			name: "bad-range",
			args: args{
				min: newServerVersion("2.0.0"),
				max: newServerVersion("1.0.0"),
			},
			want:        nil,
			wantErr:     true,
			wantErrType: BadServerVersionRange,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewServerVersionRange(tt.args.min, tt.args.max)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErrType.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestServerVersionRange_Contains(t *testing.T) {
	type args struct {
		version *ServerVersion
	}
	for _, tt := range []struct {
		name         string
		versionRange *ServerVersionRange
		args         args
		want         bool
	}{
		{
			name:         "success",
			versionRange: newServerVersionRangeFromString("1.0.0", "3.0.0"),
			args: args{
				version: newServerVersion("2.0.0"),
			},
			want: true,
		},
		{
			name:         "fail",
			versionRange: newServerVersionRangeFromString("1.0.0", "3.0.0"),
			args: args{
				version: newServerVersion("4.0.0"),
			},
			want: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.versionRange.Contains(tt.args.version)
			assert.Equal(t, got, tt.want)
		})
	}
}

func newServerVersionRangeFromString(min string, max string) *ServerVersionRange {
	minStruct := &ServerVersion{Version: min}
	maxStruct := &ServerVersion{Version: max}

	minStruct.parseVersion()
	maxStruct.parseVersion()

	serverRange, err := NewServerVersionRange(minStruct, maxStruct)
	if err != nil {
		panic(err)
	}
	return serverRange
}
