package events

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/louislouislouislouis/oasnake/app/pkg/generator/command"
)

type FinishParsingEvent struct {
	RootCmd *command.NodeCmd
	Spec    *openapi3.T
}

func (e FinishParsingEvent) Type() EventType { return FinishParsing }
