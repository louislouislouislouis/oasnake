package config

type GenerateConfig struct {
	ServerURL       string
	InputFilePath   string
	OutputDirectory string
	Module          string
	CommandName     string
	WithModel       bool
	WithMainFile    bool
	Installation    InstallationConfig
}

func NewGenerateConfig() GenerateConfig {
	return GenerateConfig{
		ServerURL:       "",
		InputFilePath:   "",
		OutputDirectory: "out",
		Module:          "",
		CommandName:     "",
		WithModel:       false,
		WithMainFile:    true,
		Installation: InstallationConfig{
			BinaryName:        "",
			InstallWithGo:     true,
			InstallWithDocker: false,
		},
	}
}

type InstallationConfig struct {
	HasToBeInstalled  bool // Indicates if the binary needs to be installed
	BinaryName        string
	InstallWithGo     bool
	InstallWithDocker bool
}

func (cfg InstallationConfig) NeedToBeInstalled() bool {
	return cfg.HasToBeInstalled || cfg.BinaryName != ""
}
