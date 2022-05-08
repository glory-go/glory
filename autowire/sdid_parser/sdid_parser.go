package sdid_parser

import (
	"strings"
)

import (
	"github.com/glory-go/glory/autowire"
	"github.com/glory-go/glory/autowire/util"
)

type defaultSDIDParser struct {
}

var defaultSDIDParserSingleton autowire.SDIDParser

func GetDefaultSDIDParser() autowire.SDIDParser {
	if defaultSDIDParserSingleton == nil {
		defaultSDIDParserSingleton = &defaultSDIDParser{}
	}
	return defaultSDIDParserSingleton
}

func (p *defaultSDIDParser) Parse(fi *autowire.FieldInfo) (string, error) {
	splitedTagValue := strings.Split(fi.TagValue, ",")
	interfaceName := fi.FieldType
	if interfaceName == "" {
		interfaceName = splitedTagValue[0]
	}
	return util.GetIdByNamePair(interfaceName, splitedTagValue[0]), nil
}
