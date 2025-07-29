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
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/util"
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
func (g *Generator) Generate() (string, error) {
	log.Debug().Msg("Starting code generation")

	// Step 1: Load the OpenAPI specification
	swagger, opts, err := g.loadSwagger()
	if err != nil {
		return "", fmt.Errorf("failed to load OpenAPI specification: %w", err)
	}

	// Step 2: Construct the command tree from the OpenAPI paths
	rootCommand, err := g.toCommandTree(swagger)
	if err != nil {
		return "", fmt.Errorf("failed to construct command tree: %w", err)
	}

	// Step 3: Render CLI command files recursively
	cmdOutputPath := filepath.Join(g.Config.OutputDirectory, commandPath)
	if err := traverseAndRenderCommands(rootCommand, cmdOutputPath); err != nil {
		return "", fmt.Errorf("failed to render command files: %w", err)
	}

	// Step 4 (Optional): Generate API models using oapi-codegen
	if g.Config.WithModel {
		if err := g.generateModel(swagger, opts); err != nil {
			return "", fmt.Errorf("failed to generate API models: %w", err)
		}
	}

	// Step 5: Generate core application templates
	if err := g.generateCoreApp(rootCommand); err != nil {
		return "", fmt.Errorf("failed to generate core application templates: %w", err)
	}

	log.Info().Msgf("Code generation completed successfully. Output directory: %s", g.Config.OutputDirectory)
	return g.GetEffectiveRootUsage(swagger), nil
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

// GetEffectiveServerUrlConfig returns the effective server URL to be used,
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
func (g *Generator) GetEffectiveServerUrl(doc *openapi3.T) (string, error) {
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

// buildCommandTree parses OpenAPI paths into a hierarchical command structure.
func (g *Generator) toCommandTree(doc *openapi3.T) (*command.NodeCmd, error) {
	baseUrl, err := g.GetEffectiveServerUrl(doc)
	if err != nil {
		return nil, fmt.Errorf("error getting effective server URL: %w", err)
	}
	resolvePath := func(subPath string) string {
		if g.Config.ToInstall {
			return subPath
		}
		return filepath.Join(g.Config.OutputDirectory, subPath)
	}
	globalConfig := command.CommandGlobalConfig{
		RootUsage:   g.GetEffectiveRootUsage(doc),
		ModuleName:  g.Config.Module,
		BaseCmdPath: resolvePath(commandPath),
		ConfigPath:  resolvePath(configPath),
		AppPath:     resolvePath(appPath),
		ServicePath: resolvePath(servicePath),
		BaseUrl:     baseUrl,
	}
	rootNode := command.NewRootNodeCmd(globalConfig)
	for path, pathItem := range doc.Paths.Map() {
		segments := strings.Split(strings.Trim(path, "/"), "/")
		current := rootNode

		for _, segment := range segments {
			if _, exists := current.Children[segment]; !exists {
				current.Children[segment] = current.NewChildrenNodeCmd(segment)
			}
			current = current.Children[segment]
		}

		ops := map[command.Method]*openapi3.Operation{
			command.GET:    pathItem.Get,
			command.POST:   pathItem.Post,
			command.PUT:    pathItem.Put,
			command.PATCH:  pathItem.Patch,
			command.DELETE: pathItem.Delete,
		}

		for method, op := range ops {
			if op != nil {
				current.Methods[method] = op
			}
		}
	}
	return rootNode, nil
}

func (g *Generator) generateModel(swagger *openapi3.T, opts *codegen.Configuration) error {
	generatedModel, err := codegen.Generate(swagger, *opts)
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

	if g.Config.ToInstall {
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

// This function loads the OpenAPI specification, and parse into an oapi-codegen configuration.
func (g *Generator) loadSwagger() (*openapi3.T, *codegen.Configuration, error) {
	opts := codegen.Configuration{
		PackageName: "client",
		Generate: codegen.GenerateOptions{
			Client: true,
			Models: true,
		},
	}
	opts = opts.UpdateDefaults()

	overlayOpts := util.LoadSwaggerWithOverlayOpts{
		Path: opts.OutputOptions.Overlay.Path,
		// default to strict, but can be overridden
		Strict: true,
	}

	if opts.OutputOptions.Overlay.Strict != nil {
		overlayOpts.Strict = *opts.OutputOptions.Overlay.Strict
	}

	swagger, err := util.LoadSwaggerWithOverlay(g.Config.InputFilePath, overlayOpts)
	if err != nil {
		return nil, nil, fmt.Errorf("error loading OpenAPI spec: %w", err)
	}
	return swagger, &opts, nil
}
