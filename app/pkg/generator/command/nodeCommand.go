package command

import (
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/louislouislouislouis/oasnake/app/pkg/utils"
	"github.com/rs/zerolog/log"
)

type CommandGlobalConfig struct {
	RootUsage   string
	ModuleName  string
	BaseUrl     string
	BaseCmdPath string
	ConfigPath  string
	AppPath     string
	ServicePath string
}

var CommonFolder = "common"

func (config CommandGlobalConfig) GetAppImportPath() string {
	return utils.TrimTrailingSlash(filepath.Join(
		config.ModuleName,
		config.AppPath,
	))
}

func (config CommandGlobalConfig) GetBaseCommandImportPath() string {
	return utils.TrimTrailingSlash(filepath.Join(
		config.ModuleName,
		config.BaseCmdPath,
	))
}

func (config CommandGlobalConfig) GetCommonImportPath() string {
	return utils.TrimTrailingSlash(filepath.Join(
		config.ModuleName,
		config.BaseCmdPath,
		CommonFolder,
	))
}

func (config CommandGlobalConfig) GetConfigImportPath() string {
	return utils.TrimTrailingSlash(filepath.Join(
		config.ModuleName,
		config.ConfigPath,
	))
}

func (config CommandGlobalConfig) GetServiceImportPath() string {
	return utils.TrimTrailingSlash(filepath.Join(
		config.ModuleName,
		config.ServicePath,
	))
}

type NodeCmd struct {
	GlobalConfig CommandGlobalConfig
	RelativePath string
	Parent       *NodeCmd
	segment      string
	Methods      map[Method]*openapi3.Operation
	Children     map[string]*NodeCmd
	depth        int
	paramDepth   int
}

func (node *NodeCmd) GetPath() string {
	if node.IsRootNodeCmd() {
		return "/"
	}

	var segments []string
	current := node
	for current != nil && !current.IsRootNodeCmd() {
		segments = append([]string{current.segment}, segments...)
		current = current.Parent
	}

	return "/" + strings.Join(segments, "/")
}

func (node *NodeCmd) GetDefaultMethod() Method {
	defaultMethod := GET
	if _, ok := node.Methods[GET]; !ok {
		for m := range node.Methods {
			defaultMethod = m
			break
		}
	}
	return defaultMethod
}

func (node *NodeCmd) GetAppModule() string {
	return path.Base(node.GlobalConfig.AppPath)
}

func (node *NodeCmd) GetLongDescription() string {
	return node.getCmdDescription(false)
}

func (node *NodeCmd) GetShortDescription() string {
	return node.getCmdDescription(true)
}

func (node *NodeCmd) GetImportPath() string {
	return utils.TrimTrailingSlash(filepath.Join(
		node.GlobalConfig.ModuleName,
		node.GlobalConfig.BaseCmdPath,
		node.RelativePath,
	))
}

func (node *NodeCmd) GetQueryParams() map[string]Parameter {
	return node.getParams("query")
}

func (node *NodeCmd) GetHeaderParams() map[string]Parameter {
	return node.getParams("header")
}

func (node *NodeCmd) getParams(paramType string) map[string]Parameter {
	params := make(map[string]Parameter)
	for _, operation := range node.Methods {
		for _, item := range operation.Parameters {
			if v := item.Value; v != nil {
				if v.In == paramType {
					params[v.Name] = Parameter{*v}
				}
			}
		}
	}
	return params
}

func (node *NodeCmd) getCmdDescription(isShort bool) string {
	var builder strings.Builder
	for method, operation := range node.Methods {
		description := ""

		if isShort {
			description = utils.RemoveBackTicks(operation.Summary)
		} else {
			description = utils.RemoveBackTicks(operation.Description)
		}
		builder.WriteString("\n" + strings.ToUpper(string(method)))
		builder.WriteString("\n" + description)
		builder.WriteString("\n----------------------")
	}

	return builder.String()
}

func (node *NodeCmd) GetFileName() string {
	if node.IsRootNodeCmd() {
		return "root.go"
	}
	return node.GetPackageName() + ".go"
}

func (node *NodeCmd) GetParamName() string {
	if node.IsParam() {
		return strings.Trim(node.segment, "{}")
	}
	return ""
}

func (node *NodeCmd) GetPackageName() string {
	var name string

	if node.IsRootNodeCmd() {
		return path.Base(node.GlobalConfig.BaseCmdPath)
	} else if node.IsParam() {
		name = strings.Trim(node.segment, "{}")
	} else {
		name = node.segment
	}

	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "-", "_")

	if name == "" {
		log.Error().Msg("Package name is empty, defaulting to 'cmd'")
	}

	// Check reserved keyword
	if slices.Contains(reservedPackageNames, name) {
		name = name + "_cmd"
	}

	return name
}

var reservedPackageNames = []string{
	"break", "default", "func", "interface", "select",
	"case", "defer", "go", "map", "struct",
	"chan", "else", "goto", "package", "switch",
	"const", "fallthrough", "if", "range", "type",
	"continue", "for", "import", "return", "var",
	"init",
	// Custom
	"main", "cmd", "config",
}

func (node *NodeCmd) GetUsage() string {
	if node.IsRootNodeCmd() {
		return node.GlobalConfig.RootUsage
	}
	if node.IsParam() {
		return "<" + strings.Trim(node.segment, "{}") + ">"
	}
	return node.GetPackageName()
}

func (node *NodeCmd) GetCobraFunctionCommandName() string {
	commandName := node.segment
	if node.IsRootNodeCmd() {
		commandName = node.GlobalConfig.RootUsage
	}
	if node.IsParam() {
		commandName = strings.Trim(node.segment, "{}")
	}
	// String manipulation to ensure the command name is properly formatted
	commandName = utils.CapitalizeFirstOnly(commandName)
	return utils.GoCodeString(commandName)
}

func (node *NodeCmd) IsRootNodeCmd() bool {
	return node.Parent == nil
}

func (node *NodeCmd) NewChildrenNodeCmd(segment string) *NodeCmd {
	childNodeCmd := newNodeCmd(segment)
	childNodeCmd.depth = node.depth + 1
	childNodeCmd.Parent = node
	childNodeCmd.RelativePath = node.RelativePath + childNodeCmd.GetPackageName() + "/"
	return childNodeCmd
}

func (node *NodeCmd) IsParam() bool {
	return strings.HasPrefix(node.segment, "{") && strings.HasSuffix(node.segment, "}")
}

func newNodeCmd(segment string) *NodeCmd {
	return &NodeCmd{
		depth:      0,
		paramDepth: 0,
		Parent:     nil,
		segment:    segment,
		Methods:    make(map[Method]*openapi3.Operation),
		Children:   make(map[string]*NodeCmd),
	}
}

func NewRootNodeCmd() *NodeCmd {
	node := newNodeCmd("")
	node.RelativePath = "/"
	return node
}

func (node *NodeCmd) SetGlobalConfig(config CommandGlobalConfig) {
	node.GlobalConfig = config
	for _, child := range node.Children {
		child.SetGlobalConfig(config)
	}
}
