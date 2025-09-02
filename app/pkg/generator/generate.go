package generator

import (
	_ "embed"
	"fmt"
	"math/rand/v2"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/louislouislouislouis/oasnake/app/pkg/generator/command"
	"github.com/louislouislouislouis/oasnake/app/pkg/utils"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"github.com/rs/zerolog/log"
)

const (
	modelFileName = "model.go"

	// Default paths for project structure
	appPath     = "/app"
	commandPath = "/app/cmd"
	configPath  = "/app/pkg/config"
	servicePath = "/app/pkg/service"
	modelPath   = "/app/pkg/model"
)

// Generator is responsible for generating the CLI based on OpenAPI specs.
type Generator struct {
	Config *GeneratorConfig
}

// NewGenerator creates a new Generator instance.
func NewGenerator(cfg *GeneratorConfig) *Generator {
	return &Generator{Config: cfg}
}

// Generate executes the complete CLI generation pipeline.
// It processes the OpenAPI specification, constructs the command tree,
// renders command files, generates models, and sets up core application templates.
//
// The generation steps include:
//  1. Loading the OpenAPI spec.
//  2. Building a command tree from the API paths.
//  3. Rendering CLI command files.
//  4. Generating models via oapi-codegen.
//  5. Creating essential core application templates.
//
// Returns an error if any stage of the process fails.
func (g *Generator) Generate(rootCommand *command.NodeCmd, spec *openapi3.T) (string, error) {
	log.Debug().Msg("Starting code generation")

	err := g.addRelevantGeneratorConfig(rootCommand, spec)
	if err != nil {
		return "", fmt.Errorf("failed to add generator config: %w", err)
	}

	// Render CLI command files recursively
	cmdOutputPath := filepath.Join(g.Config.OutputDirectory, commandPath)
	if err := traverseAndRenderCommands(rootCommand, cmdOutputPath); err != nil {
		return "", fmt.Errorf("failed to render command files: %w", err)
	}

	// Generate API models using oapi-codegen
	if g.Config.WithModel {
		if err := g.generateModel(spec); err != nil {
			return "", fmt.Errorf("failed to generate API models: %w", err)
		}
	}

	// Generate core application templates
	if err := g.generateCoreApp(rootCommand); err != nil {
		return "", fmt.Errorf("failed to generate core application templates: %w", err)
	}

	log.Info().Msgf("Code generation completed successfully. Output directory: %s", g.Config.OutputDirectory)
	return g.GetEffectiveRootUsage(spec), nil
}

func (g *Generator) addRelevantGeneratorConfig(rootCommand *command.NodeCmd, spec *openapi3.T) error {
	baseURL, err := g.GetEffectiveServerURL(spec)
	if err != nil {
		return err
	}
	globalConfig := command.CommandGlobalConfig{
		RootUsage:   g.GetEffectiveRootUsage(spec),
		ModuleName:  g.Config.Module,
		BaseUrl:     baseURL,
		BaseCmdPath: g.resolvePath(commandPath),
		ConfigPath:  g.resolvePath(configPath),
		AppPath:     g.resolvePath(appPath),
		ServicePath: g.resolvePath(servicePath),
	}
	rootCommand.SetGlobalConfig(globalConfig)
	return nil
}

func (g *Generator) resolvePath(subPath string) string {
	if g.Config.WithCompilerFile {
		return subPath
	}
	return filepath.Join(g.Config.OutputDirectory, subPath)
}

// traverseAndRenderCommands recursively traverses the CLI command tree
// and renders each command node into a corresponding Go source file.
//
// For each command node:
//   - A template is rendered using the node's data.
//   - The output is written to the specified directory structure.
//   - All child nodes are processed recursively.
//
// Parameters:
//   - cmd: The root command node to render.
//   - dir: The base directory where generated files should be placed.
//
// Returns an error if rendering or traversal fails at any point.
func traverseAndRenderCommands(cmd *command.NodeCmd, dir string) error {
	templator := NewTemplator(Command)

	output := utils.FS{
		Directory: dir,
		Filename:  cmd.GetFileName(),
	}

	if err := templator.WriteTemplateToFile(cmd, output); err != nil {
		return fmt.Errorf("failed to render command template for %q: %w", cmd.GetFileName(), err)
	}

	for _, child := range cmd.Children {
		childDir := filepath.Join(dir, child.GetPackageName())

		if err := traverseAndRenderCommands(child, childDir); err != nil {
			return fmt.Errorf("failed to render subcommand %q: %w", child.GetFileName(), err)
		}
	}

	return nil
}

// GetEffectiveServerURL returns the effective server URL to be used,
// based on the configuration and the provided OpenAPI document.
//
// Priority order:
// 1️⃣ If a server URL is explicitly defined in the config (`Config.ServerURL`), it is returned.
// 2️⃣ Otherwise, it attempts to retrieve the base path from the OpenAPI `servers` section.
//
// If neither is available, an error is returned.
//
// Parameters:
//   - doc (*openapi3.T): The OpenAPI document to read the server URL from.
//
// Returns:
//   - (string): The determined server URL.
//   - (error): An error if the server URL could not be determined.
//
// Example error:
//
//	"❌ No server URL defined in the OpenAPI spec and no --server-url flag provided"
func (g *Generator) GetEffectiveServerURL(doc *openapi3.T) (string, error) {
	// 1️⃣ Use the server URL from config if set
	if g.Config.ServerURL != "" {
		return g.Config.ServerURL, nil
	}

	// 2️⃣ Attempt to extract the base path from OpenAPI spec
	url, err := doc.Servers.BasePath()
	if err != nil {
		return "", fmt.Errorf("🚨 error extracting server URL from OpenAPI spec: %w", err)
	}

	// 3️⃣ Validate that the extracted URL is meaningful
	if url == "/" {
		return "", fmt.Errorf("❌ First server URL not defined in OpenAPI spec and no --server-url flag provided")
	}

	// ✅ Return the resolved server URL
	return url, nil
}

// GetEffectiveRootUsage determines the CLI binary name to use for generation.
//
// Priority is given to the explicit value set in the configuration. If not set,
// the method attempts to infer the binary name from the OpenAPI specification's title.
//
// Parameters:
//   - doc: The OpenAPI document used to extract metadata.
//
// Returns:
//   - The effective root ne as a string.
func (g *Generator) GetEffectiveRootUsage(doc *openapi3.T) string {
	// Priority 1: Use the name explicitly set in the configuration
	if g.Config.CommandName != "" {
		return g.Config.CommandName
	}

	// Priority 2: Use the OpenAPI spec title as the fallback name
	rootusage := utils.GoCodeString(strings.TrimSpace(doc.Info.Title))
	if rootusage == "" {
		return "oasnake-cli" + fmt.Sprintf("%d", rand.IntN(1000)) // Default name if no title is provided
	}

	return rootusage
}

func (g *Generator) generateModel(swagger *openapi3.T) error {
	generatedModel, err := codegen.Generate(swagger, *g.Config.parserCodeGenConf)
	if err != nil {
		return fmt.Errorf("error generating model: %w", err)
	}
	if err := utils.WriteFileContent(
		utils.WriterConfig{
			OutputDirectoryShouldBeEmpty: false,
			Output: utils.FS{
				Directory: filepath.Join(g.Config.OutputDirectory, modelPath),
				Filename:  modelFileName,
			},
			Content: generatedModel,
		},
	); err != nil {
		return fmt.Errorf("error writing generated model to file: %w", err)
	}
	return nil
}

// GenerateCoreApp generates the core application files based on the command tree.
func (g *Generator) generateCoreApp(root *command.NodeCmd) error {
	type fileGen struct {
		Template TemplatorType
		Path     string
		Name     string
	}

	files := []fileGen{
		{CommonCommand, filepath.Join(g.Config.OutputDirectory, commandPath, command.CommonFolder), "utils.go"},
		{ConfigCommand, filepath.Join(g.Config.OutputDirectory, configPath), "command.go"},
		{ConfigRequest, filepath.Join(g.Config.OutputDirectory, configPath), "resuest.go"},
		{ConfigMethod, filepath.Join(g.Config.OutputDirectory, configPath), "method.go"},
		{ConfigExtension, filepath.Join(g.Config.OutputDirectory, configPath), "extension.go"},
		{Service, filepath.Join(g.Config.OutputDirectory, servicePath), "service.go"},
		{App, filepath.Join(g.Config.OutputDirectory, appPath), "app.go"},
	}

	if g.Config.WithCompilerFile {
		files = append(files,
			fileGen{Mod, g.Config.OutputDirectory, "go.mod"},
			fileGen{Main, g.Config.OutputDirectory, "main.go"},
		)
	}

	for _, f := range files {
		if err := NewTemplator(f.Template).WriteTemplateToFile(root, utils.FS{
			Directory: f.Path,
			Filename:  f.Name,
		}); err != nil {
			return fmt.Errorf("template generation failed (%s): %w", f.Name, err)
		}
	}

	return nil
}
