package generator

type GeneratorConfig struct {
	ServerURL        string
	InputFilePath    string
	OutputDirectory  string
	Module           string
	CommandName      string
	WithModel        bool
	WithCompilerFile bool
}

func NewGeneratorConfig() *GeneratorConfig {
	return &GeneratorConfig{
		ServerURL:       "",
		InputFilePath:   "",
		OutputDirectory: "out",
		Module:          "",
		CommandName:     "",
		WithModel:       false,
	}
}
