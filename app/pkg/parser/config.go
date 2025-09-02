package parser

import "github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"

type Config struct {
	ParserCodeGenConf *codegen.Configuration
	InputFilePath     string
}
