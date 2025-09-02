package generator

import "github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"

type GeneratorConfig struct {
	ServerURL        string
	OutputDirectory  string
	Module           string
	CommandName      string
	WithModel        bool
	WithCompilerFile bool

	parserCodeGenConf *codegen.Configuration
}

func NewGeneratorConfig(parserCodeGenConf *codegen.Configuration) *GeneratorConfig {
	return &GeneratorConfig{
		ServerURL:         "",
		OutputDirectory:   "out",
		Module:            "",
		CommandName:       "",
		WithModel:         false,
		parserCodeGenConf: parserCodeGenConf,
	}
}
