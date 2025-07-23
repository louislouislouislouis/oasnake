package generator

import (
	"bytes"
	_ "embed"
	"text/template"

	"github.com/louislouislouislouis/oasnake/app/pkg/utils"
)

var (
	//go:embed assets/go.mod.gotmpl
	goModFileTmpl []byte

	//go:embed assets/command.gotmpl
	cobraCodeTmpl []byte

	//go:embed assets/config/command.gotmpl
	configCommand []byte

	//go:embed assets/commonCommand.gotmpl
	commonCommandTmpl []byte

	//go:embed assets/config/extension.gotmpl
	configExtension []byte

	//go:embed assets/config/method.gotmpl
	configMethod []byte

	//go:embed assets/config/request.gotmpl
	configRequest []byte

	//go:embed assets/main.gotmpl
	mainTmpl []byte

	//go:embed assets/app.gotmpl
	appTmpl []byte

	//go:embed assets/service.gotmpl
	svcTmpl []byte
)

type TemplatorType int

type Templator struct {
	t TemplatorType
}

const (
	Service TemplatorType = iota
	App
	Main
	Command
	Mod
	ConfigRequest
	ConfigMethod
	ConfigExtension
	ConfigCommand
	CommonCommand
)

func NewTemplator(t TemplatorType) *Templator {
	return &Templator{t: t}
}

func (t *Templator) getTemplate() string {
	switch t.t {
	case Service:
		return string(svcTmpl)
	case App:
		return string(appTmpl)
	case Main:
		return string(mainTmpl)
	case Command:
		return string(cobraCodeTmpl)
	case Mod:
		return string(goModFileTmpl)
	case ConfigRequest:
		return string(configRequest)
	case ConfigMethod:
		return string(configMethod)
	case ConfigExtension:
		return string(configExtension)
	case ConfigCommand:
		return string(configCommand)
	case CommonCommand:
		return string(commonCommandTmpl)
	default:
		return ""
	}
}

func (templator Templator) renderTemplate(data any) (string, error) {
	tmpl, err := template.New("template").Parse(templator.getTemplate())
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (templator Templator) WriteTemplateToFile(data any, outputPath utils.FS) error {
	renderedContent, err := templator.renderTemplate(data)
	if err != nil {
		return err
	}

	return utils.WriteFileContent(utils.WriterConfig{
		OutputDirectoryShouldBeEmpty: false,
		Output:                       outputPath,
		Content:                      renderedContent,
	})
}
