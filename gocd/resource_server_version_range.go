package gocd

import (
	"errors"
	"fmt"
)

var BadServerVersionRange = errors.New("minimum server version must not be larger than maximum server version")

type ServerVersionRange struct {
	min *ServerVersion
	max *ServerVersion
}

func NewServerVersionRange(min *ServerVersion, max *ServerVersion) (*ServerVersionRange, error) {
	min.parseVersion()
	max.parseVersion()
	if max.LessThan(min) {
		return nil, BadServerVersionRange
	}
	return &ServerVersionRange{
		min: min,
		max: max,
	}, nil
}

func (svr *ServerVersionRange) Contains(version *ServerVersion) bool {
	lowerBound := svr.min.LessThan(version)
	upperBound := version.LessThan(svr.max)
	equality := version.Equal(svr.max)
	return lowerBound && upperBound || equality
}
func (svr *ServerVersionRange) String() string {
	return fmt.Sprintf("(%s, %s]", svr.min.String(), svr.max.String())
}
