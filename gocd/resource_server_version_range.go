package gocd

import "errors"

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
	return lowerBound && upperBound
}
