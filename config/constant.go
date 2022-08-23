package config

import (
	"regexp"

	set "github.com/deckarep/golang-set/v2"
)

const (
	DEFAULT_PATH = "config/glory.yaml"
	GLORY_ENV    = "GLORY_ENV"

	CONFIG_CENTER_KEY = "config_center"
)

var (
	defaultConfigPath = DEFAULT_PATH

	keepedConfigCenterName   = set.NewSet("", "env")
	skipInitConfigCenterName = set.NewSet("env")
	keepedComponentName      = set.NewSet("", "config_center", "service_name")

	placeHolderRegexp = regexp.MustCompile(`^\$(.*)\{(.+)\}$`)
)
