package parser

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/louislouislouislouis/oasnake/app/pkg/generator/command"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/util"
	"github.com/rs/zerolog/log"
)

type Parser struct {
	Config Config
}

func NewParser(cfg Config) *Parser {
	newConf := cfg.ParserCodeGenConf.UpdateDefaults()
	cfg.ParserCodeGenConf = &newConf
	return &Parser{Config: cfg}
}

type ParseResult struct {
	rootCommand   *command.NodeCmd
	codeGenConfig *codegen.Configuration
}

func (p *Parser) ParseAndGetOpts() (*command.NodeCmd, *openapi3.T, error) {
	// Step 1: Load the OpenAPI specification
	swagger, err := p.loadSwagger()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load OpenAPI specification: %w", err)
	}

	// Step 2: Construct the command tree from the OpenAPI paths
	rootCommand, err := p.toCommandTree(swagger)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to construct command tree: %w", err)
	}
	return rootCommand, swagger, nil
}

func (p *Parser) loadSwagger() (*openapi3.T, error) {
	overlayOpts := util.LoadSwaggerWithOverlayOpts{
		Path:   p.Config.ParserCodeGenConf.OutputOptions.Overlay.Path,
		Strict: true,
	}

	if p.Config.ParserCodeGenConf.OutputOptions.Overlay.Strict != nil {
		overlayOpts.Strict = *p.Config.ParserCodeGenConf.OutputOptions.Overlay.Strict
	}

	swagger, err := util.LoadSwaggerWithOverlay(p.Config.InputFilePath, overlayOpts)
	if err != nil {
		return nil, fmt.Errorf("error loading OpenAPI spec: %w", err)
	}
	return swagger, nil
}

// buildCommandTree parses OpenAPI paths into a hierarchical command structure.
func (p *Parser) toCommandTree(doc *openapi3.T) (*command.NodeCmd, error) {
	rootNode := command.NewRootNodeCmd()
	for path, pathItem := range doc.Paths.Map() {
		segments := strings.Split(strings.Trim(path, "/"), "/")
		current := rootNode

	segmentLoop:
		for _, segment := range segments {
			if segment == "" {
				log.Debug().Msg("Add Root Segment")
				break segmentLoop
			}
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
