package filter

import (
	"github.com/pkg/errors"
)

var (
	ErrNotSetNextFilter = errors.Errorf("chain filter doesn't have next filter")
)
